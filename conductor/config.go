package conductor

import (
	"io/ioutil"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Config is the simple configuration pieces for Conductor
type Config struct {
	NodeListFilename string
	CommandPort      int
	Logger           *log.Entry
	PeerTLSCA        string
	PeerTLSCert      string
	PeerTLSKey       string
	ClientTLSCA      string
	ClientTLSCert    string
	ClientTLSKey     string
}

// NodeInfoConfig is the data stored in a node list file which can be used to
// populate the authoritative list of what should be in the cluster
type NodeInfoConfig struct {
	Name        string
	IP          string `yaml:"IP"`
	CommandPort int    `yaml:"CommandPort"`
	Insecure    bool   `yaml:"Insecure"`
	PeerPort    int    `yaml:"PeerPort"`
	ClientPort  int    `yaml:"ClientPort"`
}

// NodeListConfig is an array of node configuration
type NodeListConfig []NodeInfoConfig

func (c *Conductor) checkNodeList() (bool, error) {
	fi, err := os.Stat(c.Config.NodeListFilename)
	if err != nil {
		return false, err
	}
	if !c.lastNodeListRead.Before(fi.ModTime()) {
		return false, nil
	}
	d, err := ioutil.ReadFile("/config/node-list.yaml")
	if err != nil {
		return false, err
	}
	err = c.loadYaml(d)
	if err != nil {
		return false, err
	}
	c.lastNodeListRead = fi.ModTime()
	return true, nil
}

func (c *Conductor) loadYaml(d []byte) error {
	var data []NodeInfoConfig
	err := yaml.Unmarshal(d, &data)
	if err != nil {
		return err
	}
	seen := map[string]struct{}{}
	for _, n := range data {
		seen[n.Name] = struct{}{}
		ni := &NodeInfo{
			Name:         n.Name,
			IP:           n.IP,
			CommandPort:  n.CommandPort,
			PeerPort:     n.PeerPort,
			PeerSecure:   c.peerTLS,
			ClientPort:   n.ClientPort,
			ClientSecure: c.clientTLS,
		}

		// TODO change here ?
		if n.Insecure {
			// ni.CommandOpts = []grpc.DialOption{grpc.WithInsecure()}
		}
		c.NodeList[n.Name] = ni
	}
	for newNodeName := range c.NodeList {
		if _, ok := seen[newNodeName]; !ok {
			delete(c.NodeList, newNodeName)
		}
	}
	return err
}

func (c *Conductor) currentNodeNames() []string {
	cnl := []string{}
	for nn := range c.CurrentNodes {
		cnl = append(cnl, nn)
	}
	sort.Strings(cnl)
	return cnl
}
