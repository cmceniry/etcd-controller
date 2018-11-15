package driver

import (
	"context"
	"fmt"
	"net"

	pb "github.com/cmceniry/etcd-controller/driver/driverpb"
	"github.com/cmceniry/etcd-controller/ectl"
	"google.golang.org/grpc"
)

type SimpleDriver struct {
	Config     SimpleDriverConfig
	Process    *ectl.ETCDProcess
	Listener   net.Listener
	GRPCServer *grpc.Server
	inProgress bool
}

type SimpleDriverConfig struct {
	Binary      string
	Name        string
	DataDir     string
	IP          string
	ClientPort  int
	PeerPort    int
	CommandPort int
}

func NewSimpleDriver(c SimpleDriverConfig) (*SimpleDriver, error) {
	if c.Binary == "" {
		return nil, fmt.Errorf("Undefined binary")
	}
	if c.Name == "" {
		return nil, fmt.Errorf("Undefined name")
	}
	if c.DataDir == "" {
		return nil, fmt.Errorf("Undefined datadir")
	}
	if c.IP == "" {
		return nil, fmt.Errorf("Undefined IP")
	}
	if c.ClientPort == 0 {
		return nil, fmt.Errorf("Undefined ClientPort")
	}
	if c.PeerPort == 0 {
		return nil, fmt.Errorf("Undefined PeerPort")
	}
	if c.CommandPort == 0 {
		return nil, fmt.Errorf("Undefined CommandPort")
	}

	s := &SimpleDriver{Config: c}
	pc := ectl.ETCDConfig{
		Binary:                   s.Config.Binary,
		Name:                     s.Config.Name,
		DataDir:                  s.Config.DataDir,
		AdvertiseClientURLs:      fmt.Sprintf("http://%s:%d", s.Config.IP, s.Config.ClientPort),
		ListenClientURLs:         fmt.Sprintf("http://0.0.0.0:%d", s.Config.ClientPort),
		ClientCertAuth:           false,
		InitialAdvertisePeerURLs: fmt.Sprintf("http://%s:%d", s.Config.IP, s.Config.PeerPort),
		ListenPeerURLs:           fmt.Sprintf("http://0.0.0.0:%d", s.Config.PeerPort),
		PeerClientCertAuth:       false,
	}
	s.Process = &ectl.ETCDProcess{
		Config: pc,
	}
	return s, nil
}

func (s *SimpleDriver) runGRPCListener() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.CommandPort))
	if err != nil {
		return err
	}
	s.Listener = l
	var opts []grpc.ServerOption
	s.GRPCServer = grpc.NewServer(opts...)
	pb.RegisterDriverServer(s.GRPCServer, s)
	go func() {
		s.GRPCServer.Serve(s.Listener)
	}()
	return nil
}

func (s *SimpleDriver) Run() error {
	return s.runGRPCListener()
}

func (s *SimpleDriver) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	r := &pb.StatusResponse{}
	h, err := s.Process.GetHealth()
	if err != nil || !h {
		r.State = StateUnknown
	}
	r.State = StateRunning
	return r, nil
}

func (s *SimpleDriver) InitializeCluster(ctx context.Context, req *pb.InitClusterRequest) (*pb.InitClusterResponse, error) {
	r := &pb.InitClusterResponse{}
	if req.Snapshot != "" {
		r.ErrorMessage = "Init from snapshot not implemented"
	}
	err := s.Process.StartInitial()
	if err != nil {
		r.ErrorMessage = err.Error()
	}
	r.Success = r.ErrorMessage == ""
	return r, nil
}
