package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/processors/env"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

var (

	// DestroyCmd ...
	DestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current project and remove it from Microbox.",
		Long: `
Destroys the current project and removes it from Microbox – destroying
the filesystem mount, associated dns aliases, and local app data.
		`,
		PreRun: steps.Run("start"),
		Run:    destroyFunc,
	}
)

// destroyFunc ...
func destroyFunc(ccmd *cobra.Command, args []string) {
	envModel, err := models.FindEnvByID(config.EnvID())
	if err != nil {
		fmt.Println("This project doesn't exist on microbox.")
		return
	}

	if len(args) == 0 {
		display.CommandErr(env.Destroy(envModel))
		return
	}

	_, _, name := helpers.Endpoint(envModel, args, 2)
	appModel, err := models.FindAppBySlug(envModel.ID, name)
	if err != nil {
		fmt.Println("Could not find the application")
	}

	display.CommandErr(app.Destroy(appModel))

}
