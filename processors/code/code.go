// Package code ...
package code

import (
	"time"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"
	boxfile "github.com/mu-box/microbox-boxfile"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

// these constants represent different potential names a service can have
const (
	BUILD = "build"
)

// these constants represent different potential states an app can end up in
const (
	ACTIVE = "active"
)

func pullBuildImage() (string, error) {
	// extract the build image from the boxfile
	buildImage := buildImage()

	if docker.ImageExists(buildImage) {
		return buildImage, nil
	}

	display.StartTask("Pulling %s image", buildImage)
	defer display.StopTask()

	// generate a docker percent display
	dockerPercent := &display.DockerPercentDisplay{
		Output: display.NewStreamer("info"),
		// Prefix: buildImage,
	}

	// pull the build image
	imagePull := func() error {
		_, err := docker.ImagePull(buildImage, dockerPercent)
		return err
	}
	if err := util.Retry(imagePull, 5, time.Second); err != nil {
		lumber.Error("code:pullBuildImage:docker.ImagePull(%s, nil): %s", buildImage, err.Error())
		display.ErrorTask()
		return "", util.ErrorAppend(err, "failed to pull docker image (%s)", buildImage)
	}

	return buildImage, nil
}

// BuildImage fetches the build image from the boxfile
func buildImage() string {
	// first let's see if the user has a custom build image they want to use
	box := boxfile.NewFromPath(config.Boxfile())
	image := box.Node("run.config").StringValue("image")

	// then let's set the default if the user hasn't specified
	if image == "" {
		image = "mubox/build"
	}

	return image
}
