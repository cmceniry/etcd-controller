package conductor

import (
	"context"
	"fmt"

	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
	"google.golang.org/grpc"
)

// Client is the library abstraction for the Conductor GRPC client
type Client struct {
	IP          string
	CommandPort int
	Opts        []grpc.DialOption
	conn        *grpc.ClientConn
	client      pb.ConductorClient
}

// NewClient returns a Conductor GRPC Client
func NewClient(ip string, cp int, opts []grpc.DialOption) (*Client, error) {
	c := &Client{
		IP:          ip,
		CommandPort: cp,
		Opts:        opts,
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", c.IP, c.CommandPort),
		c.Opts...,
	)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.client = pb.NewConductorClient(c.conn)
	return c, nil
}

// ConductorInfo displays information about the Conductor
type ConductorInfo struct {
	IsConductor bool
}

// NodeStatus wraps the GRPC node status response
type NodeStatus struct {
	Name   string
	Status string
}

// ClusterStatus collect the NodeStatus of all nodes in the cluster
type ClusterStatus struct {
	Nodes []NodeStatus
}

// Info wraps GRPC GetInfo for the conductor
func (c *Client) Info() (ConductorInfo, error) {
	sr := &pb.GetInfoRequest{}
	r, err := c.client.GetInfo(context.Background(), sr)
	if err != nil {
		return ConductorInfo{}, fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	ci := ConductorInfo{
		IsConductor: r.IsConductor,
	}
	return ci, nil
}

// Status wraps GRPC GetStatus for the cluster
func (c *Client) Status() (ClusterStatus, error) {
	sr := &pb.GetStatusRequest{}
	r, err := c.client.GetStatus(context.Background(), sr)
	if err != nil {
		return ClusterStatus{}, fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	cs := ClusterStatus{
		Nodes: []NodeStatus{},
	}
	for _, n := range r.Nodes {
		cs.Nodes = append(cs.Nodes, NodeStatus{
			Name:   n.Name,
			Status: n.Status.String(),
		})
	}
	return cs, nil
}

// NodeStatus wraps the GRPC GetNodeStatus call
func (c *Client) NodeStatus(nodeName string) (NodeStatus, error) {
	sr := &pb.GetNodeStatusRequest{Name: nodeName}
	r, err := c.client.GetNodeStatus(context.Background(), sr)
	if err != nil {
		return NodeStatus{}, fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	return NodeStatus{
		Name:   r.Node.Name,
		Status: r.Node.Status.String(),
	}, nil
}
