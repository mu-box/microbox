package env

import (
	"path/filepath"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/provider"
)

// Unmount unmounts the env shares
func Unmount(env *models.Env) error {

	// break early if there is still an environment using the mounts
	if mountsInUse(env) {
		return nil
	}

	display.StartTask(env.Name)
	defer display.StopTask()

	// unmount the engine if it's a local directory
	engineDir, _ := config.EngineDir()
	if engineDir != "" {
		src := engineDir                                                // local directory
		dst := filepath.Join(provider.HostShareDir(), env.ID, "engine") // b2d "global zone"

		// unmount the env on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			display.ErrorTask()
			return util.ErrorAppend(err, "failed to remove engine mount")
		}
	}

	// unmount the app src
	src := env.Directory
	dst := filepath.Join(provider.HostShareDir(), env.ID, "code")

	// unmount the env on the provider
	if err := provider.RemoveMount(src, dst); err != nil {
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to remove code mount")
	}

	return nil
}

func UnmountEngine(env *models.Env, engineDir string) error {
	// unmount the engine if it's a local directory
	if engineDir != "" {
		src := engineDir                                                // local directory
		dst := filepath.Join(provider.HostShareDir(), env.ID, "engine") // b2d "global zone"

		// unmount the env on the provider
		if err := provider.RemoveMount(src, dst); err != nil {
			display.ErrorTask()
			return util.ErrorAppend(err, "failed to cleanup old engine mount")
		}
	}
	return nil
}

// mountsInUse returns true if any of the env's apps are running
func mountsInUse(env *models.Env) bool {
	devApp, _ := models.FindAppBySlug(env.ID, "dev")
	simApp, _ := models.FindAppBySlug(env.ID, "sim")
	return devApp.Status == "up" || simApp.Status == "up"
}

// returns true if the app or engine is mounted
func mounted(env *models.Env) bool {

	// if the engine is mounted, check that
	engineDir, _ := config.EngineDir()
	if engineDir != "" {
		dst := filepath.Join(provider.HostShareDir(), env.ID, "engine")

		if provider.HasMount(dst) {
			return true
		}
	}

	// check to see if the code is mounted
	dst := filepath.Join(provider.HostShareDir(), env.ID, "code")
	return provider.HasMount(dst)
}
