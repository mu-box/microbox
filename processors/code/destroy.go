package code

import (
	"net"

	"github.com/jcelliott/lumber"

	docker "github.com/mu-box/golang-docker-client"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/dhcp"
	"github.com/mu-box/microbox/util/display"
)

// Destroy destroys a code component from the app
func Destroy(componentModel *models.Component) error {
	display.OpenContext(componentModel.Label)
	defer display.CloseContext()

	// remove the docker container
	if err := destroyContainer(componentModel.ID); err != nil {
		return err
	}

	// detach from the host network
	if err := detachNetwork(componentModel); err != nil {
		return util.ErrorAppend(err, "failed to detach container from the host network")
	}

	// remove the componentModel from the database
	if err := componentModel.Delete(); err != nil {
		lumber.Error("code:Destroy:Component.Delete()")
		display.ErrorTask()
		return util.ErrorAppend(err, "unable to delete database model")
	}

	return nil
}

// destroys a docker container associated with this app
func destroyContainer(id string) error {
	display.StartTask("Destroying docker container")
	defer display.StopTask()

	if id == "" {
		return nil
	}

	if err := docker.ContainerRemove(id); err != nil {
		lumber.Error("component:Destroy:docker.ContainerRemove(%s): %s", id, err.Error())
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to remove docker container")
	}

	return nil
}

// detachNetwork detaches the network from the host
func detachNetwork(componentModel *models.Component) error {
	display.StartTask("Releasing IPs")
	defer display.StopTask()

	//
	if err := dhcp.ReturnIP(net.ParseIP(componentModel.IPAddr())); err != nil {
		lumber.Error("code:Destroy:dhcp.ReturnIP(%s): %s", componentModel.IPAddr(), err.Error())
		display.ErrorTask()
		return util.ErrorAppend(err, "")
	}

	return nil
}
