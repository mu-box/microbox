package dns

import (
	"fmt"

	"github.com/spf13/cobra"

	// "github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	app_dns "github.com/mu-box/microbox/processors/app/dns"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

// ListCmd ...
var ListCmd = &cobra.Command{
	Use:   "ls [local|dry-run]",
	Short: "list dns entries",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: listFn,
}

// listFn ...
func listFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	_, location, name := helpers.Endpoint(env, args, 0)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_dns.List(app))
	case "production":
		fmt.Printf(`
--------------------------------------------------------
Production dns aliasing is not yet implemented.
--------------------------------------------------------

`)
	}
}
