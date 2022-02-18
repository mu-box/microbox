package code

import (
	"github.com/jcelliott/lumber"
	boxfile "github.com/mu-box/microbox-boxfile"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/locker"
)

// Sync syncronizes an app's code components with the boxfile config
func Sync(appModel *models.App, warehouseConfig WarehouseConfig) error {
	display.OpenContext("Syncing code components")
	defer display.CloseContext()

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	if err := purgeComponents(appModel); err != nil {
		return util.ErrorAppend(err, "failed to purge code components")
	}

	if err := provisionComponents(appModel, warehouseConfig); err != nil {
		return util.ErrorAppend(err, "failed to provision components")
	}

	return nil
}

// removes the code components from the app
func purgeComponents(appModel *models.App) error {
	display.OpenContext("Removing old")
	defer display.CloseContext()

	// get all the components
	componentModels, err := appModel.Components()
	if err != nil {
		lumber.Error("code:Clean:models.App{ID:%s}.Components(): %s", appModel.ID, err.Error())
		return err
	}

	// remove components that are of the code type
	for _, componentModel := range componentModels {

		// only destroy code type containers
		if componentModel.Type == "code" {

			// run a code destroy
			if err := Destroy(componentModel); err != nil {
				return util.ErrorAppend(err, "failed to destroy code component")
			}
		}
	}

	return nil
}

// provisions the code components for the app
func provisionComponents(appModel *models.App, warehouseConfig WarehouseConfig) error {
	display.OpenContext("Starting new")
	defer display.CloseContext()

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// iterate over the code nodes and build containers for each of them
	for _, componentModel := range codeComponentModels(appModel) {

		// run the code setup process with the new config
		err := Setup(appModel, componentModel, warehouseConfig)
		if err != nil {
			return util.ErrorAppend(err, "failed to setup code (%s): %s\n", componentModel.Name, err.Error())
		}

	}

	return nil
}

// setBoxfile ...
func codeComponentModels(appModel *models.App) []*models.Component {

	componentModels := []*models.Component{}

	// look in the boxfile for code nodes and generate a stub component
	box := boxfile.New([]byte(appModel.DeployedBoxfile))
	for _, componentName := range box.Nodes("code") {

		image := box.Node(componentName).StringValue("image")
		if image == "" {
			image = "mubox/code"
		}

		componentModel := &models.Component{
			Name:  componentName,
			Label: componentName,
			Image: image,
		}

		componentModels = append(componentModels, componentModel)
	}

	return componentModels
}
