package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/display"
)

var (

	// CleanCmd ...
	CleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean out any apps that no longer exist.",
		Long: `
Clean out any apps whose working directory no longer exists. This
will remove all associated app information from your Microbox database.
`,
		PreRun: steps.Run("start"),
		Run:    cleanFn,
	}
)

// cleanFn ...
func cleanFn(ccmd *cobra.Command, args []string) {
	// get the environments
	envs, err := models.AllEnvs()
	if err != nil {
		fmt.Printf("TODO: write message for command clean: %s\n", err.Error())
		return
	}

	display.CommandErr(processors.Clean(envs))
}
