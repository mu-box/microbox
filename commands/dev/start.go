package dev

import (
	"fmt"

	docker "github.com/mu-box/golang-docker-client"
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/processors/env"
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	util_provider "github.com/mu-box/microbox/util/provider"
)

func init() {
	steps.Build("dev start", startCheck, devStart)
}

// devStart ...
func devStart(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	display.CommandErr(env.Setup(envModel))
	display.CommandErr(app.Start(envModel, appModel, "dev"))
}

func startCheck() bool {
	// check to see if the app is up

	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	if app.Status != "up" {
		return false
	}

	// make sure im mounted and ready to go
	envModel, _ := models.FindEnvByID(config.EnvID())
	if !util_provider.HasMount(fmt.Sprintf("%s%s/code", util_provider.HostShareDir(), envModel.ID)) {
		return false
	}

	// connect to docker and find out if the components exist and are running
	provider.Init()
	components, _ := app.Components()
	for _, component := range components {
		info, err := docker.ContainerInspect(component.ID)
		if err != nil || !info.State.Running {
			return false
		}
	}

	// if we found nothing wrong we are already started :)
	return true
}
