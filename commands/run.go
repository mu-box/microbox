package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/console"
	"github.com/mu-box/microbox/util/display"

	// imported because we need its steps added
	_ "github.com/mu-box/microbox/commands/dev"
)

// RunCmd ...
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Start your local development environment.",
	Long: `
Starts your local development environment and opens an
interactive console inside the environment.

You can also pass a command into 'run'. Microbox will
run the command without dropping you into a console
in your local environment.
	`,
	PreRun:  steps.Run("start", "build-runtime", "dev start", "dev deploy"),
	Run:     runFn,
	PostRun: steps.Run("dev stop"),
}

// runFn ...
func runFn(ccmd *cobra.Command, args []string) {

	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	consoleConfig := console.ConsoleConfig{}

	if len(args) > 0 {
		consoleConfig.Command = strings.Join(args, " ")
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(processors.Run(envModel, appModel, consoleConfig))
}

func init() {
	steps.Build("dev deploy", devDeployComplete, devDeploy)
}

// devDeploy ...
func devDeploy(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(envModel.ID, "dev")
	display.CommandErr(app.Deploy(envModel, appModel))
}

func devDeployComplete() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	env, _ := app.Env()
	return app.DeployedBoxfile != "" && env.BuiltBoxfile == app.DeployedBoxfile && buildComplete()
}
