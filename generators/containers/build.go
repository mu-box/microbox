package containers

import (
	"fmt"

	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/provider"
)

// BuildConfig generate the container configuration for the build container
func BuildConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	code := fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), env)
	// mounting from b2d "global zone" to container will be the same whether local engine specified or not
	engine := fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), env)

	if !provider.RequiresMount() {
		code = fmt.Sprintf("%s:/app", config.LocalDir())

		// todo: test this (likely docker-native linux)
		engineDir, _ := config.EngineDir()
		if engineDir != "" {
			engine = fmt.Sprintf("%s:/share/engine", engineDir)
		}
	}

	conf := docker.ContainerConfig{
		Name:    BuildName(),
		Image:   image,
		Network: "host",
		Binds: []string{
			code,
			engine,
			// fmt.Sprintf("%s%s/build:/mnt/build", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/deploy:/mnt/deploy", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
			fmt.Sprintf("microbox_%s_build:/mnt/build", env),
			fmt.Sprintf("microbox_%s_deploy:/mnt/deploy", env),
			fmt.Sprintf("microbox_%s_cache:/mnt/cache", env),
		},
		RestartPolicy: "no",
	}

	// Some CI's have an old kernel and require us to use the virtual network
	// this is only in effect for CI's because it automatically reserves an ip on our microbox
	// virtual network and we could have IP conflicts
	configModel, _ := models.LoadConfig()
	if configModel.CIMode {
		conf.Network = "virt"
	}

	// set http[s]_proxy and no_proxy vars
	setProxyVars(&conf)

	return conf
}

// BuildName returns the name of the build container
func BuildName() string {
	return fmt.Sprintf("microbox_%s_build", config.EnvID())
}
