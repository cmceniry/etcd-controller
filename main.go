package main

import (
	"os"
	"time"

	"github.com/cmceniry/etcd-controller/driver"
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

	err = s.Run()
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(1 * time.Second)
	for range t.C {
		time.Sleep(time.Second)
	}
}
