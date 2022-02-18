package app

import (
	"fmt"
	// "net"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/component"
	process_provider "github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/util"

	// "github.com/mu-box/microbox/util/dhcp"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/locker"
)

// Stop will stop all services associated with an app
func Stop(appModel *models.App) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is already down
	// TODO: also check if any containers are running
	if appModel.Status != "up" {
		return nil
	}

	// load the env for the display context
	envModel, err := appModel.Env()
	if err != nil {
		lumber.Error("app:Stop:models.App.Env()")
		return util.ErrorAppend(err, "failed to load app env")
	}

	display.OpenContext("%s (%s)", envModel.Name, appModel.DisplayName())
	defer display.CloseContext()

	// initialize docker for the provider
	if err := process_provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to initialize docker environment")
	}

	// stop all app components
	if err := component.StopAll(appModel); err != nil {
		return util.ErrorAppend(err, "failed to stop all app components")
	}

	display.StartTask("Pausing App")
	display.StopTask()

	// stop any dev containers
	stopDevContainer(appModel)

	// set the status to down
	appModel.Status = "down"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Stop:models.App.Save()")
		return util.ErrorAppend(err, "failed to persist app status")
	}

	return nil
}

func stopDevContainer(appModel *models.App) error {
	// grab the container info
	container, err := docker.GetContainer(fmt.Sprintf("microbox_%s", appModel.ID))
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}

	// remove the container
	if err := docker.ContainerRemove(container.ID); err != nil {
		lumber.Error("dev:console:teardown:docker.ContainerRemove(%s): %s", container.ID, err)
		// we cannot trust that the container is there even though we checked for it
		// return util.ErrorAppend(err, "failed to remove dev container")
	}

	// // extract the container IP
	// ip := docker.GetIP(container)

	// // return the container IP back to the IP pool
	// if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
	// 	lumber.Error("dev:console:teardown:dhcp.ReturnIP(%s): %s", ip, err)

	// 	lumber.Error("An error occurred during dev console teadown:%s", err.Error())
	// 	return util.ErrorAppend(err, "failed to return unused IP back to pool")
	// }
	return nil
}
