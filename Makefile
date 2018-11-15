


docker:
	go mod tidy
	docker build -t etcd-controller:snapshot .

explore:
	docker run -it --rm etcd-controller:snapshot /bin/sh