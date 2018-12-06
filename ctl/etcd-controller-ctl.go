package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"strconv"

	"github.com/cmceniry/etcd-controller/driver"
	"github.com/cmceniry/etcd-controller/conductor"
	"google.golang.org/grpc"
)

func fail(rc int, message string, args ...interface{}) {
	fmt.Printf(message, args...)
	os.Exit(rc)
}

func mustSimpleClient(ip string, port int, opts []grpc.DialOption) *driver.SimpleClient {
	s, err := driver.NewSimpleClient(ip, port, opts)
	if err != nil {
		fail(-1, "%s:%d Simple client connect failure: %s\n", ip, port, err)
	}
	return s
}

func mustConductorClient(ip string, port int, opts []grpc.DialOption) *conductor.Client {
	c, err := conductor.NewClient(ip, port, opts)
	if err != nil {
		fail(-1, "%s:%d Conductor client connect failure: %s\n", ip, port, err)
	}
	return c
}



func main() {
	if len(os.Args) < 3 {
		fail(-1, "Usage: %s node action [args]\n", os.Args[0])
	}
	node := os.Args[1]
	action := os.Args[2]

	nodeSplit := strings.Split(node, ":")
	if len(nodeSplit) != 2 {
		fail(-1, `Invalid node format "%s": not name:port`, node)
	}
	nodeIP := nodeSplit[0]
	nodePort, err := strconv.Atoi(nodeSplit[1])
	if err != nil {
		fail(-1, `Invalid node format "%s": %s`, node, err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	switch action {
	case "init":
		err := mustSimpleClient(nodeIP, nodePort, opts).InitCluster()
		if err != nil {
			fail(-1, "%s init failure: %s\n", node, err)
		}
	case "status":
		status, err := mustSimpleClient(nodeIP, nodePort, opts).Status()
		if err != nil {
			fail(-1, "%s status failure: %s\n", node, err)
		}
		fmt.Printf("%s Status: %d\n", node, status)
	case "cstatus":
		status, err := mustConductorClient(nodeIP, nodePort, opts).Status()
		if err != nil {
			fail(-1, "%s status failure: %s\n", node, err)
		}
		sortNames := []string{}
		nodeStatuses := map[string]conductor.NodeStatus{}
		for _, n := range status.Nodes {
			sortNames = append(sortNames, n.Name)
			nodeStatuses[n.Name] = n
		}
		sort.Strings(sortNames)
		for _, nodeName := range sortNames {
			n := nodeStatuses[nodeName]
			fmt.Printf("%s %s\n", n.Name, n.Status)
		}
	case "nodestatus":
		if len(os.Args) != 4 {
			fail(-1, "Usage: %s node nodestatus nodeForStatus\n", os.Args[0])
		}
		nodestatus, err := mustConductorClient(nodeIP, nodePort, opts).NodeStatus(os.Args[3])
		if err != nil {
			fail(-1, "%s node status %s failure: %s\n", node, os.Args[3], err)
		}
		fmt.Printf("%s", nodestatus.Status)
	case "join":
		if len(os.Args) != 4 {
			fail(-1, "Usage: %s node join peer\n", os.Args[0])
		}
		peers := strings.Split(os.Args[3], ",")
		err := mustSimpleClient(nodeIP, nodePort, opts).JoinCluster(peers)
		if err != nil {
			fail(-1, "%s join failure: %s\n", node, err)
		}
	case "start":
		err := mustSimpleClient(nodeIP, nodePort, opts).Start()
		if err != nil {
			fail(-1, "%s start failure: %s\n", node, err)
		}
	case "stop":
		err := mustSimpleClient(nodeIP, nodePort, opts).Stop()
		if err != nil {
			fail(-1, "%s stop failure: %s\n", node, err)
		}
	default:
		fail(-1, "Unknown action: %s\n", action)
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
