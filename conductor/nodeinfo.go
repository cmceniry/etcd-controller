package conductor

import (
	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
)

// IsRunning indicates that etcd for that node is running
func (n NodeInfo) IsRunning() bool {
	return n.Status == int32(pb.NodeInfoStatus_RUNNING)
}

// IsStopped indicates that etcd for that node is stopped (but looks like it
// could start normally)
func (n NodeInfo) IsStopped() bool {
	return n.Status == int32(pb.NodeInfoStatus_STOPPED)
}
