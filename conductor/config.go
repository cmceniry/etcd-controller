package conductor

import (
	"fmt"

	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

type NodeInfoConfig struct {
	Name        string
	IP          string `yaml:"IP"`
	CommandPort int    `yaml:"CommandPort"`
	Insecure    bool   `yaml:"Insecure"`
	PeerPort    int    `yaml:"PeerPort"`
	ClientPort  int    `yaml:"ClientPort"`
}

type NodeListConfig []NodeInfoConfig

func (c *Conductor) LoadYaml(d []byte) error {
	var data []NodeInfoConfig
	err := yaml.Unmarshal(d, &data)
	if err != nil {
		return err
	}
	for _, n := range data {
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
		fmt.Printf("%#v\n", ni)
		c.NodeList[n.Name] = ni
	}
	return err
}
