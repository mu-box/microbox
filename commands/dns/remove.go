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

// RemoveCmd ...
var RemoveCmd = &cobra.Command{
	Use:   "rm [local|dry-run] <hostname>",
	Short: "Remove dns entries",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: removeFn,
}

// removeFn ...
func removeFn(ccmd *cobra.Command, args []string) {
	// parse the dnss excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)

	if len(args) != 1 {
		fmt.Println("i need a dns")
	}

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_dns.Remove(app, args[0]))
	case "production":
		fmt.Printf(`
--------------------------------------------------------
Production dns aliasing is not yet implemented.
--------------------------------------------------------

`)
	}
}
