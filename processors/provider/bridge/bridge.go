package bridge

import (
	"fmt"
	"path/filepath"

	// "runtime"
	"net"

	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/dhcp"
	"github.com/mu-box/microbox/util/provider"
)

func BridgeConfig() string {
	// node := ""
	// if runtime.GOOS == "windows" {
	// 	node = "dev-node MyTap"
	// }

	ip, _ := provider.HostIP()
	return fmt.Sprintf(`client

dev tap
proto udp
remote %s 1194
resolv-retry infinite
nobind
persist-key
persist-tun

ca "%s"
cert "%s"
key "%s"

cipher none
auth none
verb 3
`, ip, CaCrt(), ClientCrt(), ClientKey())
}

func ConfigFile() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "openvpn.conf"))
}

func CaCrt() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "ca.crt"))
}

func ClientKey() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "client.key"))
}

func ClientCrt() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "client.crt"))
}

// check to see if the bridge is connected
func Connected() bool {
	network, err := dhcp.LocalNet()
	if err != nil {
		return false
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	// look through the interfaces on the system
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		// find all the addresses assigned to the interface
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			// check to see if the ip address is in our network
			if network.Contains(ip) {

				// now check to see if that interface is up
				if i.Flags&net.FlagUp != net.FlagUp {
					return false
				}

				return true
			}
		}
	}

	return false
}
