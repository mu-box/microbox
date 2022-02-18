package sim

import (
	// "fmt"
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/app"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

func init() {
	steps.Build("sim stop", stopCheck, stopFn)
}

// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	appModel, _ := models.FindAppBySlug(config.EnvID(), "sim")
	display.CommandErr(app.Stop(appModel))
}

func stopCheck() bool {
	// currently we always stop if we are asking weather to stop
	return false
}
