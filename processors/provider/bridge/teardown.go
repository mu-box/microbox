package bridge

import (
	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	container_generator "github.com/mu-box/microbox/generators/containers"
	"github.com/mu-box/microbox/util"
)

func Teardown() error {
	if !Connected() {
		return nil
	}

	// remove bridge client
	if err := Stop(); err != nil {
		return err
	}

	// remove component
	if err := removeComponent(); err != nil {
		return err
	}

	return nil
}

func removeComponent() error {
	// grab the container info
	container, err := docker.GetContainer(container_generator.BridgeName())
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		// if we cant talk to docker its ok too
		return nil
	}

	// remove the container
	if err := docker.ContainerRemove(container.ID); err != nil {
		lumber.Error("provider:bridge:teardown:docker.ContainerRemove(%s): %s", container.ID, err)
		return util.ErrorAppend(err, "failed to remove bridge container")
	}

	return nil
}
