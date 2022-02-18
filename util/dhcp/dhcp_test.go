package dhcp_test

import (
	"net"
	"os"
	"testing"

	"github.com/mu-box/microbox/util/dhcp"
)

// TestMain ...
func TestMain(m *testing.M) {
	dhcp.Flush()
	os.Exit(m.Run())
}

// TestReservingIps ...
func TestReservingIps(t *testing.T) {
	ipOne, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	ipTwo, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	ipThree, err := dhcp.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	if ipOne.String() != "192.168.99.51" || ipTwo.String() != "192.168.99.52" || (ipThree.String() != "172.20.0.2" && ipThree.String() != "172.21.0.2") {
		t.Errorf("incorrect ip addresses %s / %s / %s", ipOne, ipTwo, ipThree)
	}
}

// TestReturnIP ...
func TestReturnIP(t *testing.T) {
	err := dhcp.ReturnIP(net.ParseIP("192.168.99.50"))
	if err != nil {
		t.Errorf("unable to return ip %s", err)
	}
	err = dhcp.ReturnIP(net.ParseIP("192.168.99.51"))
	if err != nil {
		t.Errorf("unable to return ip %s", err)
	}
	err = dhcp.ReturnIP(net.ParseIP("192.168.0.50"))
	if err != nil {
		t.Errorf("unable to return ip %s", err)
	}
}

// TestReuseIP ...
func TestReuseIP(t *testing.T) {
	one, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	ipTwo, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	three, err := dhcp.ReserveLocal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	err = dhcp.ReturnIP(ipTwo)
	if err != nil {
		t.Errorf("unable to return ip %s", err)
	}
	ipTwoAgain, err := dhcp.ReserveGlobal()
	if err != nil {
		t.Errorf("unable to reserve ip %s", err)
	}
	if !ipTwo.Equal(ipTwoAgain) {
		t.Errorf("should have received a repeat of %s but got %s", ipTwo.String(), ipTwoAgain.String())
	}
	dhcp.ReturnIP(one)
	dhcp.ReturnIP(three)
}
