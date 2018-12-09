package conductor

import (
	"context"
	"fmt"
	"time"
)

func (c *Conductor) etcdctlRemoveNode(ctlNode *NodeInfo, rmNode *NodeInfo) error {
	rmID, err := c.etcdctlGetMemberID(ctlNode, rmNode)
	if err != nil {
		// Node already removed
		return nil
	}
	client, err := c.etcdDial(ctlNode)
	if err != nil {
		return err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = client.MemberRemove(ctx, rmID)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (c *Conductor) findMissingExtraNodes() []string {
	ret := []string{}
	for nn, ni := range c.CurrentNodes {
		if !c.IsExtra(nn) {
			continue
		}
		// Missing means it's not available for the cluster
		if ni.IsRunning() || ni.IsStopped() {
			continue
		}
		ret = append(ret, nn)
	}
	return ret
}

func (c *Conductor) removeNodeFromCluster(rmNodeName string) error {
	rmNode, ok := c.CurrentNodes[rmNodeName]
	if !ok {
		return fmt.Errorf("conductor counldnt find node to remove %s", rmNodeName)
	}
	ctlNode, ok := c.CurrentNodes[c.pickRandomUpNode()]
	if !ok {
		return fmt.Errorf("no up nodes")
	}
	err := c.etcdctlRemoveNode(ctlNode, rmNode)
	if err != nil {
		return err
	}
	delete(c.CurrentNodes, rmNodeName)
	return nil
}

func (c *Conductor) removeExtraNodesFromCluster() error {
	for _, extraNode := range c.findMissingExtraNodes() {
		err := c.removeNodeFromCluster(extraNode)
		if err != nil {
			return fmt.Errorf("error removing %s: %s", extraNode, err)
		}
	}
	return nil
}
