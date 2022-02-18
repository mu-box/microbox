package processors

import (
	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/code"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/odin"
)

//
func Deploy(envModel *models.Env, deployConfig DeployConfig) error {

	appID := deployConfig.App

	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}

	// set the app id to the directory name if it's default
	if appID == "default" {
		appID = config.AppName()
	}

	// validate access to the app
	if err := helpers.ValidateOdinApp(appID); err != nil {
		return util.ErrorAppend(err, "unable to validate app")
	}

	warehouseConfig, err := getWarehouseConfig(envModel, appID)
	if err != nil {
		return util.ErrorAppend(err, "unable to generate warehouse config")
	}

	// print the first deploy message if this is the first deploy for the app
	if warehouseConfig.PreviousBuild == "" {
		display.FirstDeploy()
	}

	// publish to remote warehouse
	if err := code.Publish(envModel, warehouseConfig); err != nil {
		return util.ErrorAppend(err, "failed to publish build to app's warehouse")
	}

	// tell odin what happened
	if err := odin.Deploy(appID, warehouseConfig.BuildID, envModel.BuiltBoxfile, deployConfig.Message); err != nil {
		lumber.Error("deploy:odin.Deploy(%s,%s,%s,%s): %s", appID, warehouseConfig.BuildID, envModel.BuiltBoxfile, deployConfig.Message, err.Error())
		return util.ErrorAppend(err, "failed to deploy code to app")
	}

	envModel.DeployedID = envModel.BuiltID
	if err := envModel.Save(); err != nil {
		lumber.Error("deploy:models:Env:Save(): %s", err.Error())
		return util.ErrorAppend(err, "failed to save build ID")
	}

	display.DeployComplete()

	return nil
}

// setWarehouseToken ...
func getWarehouseConfig(envModel *models.Env, appID string) (warehouseConfig code.WarehouseConfig, err error) {

	token, url, err := odin.GetWarehouse(appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetWarehouse(%s): %s", appID, err.Error())
		err = util.ErrorAppend(err, "failed to fetch warehouse information from microbox")
		return
	}

	// get the previous build if there was one
	prevBuild, err := odin.GetPreviousBuild(appID)
	if err != nil {
		lumber.Error("deploy:setWarehouseToken:GetPreviousBuild(%s): %s", appID, err.Error())
		err = util.ErrorAppend(err, "failed to query previous deploys from microbox")
		return
	}

	warehouseConfig.BuildID = envModel.BuiltID
	warehouseConfig.WarehouseURL = url
	warehouseConfig.WarehouseToken = token
	warehouseConfig.PreviousBuild = prevBuild

	return
}
