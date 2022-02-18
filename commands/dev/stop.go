package dev

import (
	docker "github.com/mu-box/golang-docker-client"
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	container_generator "github.com/mu-box/microbox/generators/containers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

func init() {
	steps.Build("dev stop", stopCheck, stopFn)
}

//
// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	// TODO: check the app and return some message
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	display.CommandErr(app.Stop(appModel))
}

func stopCheck() bool {
	container, err := docker.GetContainer(container_generator.DevName())

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
