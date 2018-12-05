package conductor

import "fmt"

func (c *Conductor) getClusterNodeStatus() error {
	for nn, ni := range c.CurrentNodes {
		// TODO: if ni.ExternalETCD { continue }
		dc, err := c.connectCommand(ni)
		if err != nil {
			return fmt.Errorf("connect failed %s: %s", nn, err)
		}
		s, err := dc.Status()
		if err != nil {
			return fmt.Errorf("status failued: %s: %s", nn, err)
		}
		ni.Status = s
	}
	return nil
}
