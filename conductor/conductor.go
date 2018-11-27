package conductor

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/cmceniry/etcd-controller/driver"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc"
)

// The Conductor is the unified driving entity. Only one runs at a time in the
// cluster. It interacts with the running etcd processes and with the other
// controllers in order to maintain the cluster.
//
// The Conductor's goal is to ensure that all nodes on the "official" node
// list are active members in the cluster.
type Conductor struct {
	Config           *Config
	NodeList         map[string]*NodeInfo
	CurrentNodes     map[string]*NodeInfo
	lastNodeListRead time.Time

	Listener net.Listener
	GRPCServer *grpc.Server
}

type NodeInfo struct {
	Name        string
	IP          string
	CommandPort int
	CommandOpts []grpc.DialOption
	PeerPort    int
	ClientPort  int
	Status      int32
}

func (n NodeInfo) ClientURL() string {
	return fmt.Sprintf("http://%s:%d", n.IP, n.ClientPort)
}

func (n NodeInfo) PeerURL() string {
	return fmt.Sprintf("http://%s:%d", n.IP, n.PeerPort)
}

func (n NodeInfo) PeerString() string {
	return fmt.Sprintf("%s=http://%s:%d", n.Name, n.IP, n.PeerPort)
}

// NewConductor is the general Conductor constructor.
func NewConductor(c Config) *Conductor {
	return &Conductor{
		Config:       &c,
		NodeList:     make(map[string]*NodeInfo),
		CurrentNodes: make(map[string]*NodeInfo),
	}
}

func (c *Conductor) main() {

}

func (c *Conductor) connectCommand(ni *NodeInfo) (*driver.SimpleClient, error) {
	return driver.NewSimpleClient(ni.IP, ni.CommandPort, ni.CommandOpts)
}

func (c *Conductor) etcdDial(ni *NodeInfo) (*clientv3.Client, error)  {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{ni.ClientURL()},
		DialTimeout: 5 * time.Second,
	})
}

func (c *Conductor) etcdctlStatus(ni *NodeInfo) error {
	client, err := c.etcdDial(ni)
	if err != nil {
		return err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Status(ctx, ni.ClientURL())
	cancel()
	if err == nil {
		fmt.Printf("%#v\n", resp)
	}
	switch err {
	case nil, rpctypes.ErrPermissionDenied:
		return nil
	case rpctypes.ErrTimeout:
		return err
	default:
		return err
	}
}

// Used to stall for time after operations that are likely to trigger known
// leader elections (i.e. member changes)
func (c *Conductor) etcdctlWaitForMaster(ni *NodeInfo) error {
	client, err := c.etcdDial(ni)
	if err != nil {
		return err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	_, err = client.Get(ctx, "_")
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (c *Conductor) etcdctlMemberAdd(ctlNode *NodeInfo, newNode *NodeInfo) (uint64, uint64, error) {
	client, err := c.etcdDial(ctlNode)
	if err != nil {
		return 0, 0, err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.MemberAdd(ctx, []string{newNode.PeerURL()})
	cancel()
	if err != nil {
		return 0, 0, fmt.Errorf("MemberAdd failed: %s", err)
	}
	return resp.Header.MemberId, resp.Member.ID, nil
}

func (c *Conductor) initNewCluster(initNodeName string) error {
	// TODO: shutdown all nodes
	initNode, ok := c.NodeList[initNodeName]
	if !ok {
		return fmt.Errorf("Unknown init node: %s", initNodeName)
	}
	dc, err := c.connectCommand(initNode)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}
	err = dc.InitCluster()
	if err != nil {
		return fmt.Errorf("init failed: %s", err)
	}
	time.Sleep(5 * time.Second)
	err = c.etcdctlStatus(initNode)
	if err != nil {
		return fmt.Errorf("init status failed: %s", err)
	}
	c.CurrentNodes[initNodeName] = &(*initNode)
	return nil
}

func (c *Conductor) etcdctlMemberList(ni *NodeInfo) (map[string]uint64, error) {
	client, err := c.etcdDial(ni)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := client.MemberList(ctx)
	if err != nil {
		return nil, err
	}
	ret := map[string]uint64{}
	for _, m := range resp.Members {
		if m.Name != "" {
			ret[m.Name] = m.ID
		}
	}
	return ret, nil
}

func (c *Conductor) addNodeToCluster(newNodeName string) error {
	newNode, ok := c.NodeList[newNodeName]
	if !ok {
		return fmt.Errorf("Unknown add node: %s", newNodeName)
	}
	dc, err := c.connectCommand(newNode)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}
	peerList := []string{}
	var ctlNode *NodeInfo
	for _, ni := range c.CurrentNodes {
		ctlNode = ni
		peerList = append(peerList, ni.PeerString())
	}
	adderID, newID, err := c.etcdctlMemberAdd(ctlNode, newNode)
	if err != nil {
		return fmt.Errorf("member add failed: %s", err)
	}
	fmt.Printf("AdderID: %x, New Member ID: %x\n", adderID, newID)
	peerList = append(peerList, newNode.PeerString())
	err = dc.JoinCluster(peerList)
	if err != nil {
		return fmt.Errorf("join failed: %s", err)
	}
	time.Sleep(6 * time.Second)
	err = c.etcdctlStatus(newNode)
	if err != nil {
		return fmt.Errorf("join status after master failed: %s", err)
	}
	retry := -1
	for {
		retry++
		if retry > 6 {
			return fmt.Errorf("memberlist not converge after join")
		}
		ctlMembers, err := c.etcdctlMemberList(ctlNode)
		if err != nil {
			fmt.Printf("%s memberlist failed: %s", ctlNode.PeerString(), err)
			continue
		}
		newMembers, err := c.etcdctlMemberList(newNode)
		if err != nil {
			return fmt.Errorf("%s memberlist failed: %s", newNode.PeerString(), err)
		}
		if reflect.DeepEqual(ctlMembers, newMembers) {
			fmt.Printf("members: %#v\n", newMembers)
			break
		}
		time.Sleep(1 * time.Second)
	}
	err = c.etcdctlStatus(newNode)
	if err != nil {
		return fmt.Errorf("join status failed: %s", err)
	}
	c.CurrentNodes[newNodeName] = &(*newNode)
	return nil
}

func (c *Conductor) pickRandomNodeFromList() string {
	if len(c.NodeList) == 0 {
		return ""
	}
	for n := range c.NodeList {
		return n
	}
	return ""
}

func (c *Conductor) pickRandomMissingNodeFromList() string {
	for n := range c.NodeList {
		if _, ok := c.CurrentNodes[n]; !ok {
			return n
		}
	}
	return ""
}

func (c *Conductor) Run() {
	c.runGRPCListener()

	lastRun := time.Now().Add(-1 * time.Second)
	for {
		throttle := lastRun.Add(1 * time.Second).Sub(lastRun)
		if throttle > time.Second {
			throttle = time.Second
		}
		time.Sleep(throttle)
		lastRun = time.Now()
		changed, err := c.checkNodeList()
		if err != nil {
			fmt.Printf(`Error getting node list "%s": %s`+"\n", c.Config.NodeListFilename, err)
			continue
		}
		if changed {
			fmt.Printf("New node list:\n%#v\n", c.NodeList)
			if len(c.NodeList) > 0 {
				for _, ni := range c.NodeList {
					fmt.Printf("- %#v\n", ni)
				}
			} else {
				fmt.Printf("- empty\n")
			}
		}
		// TODO: Check all current nodes health
		// TODO: Check if etcd cluster has nodes not in current node list
		// TODO: Check if there are extra nodes and remove them first
		// If empty cluster, init it
		if len(c.CurrentNodes) == 0 && len(c.NodeList) > 0 {
			initNodeName := c.pickRandomNodeFromList()
			fmt.Printf("Initializing Cluser with node %s\n", initNodeName)
			err := c.initNewCluster(initNodeName)
			if err != nil {
				fmt.Printf("Init Node Failure: %s: %s\n", initNodeName, err)
			}
			fmt.Printf("Initialization successful\n")
			continue
		}
		// If there are missing nodes from the node list, add one of them at random
		if newNodeName := c.pickRandomMissingNodeFromList(); newNodeName != "" {
			fmt.Printf("Adding missing node to cluster: %s\n", newNodeName)
			err := c.addNodeToCluster(newNodeName)
			if err != nil {
				fmt.Printf("Add Node Failure: %s: %s", newNodeName, err)
			}
			fmt.Printf("Addition successful\n")
			continue
		}
		fmt.Printf("Nothing to do\n")
		time.Sleep(1 * time.Second)
	}

	// err := c.initNewCluster("etcd-controller-test-005-1")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("init successful\n")
	// err = c.addNodeToCluster("etcd-controller-test-005-2")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("2 join successful\n")
	// err = c.addNodeToCluster("etcd-controller-test-005-3")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("3 join successful\n")
}
