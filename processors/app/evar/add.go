package evar

import (
	"fmt"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

func Add(envModel *models.Env, appModel *models.App, evars map[string]string) error {

	if err := app.Setup(envModel, appModel, appModel.Name); err != nil {
		return util.ErrorAppend(err, "failed to setup app")
	}

	// iterate through the evars and add them to the app
	for key, val := range evars {
		appModel.Evars[key] = val
	}

	// save the app
	if err := appModel.Save(); err != nil {
		return util.ErrorAppend(err, "failed to persist evars")
	}

	// iterate one more time for display
	fmt.Println()
	for key := range evars {
		fmt.Printf("%s %s added\n", display.TaskComplete, key)
	}
	fmt.Println()

	return nil
}
