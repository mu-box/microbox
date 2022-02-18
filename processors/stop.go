package processors

import (
	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/processors/server"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	util_provider "github.com/mu-box/microbox/util/provider"
)

// Stop stops the running apps, unmounts all envs, and stops the provider
func Stop() error {
	// if the util provider isnt ready it doesnt need to stop
	if !util_provider.IsReady() {
		return nil
	}

	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// stop all running apps
	if err := stopAllApps(); err != nil {
		return util.ErrorAppend(err, "failed to stop running apps")
	}

	// env unmounting shouldnt be a problem any more
	// // unmount envs
	// if err := unmountEnvs(); err != nil {
	// 	return util.ErrorAppend(err, "failed to unmount envs")
	// }

	// stop the provider
	if err := provider.Stop(); err != nil {
		return util.ErrorAppend(err, "failed to stop the provider")
	}

	// stop the server
	if err := server.Stop(); err != nil {
		return util.ErrorAppend(err, "failed to stop server")
	}

	return nil
}

// stopAllApps stops all of the apps that are currently running
func stopAllApps() error {

	// load all the apps that think they're currently up
	apps, err := models.AllAppsByStatus("up")
	if err != nil {
		lumber.Error("stopAllApps:models.AllAppsByStatus(up)")
		return util.ErrorAppend(err, "failed to load running apps")
	}

	if len(apps) == 0 {
		return nil
	}

	display.OpenContext("Stopping Apps and Components")
	defer display.CloseContext()

	// run the app stop on all running apps
	for _, a := range apps {
		if err := app.Stop(a); err != nil {
			return util.ErrorAppend(err, "failed to stop running app")
		}
	}

	return nil
}

// unmountEnvs unmounts all of the environments
// func unmountEnvs() error {
// 	// unmount all the environments so stoping doesnt take forever
//
// 	envs, err := models.AllEnvs()
// 	if err != nil {
// 		display.ErrorTask()
// 		return util.ErrorAppend(err, "failed to load all envs")
// 	}
//
// 	if len(envs) == 0 {
// 		return nil
// 	}
//
// 	display.OpenContext("Removing mounts")
// 	defer display.CloseContext()
//
// 	for _, e := range envs {
// 		if err := env.Unmount(e); err != nil {
// 			display.ErrorTask()
// 			return util.ErrorAppend(err, "failed to unmount env")
// 		}
// 	}
//
// 	return nil
// }
