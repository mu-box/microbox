package app

import (
	"fmt"
	"net"
	"strings"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app/dns"
	"github.com/mu-box/microbox/processors/component"
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/dhcp"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/locker"
)

// Destroy removes the app from the provider and the database
func Destroy(appModel *models.App) error {
	// init docker client
	if err := provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if this app isn't created
	if appModel.IsNew() {
		return nil
	}

	// load the env for the display context
	envModel, err := appModel.Env()
	if err != nil {
		lumber.Error("app:Start:models.App.Env()")
		return util.ErrorAppend(err, "failed to load app env")
	}

	if err := dns.RemoveAll(appModel); err != nil {
		return util.ErrorAppend(err, "failed to remove dns aliases")
	}

	display.OpenContext("%s (%s)", envModel.Name, appModel.DisplayName())
	defer display.CloseContext()

	// remove the dev container if there is one
	docker.ContainerRemove(fmt.Sprintf("microbox_%s", appModel.ID))

	// destroy the associated components
	if err := destroyComponents(appModel); err != nil {
		return util.ErrorAppend(err, "failed to destroy components")
	}

	// release IPs
	if err := releaseIPs(appModel); err != nil {
		return util.ErrorAppend(err, "failed to release IPs")
	}

	// destroy the app model
	if err := appModel.Delete(); err != nil {
		lumber.Error("app:Destroy:models.App{ID:%s}.Destroy(): %s", appModel.ID, err.Error())
		return util.ErrorAppend(err, "failed to delete app model")
	}

	cleanImages()
	return nil
}

// destroyComponents destroys all the components of this app
func destroyComponents(appModel *models.App) error {
	display.OpenContext("Removing components")
	defer display.CloseContext()

	componentModels, err := appModel.Components()
	if err != nil {
		lumber.Error("app:destroyComponents:models.App{ID:%s}.Components() %s", appModel.ID, err.Error())
		return util.ErrorAppend(err, "unable to retrieve components")
	}

	if len(componentModels) == 0 {
		display.StartTask("Skipping (no components)")
		display.StopTask()
		return nil
	}

	for _, componentModel := range componentModels {
		if err := component.Destroy(appModel, componentModel); err != nil {
			return util.ErrorAppend(err, "failed to destroy app component")
		}
	}

	return nil
}

// releaseIPs releases the app-level ip addresses
func releaseIPs(appModel *models.App) error {
	display.StartTask("Releasing IPs")
	defer display.StopTask()

	// release all of the local IPs
	for _, ip := range appModel.LocalIPs {
		// release the IP
		if err := dhcp.ReturnIP(net.ParseIP(ip)); err != nil {
			display.ErrorTask()
			lumber.Error("app:Destroy:releaseIPs:dhcp.ReturnIP(%s): %s", ip, err.Error())
			return util.ErrorAppend(err, "failed to release IP")
		}
	}

	return nil
}

func cleanImages() {
	images, err := docker.ImageList()
	if err != nil {
		return
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			// if there is a tag that is not our build and it is one of ours.. try removing it (without force)
			if tag != "" && !strings.HasPrefix(tag, "mubox/build") && strings.HasPrefix(tag, "mubox/") {
				tag = strings.Replace(tag, ":latest", "", 1)
				docker.ImageRemove(tag, false)
			}
		}
	}
}
