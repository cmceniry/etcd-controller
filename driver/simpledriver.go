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
	peerTLS    bool
	clientTLS  bool
	failed     bool
}

type SimpleDriverConfig struct {
	Binary      string
	Name        string
	DataDir     string
	IP          string
	ClientPort  int
	PeerPort    int
	CommandPort int

	PeerTLSCA   string
	PeerTLSCert string
	PeerTLSKey  string

	ClientTLSCA   string
	ClientTLSCert string
	ClientTLSKey  string
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

	if c.PeerTLSCA != "" && c.PeerTLSCert != "" && c.PeerTLSKey != "" {
		s.peerTLS = true
	}
	if c.ClientTLSCA != "" && c.ClientTLSCert != "" && c.ClientTLSKey != "" {
		s.clientTLS = true
	}

	peerProto := "http"
	if s.peerTLS {
		peerProto = "https"
	}
	clientProto := "http"
	if s.clientTLS {
		clientProto = "https"
	}

	pc := ectl.ETCDConfig{
		Binary:                   s.Config.Binary,
		Name:                     s.Config.Name,
		DataDir:                  s.Config.DataDir,
		AdvertiseClientURLs:      fmt.Sprintf("%s://%s:%d", clientProto, s.Config.IP, s.Config.ClientPort),
		ListenClientURLs:         fmt.Sprintf("%s://0.0.0.0:%d", clientProto, s.Config.ClientPort),
		InitialAdvertisePeerURLs: fmt.Sprintf("%s://%s:%d", peerProto, s.Config.IP, s.Config.PeerPort),
		ListenPeerURLs:           fmt.Sprintf("%s://0.0.0.0:%d", peerProto, s.Config.PeerPort),
	}
	if s.peerTLS {
		pc.PeerClientCertAuth = true
		pc.PeerCAFile = c.PeerTLSCA
		pc.PeerCertFile = c.PeerTLSCert
		pc.PeerKeyFile = c.PeerTLSKey
	}
	if s.clientTLS {
		pc.ClientCertAuth = true
		pc.CAFile = c.ClientTLSCA
		pc.CertFile = c.ClientTLSCert
		pc.KeyFile = c.ClientTLSKey
	}
	s.Process = &ectl.ETCDProcess{
		Config: pc,
	}
	return s, nil
}

// RegisterWithGRPCServer handles the connection of this service with the
// CommandPort
func (s *SimpleDriver) RegisterWithGRPCServer(g *grpc.Server) {
	pb.RegisterDriverServer(g, s)
}

func (s *SimpleDriver) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	r := &pb.StatusResponse{}
	if s.Process == nil || !s.Process.IsInitialized() {
		r.State = StateUnknown
		return r, nil
	}
	if !s.Process.IsRunning() {
		r.State = StateStopped
		return r, nil
	}
	h, err := s.Process.GetHealth()
	if err != nil || !h {
		fmt.Printf("Health failed: %t - err: %s\n", h, err)
		r.State = StateUnknown
		return r, nil
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

func (s *SimpleDriver) JoinCluster(ctx context.Context, req *pb.JoinClusterRequest) (*pb.JoinClusterResponse, error) {
	r := &pb.JoinClusterResponse{}
	if len(req.Peers) == 0 {
		r.ErrorMessage = "Must specify peer"
		return r, nil
	}
	peerURLs := map[string]string{}
	for _, p := range req.Peers {
		if p.Name == "" {
			r.ErrorMessage = "Missing peer name"
			return r, nil
		}
		if p.URL == "" {
			r.ErrorMessage = "missing peer URL"
			return r, nil
		}
		peerURLs[p.Name] = p.URL
	}
	success, err := s.Process.JoinCluster(peerURLs)
	if err != nil {
		r.ErrorMessage = err.Error()
		return r, nil
	}
	if !success {
		r.ErrorMessage = "Unsuccessful"
	}
	if success && r.ErrorMessage == "" {
		r.Success = true
	}
	return r, nil
}

// StartServer is the RPC call receiver which will start an ETCD server under
// normal startup circumstances (i.e. no join/init/etc).
func (s *SimpleDriver) StartServer(ctx context.Context, req *pb.StartServerRequest) (*pb.StartServerResponse, error) {
	r := &pb.StartServerResponse{}
	err := s.Process.StartServer()
	if err != nil {
		r.ErrorMessage = err.Error()
	} else {
		r.Success = true
	}
	return r, nil
}

func (s *SimpleDriver) StopServer(ctx context.Context, req *pb.StopServerRequest) (*pb.StopServerResponse, error) {
	r := &pb.StopServerResponse{}
	err := s.Process.StopServer()
	if err != nil {
		r.ErrorMessage = err.Error()
	} else {
		r.Success = true
	}
	return r, nil
}
