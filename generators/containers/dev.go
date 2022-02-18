package containers

import (
	"fmt"
	"os"
	"runtime"

	docker "github.com/mu-box/golang-docker-client"
	boxfile "github.com/mu-box/microbox-boxfile"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/provider"
)

// DevConfig generate the container configuration for the build container
func DevConfig(appModel *models.App) docker.ContainerConfig {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))

	image := boxfile.Node("run.config").StringValue("image")

	if image == "" {
		image = "mubox/build"
	}

	code := fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), appModel.EnvID)

	if !provider.RequiresMount() {
		code = fmt.Sprintf("%s:/app", config.LocalDir())
	}

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("microbox_%s", appModel.ID),
		Image:   image, // this will need to be configurable some time
		Network: "virt",
		IP:      appModel.LocalIPs["env"],
		Binds: []string{
			code,
			// fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), appModel.EnvID),
			// fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), appModel.EnvID),
			fmt.Sprintf("microbox_%s_build:/data", appModel.EnvID),
			fmt.Sprintf("microbox_%s_cache:/mnt/cache", appModel.EnvID),
		},
		RestartPolicy: "no",
	}

	// set the terminal veriable
	if runtime.GOOS == "windows" {
		config.Env = []string{"TERM=cygwin"}
	}

	termEvar := os.Getenv("TERM")
	// msys doesnt work on linux so we will leave cygwin
	if termEvar != "" && termEvar != "msys" {
		config.Env = []string{"TERM=" + termEvar}
	}

	// set http[s]_proxy and no_proxy vars
	setProxyVars(&config)

	// // add cache_dirs into the container binds
	// libDirs := boxfile.Node("run.config").StringSliceValue("cache_dirs")

	// for _, libDir := range libDirs {
	// 	// TODO: the cache source should come from the provider
	// 	path := fmt.Sprintf("%s/%s/cache/cache_dirs/%s:/app/%s", provider.HostMntDir(), appModel.EnvID, libDir, libDir)
	// 	config.Binds = append(config.Binds, path)
	// }

	return config
}

// DevName returns the name of the build container
func DevName() string {
	return fmt.Sprintf("microbox_%s_dev", config.EnvID())
}
