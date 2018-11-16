package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cmceniry/etcd-controller/driver"
	pb "github.com/cmceniry/etcd-controller/driver/driverpb"
	"google.golang.org/grpc"
)

func fail(rc int, message string, args ...interface{}) {
	fmt.Printf(message, args...)
	os.Exit(rc)
}

func main() {
	if len(os.Args) < 3 {
		fail(-1, "Usage: %s node action [args]\n", os.Args[0])
	}
	node := os.Args[1]
	action := os.Args[2]

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(node, opts...)
	if err != nil {
		fail(-1, "%s GRPC dial failure: %s\n", node, err)
	}
	defer conn.Close()
	client := pb.NewDriverClient(conn)

	switch action {
	case "init":
		icr := &pb.InitClusterRequest{}
		r, err := client.InitializeCluster(context.Background(), icr)
		if err != nil {
			fail(-1, "%s GRPC call failure: %s\n", node, err)
		}
		if !r.Success {
			fail(-1, "%s init failure: %s\n", node, r.ErrorMessage)
		}
	case "status":
		sr := &pb.StatusRequest{}
		r, err := client.GetStatus(context.Background(), sr)
		if err != nil {
			fail(-1, "%s GRPC call failure: %s\n", node, err)
		}
		if r.State != driver.StateRunning {
			fail(-1, "%s unhealthy", node)
		}
	case "join":
		if len(os.Args) != 4 {
			fail(-1, "Usage: %s node join peer\n", os.Args[0])
		}
		pL := []*pb.PeerInfo{}
		for _, pS := range strings.Split(os.Args[3], ",") {
			pI := strings.Split(pS, "=")
			if len(pI) != 2 {
				fail(-1, "Invalid format %s", pS)
			}
			pL = append(pL, &pb.PeerInfo{
				Name: pI[0],
				URL: pI[1],
			})
		}
		jr := &pb.JoinClusterRequest{
			Peers: pL,
		}
		r, err := client.JoinCluster(context.Background(), jr)
		if err != nil {
			fail(-1, "%s GRPC call failure: %s\n", node, err)
		}
		if !r.Success {
			fail(-1, "%s join failure: %s\n", node, r.ErrorMessage)
		}
	case "stop":
		r, err := client.StopServer(context.Background(), &pb.StopServerRequest{})
		if err != nil {
			fail(-1, "%s GRPC call failure: %s\n", node, err)
		}
		if !r.Success {
			fail(-1, "%s stop failure: %s\n", node, r.ErrorMessage)
		}
	default:
		fail(-1, "Unknown action: %s", action)
	}

	os.Exit(0)
	// cli, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
	// 	DialTimeout: 5 * time.Second,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer cli.Close()
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// resp, err := cli.Put(ctx, "sample_key", "sample_value")
	// cancel()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%#v\n", resp)

}
