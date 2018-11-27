package conductor

import (
	"io/ioutil"
	"os"

	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	NodeListFilename string
	CommandPort      int
}

type NodeInfoConfig struct {
	Name        string
	IP          string `yaml:"IP"`
	CommandPort int    `yaml:"CommandPort"`
	Insecure    bool   `yaml:"Insecure"`
	PeerPort    int    `yaml:"PeerPort"`
	ClientPort  int    `yaml:"ClientPort"`
}

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
			Name:        n.Name,
			IP:          n.IP,
			CommandPort: n.CommandPort,
			PeerPort:    n.PeerPort,
			ClientPort:  n.ClientPort,
		}
		if n.Insecure {
			ni.CommandOpts = []grpc.DialOption{grpc.WithInsecure()}
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
