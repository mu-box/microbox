package env

import (
	"fmt"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/locker"
	util_provider "github.com/mu-box/microbox/util/provider"
)

// Destroy brings down the environment setup
func Destroy(env *models.Env) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// find apps
	apps, err := env.Apps()
	if err != nil {
		lumber.Error("env:Destroy:models.Env{ID:%s}.Apps(): %s", env.ID, err)
		return util.ErrorAppend(err, "failed to load app collection")
	}

	// destroy apps
	for _, a := range apps {

		err := app.Destroy(a)
		if err != nil {
			return util.ErrorAppend(err, "failed to remove app")
		}
	}

	// unmount the environment
	if err := Unmount(env); err != nil {
		return util.ErrorAppend(err, "failed to unmount env")
	}

	// TODO: remove folder from host /mnt/sda1/env_id
	if err := util_provider.RemoveEnvDir(env.ID); err != nil {
		// it is ok if the cleanup fails its not worth erroring here
		// return util.ErrorAppend(err, "failed to remove the environment from host")
	}

	// remove volumes
	docker.VolumeRemove(fmt.Sprintf("microbox_%s_app", env.ID))
	docker.VolumeRemove(fmt.Sprintf("microbox_%s_cache", env.ID))
	docker.VolumeRemove(fmt.Sprintf("microbox_%s_mount", env.ID))
	docker.VolumeRemove(fmt.Sprintf("microbox_%s_deploy", env.ID))
	docker.VolumeRemove(fmt.Sprintf("microbox_%s_build", env.ID))

	// remove the environment
	if err := env.Delete(); err != nil {
		return util.ErrorAppend(err, "failed to remove env")
	}

	return nil
}
