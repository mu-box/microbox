package code

import (
	"strings"
	"time"

	"github.com/jcelliott/lumber"
	docker "github.com/mu-box/golang-docker-client"

	container_generator "github.com/mu-box/microbox/generators/containers"
	hook_generator "github.com/mu-box/microbox/generators/hooks/build"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/hookit"
)

// Compile builds the codebase that can then be deployed
func Compile(envModel *models.Env) error {
	display.OpenContext("Compiling application")
	defer display.CloseContext()

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return util.ErrorAppend(err, "failed to pull the compile image")
	}

	// if a compile container was leftover from a previous compile, let's remove it
	docker.ContainerRemove(container_generator.CompileName())

	display.StartTask("Starting docker container")

	// start the container
	config := container_generator.CompileConfig(buildImage)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Compile:docker.CreateContainer(%+v): %s", config, err.Error())
		return util.ErrorAppend(err, "failed to start docker container")
	}

	display.StopTask()

	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container_generator.CompileName())

	if err := prepareCompileEnvironment(container.ID); err != nil {
		return err
	}

	if err := compileCode(container.ID); err != nil {
		return err
	}

	// update the compiled flag
	envModel.LastCompile = time.Now()
	envModel.BuiltID = util.RandomString(30)

	return envModel.Save()
}

// prepareCompileEnvironment runs hooks to prepare the build environment
func prepareCompileEnvironment(containerID string) error {
	display.StartTask("Preparing environment for compile")
	defer display.StopTask()

	// run the user hook
	if out, err := hookit.DebugExec(containerID, "user", hook_generator.UserPayload(), "info"); err != nil {
		// handle 'exec failed: argument list too long' error
		if strings.Contains(out, "argument list too long") {
			if err2, ok := err.(util.Err); ok {
				err2.Suggest = "You may have too many ssh keys, please specify the one you need with `microbox config set ssh-key ~/.ssh/id_rsa`"
				err2.Output = out
				err2.Code = "1001"
				return util.ErrorAppend(err2, "failed to run the (compile)user hook")
			}
		}
		return util.ErrorAppend(err, "failed to run the (compile)user hook")
	}

	// run the configure hook
	if out, err := hookit.DebugExec(containerID, "configure", hook_generator.ConfigurePayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the (compile)configure hook")
		}
		return util.ErrorAppend(err, "failed to run the (compile)configure hook")
	}

	// run the boxfile hook
	if out, err := hookit.DebugExec(containerID, "boxfile", hook_generator.BoxfilePayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the (compile)boxfile hook")
		}
		return util.ErrorAppend(err, "failed to run the (compile)boxfile hook")
	}

	// run the mount hook
	if out, err := hookit.DebugExec(containerID, "mount", hook_generator.MountPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the (compile)mount hook")
		}
		return util.ErrorAppend(err, "failed to run the (compile)mount hook")
	}

	return nil
}

// compileCode runs the hooks to compile the codebase
func compileCode(containerID string) error {

	display.StartTask("Compiling code")
	defer display.StopTask()

	// run the compile hook
	if out, err := hookit.DebugExec(containerID, "compile", hook_generator.CompilePayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the (compile)compile hook")
		}
		return util.ErrorAppend(err, "failed to run the (compile)compile hook")
	}

	// run the pack-app hook
	if out, err := hookit.DebugExec(containerID, "pack-app", hook_generator.PackAppPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the (compile)pack-app hook")
		}
		return util.ErrorAppend(err, "failed to run the (compile)pack-app hook")
	}

	return nil
}
