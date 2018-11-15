package main

import (
	"os"
	"time"

	"github.com/cmceniry/etcd-controller/driver"
)

func main() {
	if os.Getenv("ETCDCONTROLLER_IP") == "" {
		panic("ETCDCONTROLLER_IP must be set")
	}
	s, err := driver.NewSimpleDriver(
		driver.SimpleDriverConfig{
			Binary:      "/usr/local/bin/etcd",
			Name:        "test01",
			DataDir:     "/var/lib/etcd",
			IP:          os.Getenv("ETCDCONTROLLER_IP"),
			ClientPort:  2379,
			PeerPort:    2380,
			CommandPort: 4270,
		},
	)
	if err != nil {
		panic(err)
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(1 * time.Second)
	for range t.C {
		time.Sleep(time.Second)
	}
}
