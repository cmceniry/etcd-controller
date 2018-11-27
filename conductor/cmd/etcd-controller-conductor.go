package main

import (
	"github.com/cmceniry/etcd-controller/conductor"
)

func main() {
	c := conductor.NewConductor(conductor.Config{
		NodeListFilename: "/config/node-list.yaml",
		CommandPort:      4270,
	})
	c.Run()
}
