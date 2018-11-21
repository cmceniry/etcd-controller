FROM golang:1.11 AS build

WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' && \
    CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' ctl/etcd-controller-ctl.go && \
    CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' conductor/cmd/etcd-controller-conductor.go


FROM quay.io/coreos/etcd:v3.1.19
WORKDIR /
COPY --from=build /src/etcd-controller /etcd-controller
COPY --from=build /src/etcd-controller-ctl /etcd-controller-ctl
COPY --from=build /src/etcd-controller-conductor /etcd-controller-conductor
