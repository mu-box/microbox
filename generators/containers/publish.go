package containers

import (
	"fmt"

	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/config"
	// "github.com/mu-box/microbox/util/provider"
)

// PublishConfig generate the container configuration for the build container
func PublishConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	config := docker.ContainerConfig{
		Name:    PublishName(),
		Image:   image,
		Network: "host",
		Binds: []string{
			// fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/deploy:/mnt/deploy", provider.HostMntDir(), env),
			fmt.Sprintf("microbox_%s_app:/mnt/app", env),
			fmt.Sprintf("microbox_%s_cache:/mnt/cache", env),
			fmt.Sprintf("microbox_%s_deploy:/mnt/deploy", env),
		},
		RestartPolicy: "no",
	}

	// Some CI's have an old kernel and require us to use the virtual network
	// this is only in effect for CI's because it automatically reserves an ip on our microbox
	// virtual network and we could have IP conflicts
	configModel, _ := models.LoadConfig()
	if configModel.CIMode {
		config.Network = "virt"
	}

	// set http[s]_proxy and no_proxy vars
	setProxyVars(&config)

	return config
}

// PublishName returns the name of the build container
func PublishName() string {
	return fmt.Sprintf("microbox_%s_publish", config.EnvID())
}
