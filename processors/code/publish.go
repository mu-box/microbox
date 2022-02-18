package code

import (
	"strings"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	container_generator "github.com/mu-box/microbox/generators/containers"
	"github.com/mu-box/microbox/generators/hooks/build"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/hookit"
)

// Publish ...
func Publish(envModel *models.Env, WarehouseConfig WarehouseConfig) error {
	display.OpenContext("Deploying app")
	defer display.CloseContext()

	// initialize the docker client
	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return util.ErrorAppend(err, "failed to pull the build image")
	}

	display.StartTask("Starting docker container")

	// if a publish container was leftover from a previous publish, let's remove it
	docker.ContainerRemove(container_generator.PublishName())

	// start the container
	config := container_generator.PublishConfig(buildImage)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Build:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to start docker container")
	}
	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container.ID)

	display.StopTask()

	lumber.Prefix("code:Publish")
	defer lumber.Prefix("")

	display.StartTask("Uploading")

	// run user hook
	// TODO: display output from hooks
	payload := build.UserPayload()
	// todo: should this be if payload == ""
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload()")
		return util.ErrorAppend(err, "unable to retrieve user payload")
	}
	if out, err1 := hookit.DebugExec(container.ID, "user", payload, "info"); err1 != nil {
		// handle 'exec failed: argument list too long' error
		if strings.Contains(out, "argument list too long") {
			if err2, ok := err1.(util.Err); ok {
				err2.Suggest = "You may have too many ssh keys, please specify the one you need with `microbox config set ssh-key ~/.ssh/id_rsa`"
				err2.Output = out
				err2.Code = "1001"
				return util.ErrorAppend(err2, "failed to run the (publish)user hook")
			}
		}
		return util.ErrorAppend(err, "failed to run the (publish)user hook")
	}

	buildWarehouseConfig := build.WarehouseConfig{
		BuildID:        WarehouseConfig.BuildID,
		WarehouseURL:   WarehouseConfig.WarehouseURL,
		WarehouseToken: WarehouseConfig.WarehouseToken,
		PreviousBuild:  WarehouseConfig.PreviousBuild,
	}

	payload = build.PublishPayload(envModel, buildWarehouseConfig)
	if err != nil {
		lumber.Error("code:Publish:build.UserPayload()")
		display.ErrorTask()
		return util.ErrorAppend(err, "unable to retrieve user payload")
	}
	if out, err1 := hookit.DebugExec(container.ID, "publish", payload, "info"); err1 != nil {
		if err2, ok := err1.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the (publish)publish hook")
		}
		return util.ErrorAppend(err1, "failed to run the (publish)publish hook")
	}

	display.StopTask()

	return nil
}
