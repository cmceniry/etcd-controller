package conductor

import (
	"context"
	"fmt"
	"sort"
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
	eNoRun := []string{}
	eRun := []string{}
	for nn, ni := range c.CurrentNodes {
		if !c.IsExtra(nn) {
			continue
		}
		if ni.IsRunning() {
			eRun = append(eRun, nn)
		} else {
			eNoRun = append(eNoRun, nn)
		}
	}
	sort.Strings(eNoRun)
	sort.Strings(eRun)
	// List the not running ones as extra first
	return append(eNoRun, eRun...)
}

// removeNodeFromCluster attempts to update two pieces: the etcd cluster, and
// the etcd-controller group.
func (c *Conductor) removeNodeFromCluster(rmNodeName string) error {
	rmNode, ok := c.CurrentNodes[rmNodeName]
	if !ok {
		return fmt.Errorf("conductor counldn't find node to remove %s", rmNodeName)
	}
	ctlNode, ok := c.CurrentNodes[c.pickRandomUpNode()]
	if !ok {
		return fmt.Errorf("no up nodes")
	}
	// TODO: don't remove if it'll upset etcd cluster availability - up nodes < quorum
	//  - I think it should only be an issue if there are stopped/failed nodes and
	//    it's attempting to remove an up node
	err := c.etcdctlRemoveNode(ctlNode, rmNode)
	if err != nil {
		return err
	}

	if rmNode.Status == 0 || rmNode.Status == 2 || rmNode.Status == 3 {
		delete(c.CurrentNodes, rmNodeName)
	}

	return nil
}

func (c *Conductor) removeExtraNodesFromCluster() error {
	for _, extraNode := range c.findMissingExtraNodes() {
		err := c.removeNodeFromCluster(extraNode)
		if err != nil {
			return fmt.Errorf("error removing %s: %s", extraNode, err)
		}
		fmt.Printf("removed %s\n", extraNode)
	}
	return nil
}
