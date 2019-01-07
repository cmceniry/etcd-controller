package conductor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"time"

	"github.com/cmceniry/etcd-controller/driver"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	running          bool

	Listener   net.Listener
	GRPCServer *grpc.Server
	clientTLS  bool
	peerTLS    bool
}

// NodeInfo holds simple information about the nodes that this Conductor is
// watching.
type NodeInfo struct {
	Name         string
	IP           string
	CommandPort  int
	PeerPort     int
	PeerSecure   bool
	ClientPort   int
	ClientSecure bool
	Status       int32
}

// ClientProto returns the protocol (http/https) for the client connection
func (n NodeInfo) ClientProto() string {
	if n.ClientSecure {
		return "https"
	}
	return "http"
}

// ClientURL returns the URL for etcd clients to connect to on this node
func (n NodeInfo) ClientURL() string {
	return fmt.Sprintf("%s://%s:%d", n.ClientProto(), n.IP, n.ClientPort)
}

// PeerProto returns the protocol (http/https) for the peer connections
func (n NodeInfo) PeerProto() string {
	if n.PeerSecure {
		return "https"
	}
	return "http"
}

// PeerURL returns the URL for etcd peer to connect to on this node
func (n NodeInfo) PeerURL() string {
	return fmt.Sprintf("%s://%s:%d", n.PeerProto(), n.IP, n.PeerPort)
}

// PeerString returns the value to use in the cluster node list
func (n NodeInfo) PeerString() string {
	return fmt.Sprintf("%s=%s", n.Name, n.PeerURL())
}

// IsExtra returns if a node is in the CurrentNodes but not in the official
// node list
func (c *Conductor) IsExtra(nodeName string) bool {
	_, ok := c.NodeList[nodeName]
	return !ok
}

// NewConductor is the general Conductor constructor.
func NewConductor(c Config) *Conductor {
	con := &Conductor{
		Config:       &c,
		NodeList:     make(map[string]*NodeInfo),
		CurrentNodes: make(map[string]*NodeInfo),
	}
	if c.PeerTLSCA != "" && c.PeerTLSCert != "" && c.PeerTLSKey != "" {
		con.peerTLS = true
	}
	if c.ClientTLSCA != "" && c.ClientTLSCert != "" && c.ClientTLSKey != "" {
		con.clientTLS = true
	}
	return con
}

func (c *Conductor) main() {

}

func (c *Conductor) connectCommand(ni *NodeInfo) (*driver.SimpleClient, error) {
	do := []grpc.DialOption{}
	if c.peerTLS {
		cert, err := tls.LoadX509KeyPair(c.Config.PeerTLSCert, c.Config.PeerTLSKey)
		if err != nil {
			return nil, err
		}
		caPool := x509.NewCertPool()
		caData, err := ioutil.ReadFile(c.Config.PeerTLSCA)
		if err != nil {
			return nil, err
		}
		if ok := caPool.AppendCertsFromPEM(caData); !ok {
			return nil, fmt.Errorf("unable to load ca certs")
		}
		do = append(do, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caPool,
			ServerName:   ni.IP,
		})))
	} else {
		do = append(do, grpc.WithInsecure())
	}
	return driver.NewSimpleClient(ni.IP, ni.CommandPort, do)
}

func (c *Conductor) etcdDial(ni *NodeInfo) (*clientv3.Client, error) {
	var tlsConfig *tls.Config
	if c.clientTLS {
		tlsInfo := transport.TLSInfo{
			TrustedCAFile: c.Config.ClientTLSCA,
			CertFile:      c.Config.ClientTLSCert,
			KeyFile:       c.Config.ClientTLSKey,
		}
		t, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		tlsConfig = t
	}
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{ni.ClientURL()},
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
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

func (c *Conductor) etcdctlGetMemberID(ctlNode *NodeInfo, needleNode *NodeInfo) (uint64, error) {
	client, err := c.etcdDial(ctlNode)
	if err != nil {
		return 0, err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	list, err := client.MemberList(ctx)
	cancel()
	if err != nil {
		return 0, err
	}
	for _, m := range list.Members {
		if m.GetName() == needleNode.Name {
			return m.GetID(), nil
		}
	}
	return 0, fmt.Errorf("node not found %s", needleNode.Name)
}

func (c *Conductor) initNewCluster(initNodeName string) error {
	// TODO: shutdown all nodes
	initNode, ok := c.CurrentNodes[initNodeName]
	if !ok {
		return errors.Errorf("unknown init node %s", initNodeName)
	}
	dc, err := c.connectCommand(initNode)
	if err != nil {
		return errors.Wrap(err, "failed to connect")
	}
	err = dc.InitCluster()
	if err != nil {
		return errors.Wrap(err, "init failed")
	}
	time.Sleep(5 * time.Second)
	err = c.etcdctlStatus(initNode)
	if err != nil {
		return errors.Wrap(err, "init status failed")
	}
	return nil
}

func (c *Conductor) startNode(nodeName string) error {
	node, ok := c.CurrentNodes[nodeName]
	if !ok {
		return fmt.Errorf("Unknown node: %s", nodeName)
	}
	dc, err := c.connectCommand(node)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}
	err = dc.Start()
	return err
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
	newNode, ok := c.CurrentNodes[newNodeName]
	if !ok {
		return fmt.Errorf("Unknown add node: %s", newNodeName)
	}
	dc, err := c.connectCommand(newNode)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}
	ctlNode, ok := c.CurrentNodes[c.pickRandomUpNode()]
	if !ok {
		return fmt.Errorf("no up nodes")
	}
	adderID, newID, err := c.etcdctlMemberAdd(ctlNode, newNode)
	if err != nil {
		return fmt.Errorf("member add failed: %s", err)
	}
	fmt.Printf("AdderID: %x, New Member ID: %x\n", adderID, newID)
	peerList := c.generatePeerList()
	if len(peerList) == 0 {
		return fmt.Errorf("peer list is empty")
	}
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
	return nil
}

// check if we have no running or runable members and indicate that we need
// to build a cluster
func (c *Conductor) checkNeedNewCluster() bool {
	if len(c.CurrentNodes) == 0 {
		return false
	}
	for _, ni := range c.CurrentNodes {
		if ni.IsRunning() || ni.IsStopped() {
			return false
		}
	}
	return true
}

func (c *Conductor) pickRandomNode() string {
	for n := range c.CurrentNodes {
		return n
	}
	return ""
}

func (c *Conductor) pickRandomUpNode() string {
	for nn, ni := range c.CurrentNodes {
		if ni.IsRunning() {
			return nn
		}
	}
	return ""
}

func (c *Conductor) pickRandomMissingNode() string {
	for nn, ni := range c.CurrentNodes {
		if !ni.IsRunning() && !ni.IsStopped() {
			return nn
		}
	}
	return ""
}

func (c *Conductor) generatePeerList() []string {
	ret := []string{}
	for _, ni := range c.CurrentNodes {
		if ni.IsRunning() {
			ret = append(ret, ni.PeerString())
		}
	}
	return ret
}

func (c *Conductor) IsRunning() bool {
	return c.running
}

// Run starts the main Conductor work loop
func (c *Conductor) Run() {
	c.running = true
	for t := range time.NewTicker(5 * time.Second).C {
		fmt.Printf("%s TICK!\n", t)
		changed, err := c.checkNodeList()
		if err != nil {
			fmt.Printf(`Error getting node list "%s": %s`+"\n", c.Config.NodeListFilename, err)
			continue
		}
		if changed {
			var dnl string
			if len(c.NodeList) > 0 {
				dnl = "New node list:\n"
				for _, ni := range c.NodeList {
					if _, ok := c.CurrentNodes[ni.Name]; !ok {
						c.CurrentNodes[ni.Name] = &(*ni)
					}
					dnl += fmt.Sprintf("    - %#v\n", ni)
				}
			} else {
				dnl = "New node list: empty\n"
			}
			fmt.Printf(dnl)
		}
		// TODO: Check all current nodes health
		err = c.getClusterNodeStatus()
		if err != nil {
			fmt.Printf(`Error getting cluster node status: %s`+"\n", err)
		}
		// Show status
		fmt.Printf(c.printNodesStatus())
		// Check if any current nodes have stopped and should attempt a restart
		for nodeName, nodeInfo := range c.CurrentNodes {
			if nodeInfo.IsStopped() {
				fmt.Printf("Starting stopped node: %s\n", nodeName)
				err := c.startNode(nodeName)
				if err != nil {
					fmt.Printf("Error trying to ensure node is running: %s: %s\n", nodeName, err)
				}
			}
		}
		// TODO: Check if there are extra nodes and remove them first
		err = c.removeExtraNodesFromCluster()
		if err != nil {
			fmt.Printf("removeExtraNodesFromCluster: %s", err)
			continue
		}
		// TODO: Check if etcd cluster has nodes not in current node list
		// If empty cluster, init it
		if c.checkNeedNewCluster() {
			initNodeName := c.pickRandomNode()
			fmt.Printf("Initializing Cluser with node %s\n", initNodeName)
			err := c.initNewCluster(initNodeName)
			if err != nil {
				fmt.Printf("Init Node Failure: %s: %s\n", initNodeName, err)
			} else {
				fmt.Printf("Initialization successful\n")
			}
			continue
		}
		// If there are uninitialized/failed nodes from the node list, add one of them at random
		if newNodeName := c.pickRandomMissingNode(); newNodeName != "" {
			fmt.Printf("Adding missing node to cluster: %s\n", newNodeName)
			err := c.addNodeToCluster(newNodeName)
			if err != nil {
				fmt.Printf("Add Node Failure: %s: %s\n", newNodeName, err)
			}
			fmt.Printf("Addition successful\n")
			continue
		}
		fmt.Printf("Nothing to do\n")
		time.Sleep(1 * time.Second)
	}
}
