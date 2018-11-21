package main

import (
	"io/ioutil"

	"github.com/cmceniry/etcd-controller/conductor"
)

func main() {
	c := conductor.NewConductor()
	d, err := ioutil.ReadFile("/config/node-list.yaml")
	if err != nil {
		panic(err)
	}
	err = c.LoadYaml(d)
	if err != nil {
		panic(err)
	}
	c.Run()
}
