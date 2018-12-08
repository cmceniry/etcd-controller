package conductor

import (
	"fmt"

	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
)

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
