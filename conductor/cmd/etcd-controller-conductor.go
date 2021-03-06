package main

import (
	"fmt"
	"net"

	"github.com/cmceniry/etcd-controller/conductor"
	pb "github.com/cmceniry/etcd-controller/conductor/conductorpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	c := conductor.NewConductor(conductor.Config{
		NodeListFilename: "/config/node-list.yaml",
		CommandPort:      4270,
		Logger:           log.WithFields(log.Fields{"component": "conductor"}),
	})

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 4270))
	if err != nil {
		panic(err)
	}
	var opts []grpc.ServerOption
	gserver := grpc.NewServer(opts...)
	pb.RegisterConductorServer(gserver, c)
	go func() {
		err := gserver.Serve(l)
		panic(err)
	}()

	c.Run()
}
