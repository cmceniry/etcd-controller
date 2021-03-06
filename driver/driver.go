package driver

//go:generate protoc -I driverpb --go_out=plugins=grpc:driverpb driver.proto

// Driver is the generic interface for controlling a local process. In reality
// a Driver must also implement the driverpb.DriverServer interface for it to
// fit into the scheme of everything
type Driver interface {
	Run() error
}

const (
	// State represents current status

	// StateUnknown is the starting point
	StateUnknown = iota

	// StateRunning means running
	StateRunning

	// StateStopped means ready but not running
	StateStopped

	// StateFailed means not ready (and not running)
	StateFailed
)
