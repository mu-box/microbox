package commands

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

var (

	// CompileCmd ...
	CompileCmd = &cobra.Command{
		Use:   "compile-app",
		Short: "Compile your application.",
		Long: `
Compiles your application source code into a deployable package.
		`,
		PreRun:  steps.Run("start", "build-runtime"),
		Run:     compileFn,
		Aliases: []string{"compile"},
	}
)

func init() {
	steps.Build("compile-app", compileComplete, compileFn)
}

// compileFn ...
func compileFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	display.CommandErr(processors.Compile(env))
}

func compileComplete() bool {
	env, _ := models.FindEnvByID(config.EnvID())
	// if the last compile has been set and they want to skip
	return !env.LastCompile.Equal(time.Time{}) && registry.GetBool("skip-compile")
}
