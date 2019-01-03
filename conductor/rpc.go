package conductor

import (
	"context"
	"fmt"

	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
	"google.golang.org/grpc"
)

//go:generate protoc -I conductorpb --go_out=plugins=grpc:conductorpb conductor.proto

// GetInfo returns information about the conductor
func (c *Conductor) GetInfo(ctx context.Context, req *pb.GetInfoRequest) (*pb.GetInfoResponse, error) {
	resp := &pb.GetInfoResponse{
		IsConductor: c.IsRunning(),
	}
	return resp, nil
}

// GetStatus returns the condition of the entire cluster
func (c *Conductor) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	resp := &pb.GetStatusResponse{
		Nodes: []*pb.NodeInfo{},
	}
	for _, ni := range c.CurrentNodes {
		resp.Nodes = append(resp.Nodes, &pb.NodeInfo{
			Name:   ni.Name,
			Status: pb.NodeInfoStatus(ni.Status),
		})
	}
	return resp, nil
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

// RegisterWithGRPCServer handles the connection of this service with the
// CommandPort
func (c *Conductor) RegisterWithGRPCServer(g *grpc.Server) {
	pb.RegisterConductorServer(g, c)
}
