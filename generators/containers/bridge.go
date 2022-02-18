package containers

import (
	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/util/dhcp"
)

// BridgeConfig generates the container configuration for a component container
func BridgeConfig() docker.ContainerConfig {
	return docker.ContainerConfig{
		Name:          BridgeName(),
		Image:         "mubox/bridge",
		Network:       "virt",
		IP:            reserveIP(),
		RestartPolicy: "always",
		Ports:         []string{"1194:1194/udp"},
	}
}

// BridgeName returns the name of the component container
func BridgeName() string {
	return "microbox_bridge"
}

// reserveIP reserves a local IP for the build container
func reserveIP() string {
	ip, _ := dhcp.ReserveLocal()
	return ip.String()
}
