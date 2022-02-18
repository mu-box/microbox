package containers

import (
	"fmt"

	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
)

// ComponentConfig generates the container configuration for a component container
func ComponentConfig(componentModel *models.Component) docker.ContainerConfig {
	config := docker.ContainerConfig{
		Name:          ComponentName(componentModel),
		Image:         componentModel.Image,
		Network:       "virt",
		IP:            componentModel.IPAddr(),
		RestartPolicy: "no",
	}

	// set http[s]_proxy and no_proxy vars
	setProxyVars(&config)

	return config
}

// ComponentName returns the name of the component container
func ComponentName(componentModel *models.Component) string {
	return fmt.Sprintf("microbox_%s_%s", componentModel.AppID, componentModel.Name)
}
