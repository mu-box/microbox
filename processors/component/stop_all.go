package component

import (
	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

// StopAll stops all app components
func StopAll(appModel *models.App) error {

	// get all the components that belong to this app
	componentModels, err := appModel.Components()
	if err != nil {
		lumber.Error("component:StopAll:models.App{ID:%s}.Components() %s", appModel.ID, err.Error())
		return util.ErrorAppend(err, "unable to retrieve components")
	}

	if len(componentModels) == 0 {
		return nil
	}

	display.OpenContext("Stopping components")
	defer display.CloseContext()

	// stop each component
	for _, componentModel := range componentModels {
		if err := Stop(componentModel); err != nil {
			return util.ErrorAppend(err, "unable to stop component(%s)", componentModel.Name)
		}
	}

	return nil
}
