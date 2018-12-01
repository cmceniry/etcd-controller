package conductor

import (
	"context"
	"fmt"
	"net"

	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
	"google.golang.org/grpc"
)

//go:generate protoc -I conductorpb --go_out=plugins=grpc:conductorpb conductor.proto

// GetStatus returns the condition of the entire cluster
func (c *Conductor) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetNodeStatus returns the condition of a specific node
func (c *Conductor) GetNodeStatus(ctx context.Context, req *pb.GetNodeStatusRequest) (*pb.GetNodeStatusResponse, error) {
	ni, ok := c.CurrentNodes[req.Name]
	if !ok {
		return nil, fmt.Errorf("node %s not known", req.Name)
	}
	resp := &pb.GetNodeStatusResponse{
		Node: &pb.NodeInfo{
			Name:   ni.Name,
			Url:    ni.PeerURL(),
			Status: pb.NodeInfoStatus(ni.Status),
		},
	}
	return resp, nil
}

func (c *Conductor) runGRPCListener() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Config.CommandPort))
	if err != nil {
		return err
	}
	c.Listener = l
	var opts []grpc.ServerOption
	c.GRPCServer = grpc.NewServer(opts...)
	pb.RegisterConductorServer(c.GRPCServer, c)
	go func() {
		c.GRPCServer.Serve(c.Listener)
	}()
	return nil
}
