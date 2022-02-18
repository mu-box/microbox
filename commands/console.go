package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/processors/env"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/console"
	"github.com/mu-box/microbox/util/display"
)

var (

	// ConsoleCmd ...
	ConsoleCmd = &cobra.Command{
		Use:   "console [<local | dry-run | {remote-alias}>] <component.id>",
		Short: "Open an interactive console inside a component.",
		Long:  ``,
		Run:   consoleFn,
	}
	user string
)

func init() {
	ConsoleCmd.Flags().StringVarP(&user, "user", "u", "", "user you would like to console in as")
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {
	if user != "" {
		registry.Set("console_user", user)
	}
	envModel, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(envModel, args, 2)

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the component you wish to console into:

ex: microbox console dry-run web.site

`, len(args))
		return
	}

	if name == "dev" && isCode(strings.Join(args, " ")) {
		display.ConsoleLocalCode()
		display.CommandErr(util.Err{
			Message: "Console to local code node not valid",
			Code:    "USER",
			Stack:   []string{"failed to console"},
			Suggest: "It appears you are trying to console to a local code node. Please use `microbox run` instead.",
		})
		return
	}

	switch location {
	case "local":
		appModel, _ := models.FindAppBySlug(config.EnvID(), name)
		if appModel.Status != "up" {
			fmt.Println("unable to continue until the app is up")
			return
		}

		componentModel, _ := models.FindComponentBySlug(config.EnvID()+"_"+name, args[0])
		// todo: determine ways this errors and handle it, use util.Err for suggestions
		// componentModel, err := models.FindComponentBySlug(config.EnvID()+"_"+name, args[0])
		// if err != nil {
		// 	display.CommandErr(err)
		// 	return
		// }

		display.CommandErr(env.Console(componentModel, console.ConsoleConfig{}))

	case "production":

		consoleConfig := processors.ConsoleConfig{
			App:  name,
			Host: args[0],
		}

		// set the meta arguments to be used in the processor and run the processor
		display.CommandErr(processors.Console(envModel, consoleConfig))

	}
}

func isCode(args string) bool {
	if strings.Contains(args, "web.") || strings.Contains(args, "worker.") {
		return true
	}
	return false
}
