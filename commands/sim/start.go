package sim

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
	steps.Build("sim start", startCheck, simStart)
}

// simStart ...
func simStart(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "sim")

	display.CommandErr(env.Setup(envModel))
	display.CommandErr(app.Start(envModel, appModel, "sim"))
}

func startCheck() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	if app.Status != "up" {
		return false
	}

	// make sure im mounted and ready to go
	envModel, _ := models.FindEnvByID(config.EnvID())
	if !util_provider.HasMount(fmt.Sprintf("%s%s/code", util_provider.HostShareDir(), envModel.ID)) {
		return false
	}

	provider.Init()
	components, _ := app.Components()
	for _, component := range components {
		info, err := docker.ContainerInspect(component.ID)
		if err != nil || !info.State.Running {
			return false
		}
	}
	return true
}
