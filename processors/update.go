// This update processor updates the images on microbox local

package processors

import (
	"strings"
	"time"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	process_provider "github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

func Update() error {

	// init docker client
	if err := process_provider.Init(); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// update all the microbox images
	if err := pullImages(); err != nil {
		return util.ErrorAppend(err, "failed to pull images")
	}

	return nil
}

func pullImages() error {
	display.OpenContext("Updating Images")
	defer display.CloseContext()

	images, err := docker.ImageList()
	if err != nil {
		return err
	}

	for _, image := range images {
		if image.Slug == "" {
			continue
		}
		if !strings.Contains(image.Slug, "mubox/") {
			continue
		}
		display.StartTask("Pulling %s image", image.Slug)

		// generate a docker percent display
		dockerPercent := &display.DockerPercentDisplay{
			Output: display.NewStreamer("info"),
		}

		// pull the build image
		imagePullFunc := func() error {
			_, err := docker.ImagePull(image.Slug, dockerPercent)
			return err
		}

		if err := util.Retry(imagePullFunc, 5, time.Second); err != nil {
			lumber.Error("code:pullBuildImage:docker.ImagePull(%s, nil): %s", image.Slug, err.Error())
			display.ErrorTask()
			return util.ErrorAppend(err, "failed to pull docker image (%s)", image.Slug)
		}

		display.StopTask()
	}

	return nil
}
