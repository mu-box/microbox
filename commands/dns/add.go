package dns

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	app_dns "github.com/mu-box/microbox/processors/app/dns"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

// AddCmd ...
var AddCmd = &cobra.Command{
	Use:   "add [local|dry-run] <hostname>",
	Short: "Adds dns entries",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: addFn,
}

// addFn ...
func addFn(ccmd *cobra.Command, args []string) {

	// parse the dnss excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 2)

	if len(args) != 1 {
		fmt.Println("i need a dns")
	}

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		app.Generate(env, name)
		display.CommandErr(app_dns.Add(env, app, args[0]))
	case "production":
		fmt.Printf(`
--------------------------------------------------------
Production dns aliasing is not yet implemented.
--------------------------------------------------------

`)
	}
}
