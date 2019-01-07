package conductor

import (
	"fmt"
	"sort"

	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
)

func (c *Conductor) getNodeNames() []string {
	t := map[string]struct{}{}
	for n := range c.NodeList {
		t[n] = struct{}{}
	}
	for n := range c.CurrentNodes {
		t[n] = struct{}{}
	}
	ret := []string{}
	for n := range t {
		ret = append(ret, n)
	}
	sort.Strings(ret)
	return ret
}

func (c *Conductor) printNodeStatus(nodeName, format string) string {
	n, ok := c.CurrentNodes[nodeName]
	if !ok {
		return fmt.Sprintf(format, n.Name, "NOTCURRENT")
	}
	s, ok := pb.NodeInfoStatus_name[n.Status]
	if !ok {
		return fmt.Sprintf(format, n.Name, "BADSTATUS")
	}
	return fmt.Sprintf(format, n.Name, s)
}

func (c *Conductor) printNodesStatus() string {
	width := 30
	nl := c.getNodeNames()
	for _, nn := range nl {
		if len(nn) < width {
			width = len(nn)
		}
	}
	format := fmt.Sprintf("- %%-%ds %%s\n", width)
	ret := "NodeList:\n"
	for _, nn := range nl {
		ret += c.printNodeStatus(nn, format)
	}
	return ret
}

func (c *Conductor) getClusterNodeStatus() error {
	for _, nn := range c.currentNodeNames() {
		ni, ok := c.CurrentNodes[nn]
		if !ok {
			return fmt.Errorf("missing current node: %s", nn)
		}
		// TODO: if ni.ExternalETCD { continue }
		dc, err := c.connectCommand(ni)
		if err != nil {
			ni.Status = int32(pb.NodeInfoStatus_UNKNOWN)
			return fmt.Errorf("connect failed %s: %s", nn, err)
		}
		s, err := dc.Status()
		if err != nil {
			ni.Status = int32(pb.NodeInfoStatus_UNKNOWN)
			return fmt.Errorf("status failued: %s: %s", nn, err)
		}
		ni.Status = s
	}
	return nil
}
