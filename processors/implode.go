package processors

import (
	"fmt"
	"os"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/env"
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/processors/server"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	util_provider "github.com/mu-box/microbox/util/provider"
)

// Implode destroys the provider and cleans microbox off of the system
func Implode() error {

	display.OpenContext("Imploding Microbox")
	defer display.CloseContext()

	// remove all environments
	envModels, _ := models.AllEnvs()
	for _, envModel := range envModels {
		// remove all environments
		if err := env.Destroy(envModel); err != nil {
			fmt.Printf("unable to remove mounts: %s", err)
		}
	}

	// destroy the provider
	if err := provider.Destroy(); err != nil {
		return util.ErrorAppend(err, "failed to destroy the provider")
	}

	// destroy the provider (VM), remove images, remove containers
	if err := util_provider.Implode(); err != nil {
		return util.ErrorAppend(err, "failed to implode the provider")
	}

	// check to see if we need to uninstall microbox
	// or just remove apps
	if registry.GetBool("full-implode") {

		// teardown the server
		//lint:ignore SA9003 if we cant tear down the server dont worry about it
		if err := server.Teardown(); err != nil {
			// return util.ErrorAppend(err, "failed to remove server")
		}

		purgeConfiguration()
	}

	return nil
}

// purges the config data and dns entries
func purgeConfiguration() error {

	display.StartTask("Purging configuration")
	defer display.StopTask()

	// implode the global dir
	if err := os.RemoveAll(config.GlobalDir()); err != nil {
		return util.ErrorAppend(util.ErrorQuiet(err), "failed to purge the data directory")
	}

	return nil
}
