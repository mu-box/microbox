package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// InfoCmd ...
	InfoCmd = &cobra.Command{
		Use:    "info",
		Short:  "Displays information about the running sim app and its components.",
		Long:   ``,
		Run:    infoFn,
	}
)

// infoFn will run the DNS processor for adding DNS entires to the "hosts" file
func infoFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	
	display.CommandErr(sim.Info(env, app))
}
