package driver

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/cmceniry/etcd-controller/driver/driverpb"
	"google.golang.org/grpc"
)

type SimpleClient struct {
	IP          string
	CommandPort int
	Opts        []grpc.DialOption
	conn        *grpc.ClientConn
	client      pb.DriverClient
}

func NewSimpleClient(ip string, cp int, opts []grpc.DialOption) (*SimpleClient, error) {
	c := &SimpleClient{
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
	// defer c.conn.Close()
	c.client = pb.NewDriverClient(c.conn)
	return c, nil
}

func (c *SimpleClient) Status() (int32, error) {
	sr := &pb.StatusRequest{}
	r, err := c.client.GetStatus(context.Background(), sr)
	if err != nil {
		return 0, fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	return r.State, nil
	// if r.State != StateRunning {
	// 	return 0, fmt.Errorf("%s:%d unhealthy", c.IP, c.CommandPort)
	// }
	// return 1, nil
}

func (c *SimpleClient) Stop() error {
	r, err := c.client.StopServer(context.Background(), &pb.StopServerRequest{})
	if err != nil {
		return fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	if !r.Success {
		return fmt.Errorf("%s:%d stop failure: %s", c.IP, c.CommandPort, r.ErrorMessage)
	}
	return nil
}

func (c *SimpleClient) InitCluster() error {
	icr := &pb.InitClusterRequest{}
	r, err := c.client.InitializeCluster(context.Background(), icr)
	if err != nil {
		return fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	if !r.Success {
		return fmt.Errorf("%s:%d init failure: %s", c.IP, c.CommandPort, err)
	}
	return nil
}

func (c *SimpleClient) JoinCluster(peers []string) error {
	peerInfos := []*pb.PeerInfo{}
	for _, pStr := range peers {
		pStrSplit := strings.Split(pStr, "=")
		if len(pStrSplit) != 2 {
			return fmt.Errorf("Invalid peer: %s", pStr)
		}
		peerInfos = append(peerInfos, &pb.PeerInfo{
			Name: pStrSplit[0],
			URL:  pStrSplit[1],
		})
	}
	jr := &pb.JoinClusterRequest{
		Peers: peerInfos,
	}
	r, err := c.client.JoinCluster(context.Background(), jr)
	if err != nil {
		return fmt.Errorf("%s:%d GRPC call failure: %s", c.IP, c.CommandPort, err)
	}
	if !r.Success {
		return fmt.Errorf("%s:%d join failure: %s", c.IP, c.CommandPort, r.ErrorMessage)
	}
	return nil
}
