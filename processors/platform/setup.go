package platform

import (
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

// Setup provisions platform components needed for an app setup
func Setup(appModel *models.App) error {
	display.OpenContext("Starting components")
	defer display.CloseContext()

	for _, component := range setupComponents {
		if err := provisionComponent(appModel, component); err != nil {
			return util.ErrorAppend(err, "failed to provision platform component")
		}
	}

	return nil
}
