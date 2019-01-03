package group

import (
	"testing"
)

func TestNetID(t *testing.T) {
	a := netID("0.0.0.0", 0)
	if a != 0 {
		t.Errorf("0.0.0.0:0 != 0. is %#x", uint64(a))
	}
	a = netID("127.0.0.1", 123)
	if a != 0x7f000001007b {
		t.Errorf("128.0.0.0:123 != 0x7f000001007b. is %#x", a)
	}
}
