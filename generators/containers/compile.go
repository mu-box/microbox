package containers

import (
	"fmt"

	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/provider"
)

// CompileConfig generate the container configuration for the build container
func CompileConfig(image string) docker.ContainerConfig {
	env := config.EnvID()
	code := fmt.Sprintf("%s%s/code:/share/code", provider.HostShareDir(), env)
	engine := fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), env)

	if !provider.RequiresMount() {
		code = fmt.Sprintf("%s:/share/code", config.LocalDir())
		engineDir, _ := config.EngineDir()
		if engineDir != "" {
			engine = fmt.Sprintf("%s:/share/engine", engineDir)
		}
	}

	conf := docker.ContainerConfig{
		Name:    CompileName(),
		Image:   image,
		Network: "host",
		Binds: []string{
			code,
			engine,
			// fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), env),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), env),
			fmt.Sprintf("microbox_%s_build:/data", env),
			fmt.Sprintf("microbox_%s_app:/mnt/app", env),
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

// CompileName returns the name of the build container
func CompileName() string {
	return fmt.Sprintf("microbox_%s_compile", config.EnvID())
}
