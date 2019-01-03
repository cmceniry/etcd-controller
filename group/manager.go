package group

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"time"

	"github.com/cmceniry/etcd-controller/conductor"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	yaml "gopkg.in/yaml.v2"
)

type nodeinfo struct {
	Name     string
	IP       string
	SerfPort int
	Info     *conductor.NodeInfo
}

// HACK - remove these later - dup with conductor
func (m *Manager) checkNodeList() (bool, error) {
	fi, err := os.Stat(m.Config.NodeListFilename)
	if err != nil {
		return false, err
	}
	if !m.lastNodeListRead.Before(fi.ModTime()) {
		return false, nil
	}
	d, err := ioutil.ReadFile(m.Config.NodeListFilename)
	if err != nil {
		return false, err
	}
	err = m.loadYaml(d)
	if err != nil {
		return false, err
	}
	m.lastNodeListRead = fi.ModTime()
	return true, nil
}

func (m *Manager) loadYaml(d []byte) error {
	var data []conductor.NodeInfoConfig
	err := yaml.Unmarshal(d, &data)
	if err != nil {
		return err
	}
	seen := map[string]struct{}{}
	for _, n := range data {
		seen[n.Name] = struct{}{}
		ni := nodeinfo{
			Name:     n.Name,
			IP:       n.IP,
			SerfPort: m.Config.SerfPort,
		}
		m.Nodes[n.Name] = &ni
	}
	for newNodeName := range m.Nodes {
		if _, ok := seen[newNodeName]; !ok {
			delete(m.Nodes, newNodeName)
		}
	}
	return err
}

// HACK - remove above

// Manager is the process which maintains the overall cluster membership and
// decides on who should be the Conductor
type Manager struct {
	Config           Config
	Nodes            map[string]*nodeinfo
	SerfEvents       chan serf.Event
	MemberListConfig *memberlist.Config
	ml               *memberlist.Memberlist
	SerfConfig       *serf.Config
	s                *serf.Serf
	lastNodeListRead time.Time
}

// Config holds the parameters for the Manager
type Config struct {
	Name             string
	IP               string
	SerfPort         int
	NodeListFilename string
}

// NewManager returns a Grouper Manager
func NewManager(c Config) (*Manager, error) {
	m := &Manager{}

	m.Config = c
	m.Nodes = map[string]*nodeinfo{}

	m.MemberListConfig = memberlist.DefaultLANConfig()
	m.MemberListConfig.BindAddr = "0.0.0.0"
	m.MemberListConfig.AdvertiseAddr = c.IP
	m.MemberListConfig.BindPort = c.SerfPort
	m.MemberListConfig.AdvertisePort = c.SerfPort
	m.MemberListConfig.LogOutput = os.Stdout

	m.SerfEvents = make(chan serf.Event, 8)
	m.SerfConfig = serf.DefaultConfig()
	m.SerfConfig.NodeName = c.Name
	m.SerfConfig.EventCh = m.SerfEvents
	m.SerfConfig.MemberlistConfig = m.MemberListConfig
	m.SerfConfig.LogOutput = os.Stdout

	s, err := serf.Create(m.SerfConfig)
	if err != nil {
		return nil, err
	}
	m.s = s

	return m, nil
}

// Run starts the main loop for the Grouper
func (m *Manager) Run() {
	go m.main()
}

// UpdateNodeList matches up with what is currently there
func (m *Manager) UpdateNodeList(nl map[string]*conductor.NodeInfo) error {
	return nil
}

func netID(ipStr string, port int) uint64 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return math.MaxUint64
	}
	return uint64(binary.BigEndian.Uint32(ip[12:16]))<<16 + uint64(port)
}

func (m *Manager) isLowestNet() bool {
	myNI, ok := m.Nodes[m.Config.Name]
	if !ok {
		return false
	}
	myID := netID(myNI.IP, myNI.SerfPort)
	for nn, ni := range m.Nodes {
		if nn == myNI.Name {
			continue
		}
		otherNetID := netID(ni.IP, ni.SerfPort)
		if myID > otherNetID {
			return false
		}
	}
	return true
}

// IsConductor indicates that this node should be running the conductor
// component
func (m *Manager) IsConductor() (bool, error) {
	// if manual!=nil && manual!= me, return false
	// if manual==me, return true

	// If total cluster, use lowest
	seen := []string{}
	for _, member := range m.s.Members() {
		// if member.isnothealth { continue }
		if _, ok := m.Nodes[member.Name]; ok {
			seen = append(seen, member.Name)
			continue
		}
		// Extra member
	}
seenloop:
	for name := range m.Nodes {
		for _, sn := range seen {
			if sn == name {
				continue seenloop
			}
		}
		return false, fmt.Errorf("not all nodes present")
	}
	// Am I the lowest
	isLow := m.isLowestNet()
	return isLow, fmt.Errorf("isLowestNet=%t", isLow)
}

func (m *Manager) checkSerfPeers() error {
	// attempt to join members in node list but not in serf
	join := []string{}
nodeloop:
	for nn, ni := range m.Nodes {
		for _, member := range m.s.Members() {
			if member.Name == nn {
				continue nodeloop
			}
		}
		join = append(join, fmt.Sprintf("%s:%d", ni.IP, ni.SerfPort))
	}
	if len(join) == 0 {
		return nil
	}
	_, err := m.s.Join(join, true)
	return err
}

func (m *Manager) main() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("Grouper TICK!\n")
			changed, err := m.checkNodeList()
			// fmt.Printf("nodelist: %#v\n", m.Nodes)
			// fmt.Printf("serfnode: %#v\n", m.s.Members())
			if err != nil {
				fmt.Printf("error checking node list: %s\n", err)
				continue
			}
			if !changed {
				continue
			}
			if err := m.checkSerfPeers(); err != nil {
				fmt.Printf("error checking peers: %s\n", err)
			}
		case e := <-m.SerfEvents:
			fmt.Printf("%#v\n", e)
			if me, ok := e.(serf.MemberEvent); ok {
				for _, member := range me.Members {
					fmt.Printf("MemberEvent: %#v\n", member)
					switch me.EventType() {
					case serf.EventMemberJoin:
						fmt.Printf("Join: %#v\n", me)
					case serf.EventMemberLeave, serf.EventMemberFailed, serf.EventMemberReap:
						fmt.Printf("Out: %#v\n", me)
					}
				}
			}
		}
	}
}
