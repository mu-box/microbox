package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"

	// added because we need its steps
	_ "github.com/mu-box/microbox/commands/sim"
)

var (

	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:   "deploy [dry-run|remote-alias]",
		Short: "Deploy your application to a live remote or a dry-run environment.",
		Long:  ``,
		PreRun: func(ccmd *cobra.Command, args []string) {
			registry.Set("skip-compile", deployCmdFlags.skipCompile)
			steps.Run("configure", "start", "build-runtime", "compile-app")(ccmd, args)
		},
		Run: deployFn,
	}

	// deployCmdFlags ...
	deployCmdFlags = struct {
		skipCompile bool
		message     string
		force       bool
	}{}
)

//
func init() {
	DeployCmd.Flags().BoolVarP(&deployCmdFlags.skipCompile, "skip-compile", "", false, "skip compiling the app")
	DeployCmd.Flags().BoolVarP(&deployCmdFlags.force, "force", "", false, "force the deploy even if you have used this build on a previous deploy")
	DeployCmd.Flags().StringVarP(&deployCmdFlags.message, "message", "m", "", "Allows you to append a message to the deploy. These messages appear in your app's deploy history in your dashboard.")
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(envModel, args, 1)

	switch location {
	case "local":
		switch name {
		case "dev":
			fmt.Println("deploying is not necessary in this context, 'microbox run' instead")
			return
		case "sim":
			steps.Run("sim start")(ccmd, args)
			appModel, _ := models.FindAppBySlug(envModel.ID, "sim")
			display.CommandErr(app.Deploy(envModel, appModel))
			steps.Run("sim stop")(ccmd, args)
		}
	case "production":
		steps.Run("login")(ccmd, args)
		deployConfig := processors.DeployConfig{
			App:     name,
			Message: deployCmdFlags.message,
			Force:   deployCmdFlags.force,
		}

		// set the meta arguments to be used in the processor and run the processor
		display.CommandErr(processors.Deploy(envModel, deployConfig))
	}
}
