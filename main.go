package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cmceniry/etcd-controller/conductor"
	"github.com/cmceniry/etcd-controller/driver"
	"github.com/cmceniry/etcd-controller/group"
	"google.golang.org/grpc"
)

func main() {
	v := map[string]string{
		"ETCDCONTROLLER_NAME":    "test001",
		"ETCDCONTROLLER_IP":      "127.0.0.1",
		"ETCDCONTROLLER_BINARY":  "/usr/local/bin/etcd",
		"ETCDCONTROLLER_DATADIR": "/var/lib/etcd",
	}
	for k := range v {
		if os.Getenv(k) != "" {
			v[k] = os.Getenv(k)
		}
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 4270))
	if err != nil {
		panic(err)
	}
	var opts []grpc.ServerOption
	gserver := grpc.NewServer(opts...)

	s, err := driver.NewSimpleDriver(
		driver.SimpleDriverConfig{
			Binary:      v["ETCDCONTROLLER_BINARY"],
			Name:        v["ETCDCONTROLLER_NAME"],
			DataDir:     v["ETCDCONTROLLER_DATADIR"],
			IP:          v["ETCDCONTROLLER_IP"],
			ClientPort:  2379,
			PeerPort:    2380,
			CommandPort: 4270,
		},
	)
	if err != nil {
		panic(err)
	}
	s.RegisterWithGRPCServer(gserver)

	c := conductor.NewConductor(conductor.Config{
		NodeListFilename: "/config/node-list.yaml",
		CommandPort:      4270,
	})
	c.RegisterWithGRPCServer(gserver)

	go func() {
		for {
			err := gserver.Serve(l)
			if err != nil {
				fmt.Printf("grpc serve fail: %s\n", err)
			}
		}
	}()

	m, err := group.NewManager(
		group.Config{
			Name:             v["ETCDCONTROLLER_NAME"],
			IP:               v["ETCDCONTROLLER_IP"],
			SerfPort:         4271,
			NodeListFilename: "/config/node-list.yaml",
		},
	)
	if err != nil {
		panic(err)
	}
	m.Run()

	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-t.C:
			fmt.Printf("main TICK!\n")
			if isCon, notCon := m.IsConductor(); isCon {
				fmt.Printf("IS CONDUCTOR\n")
				if !c.IsRunning() {
					go c.Run()
				}
			} else {
				fmt.Printf("NOT CONDUCTOR: %s\n", notCon)
				if c != nil {
					fmt.Printf("TODO: Should stop conductor\n")
				}
			}
		}
	}
}
