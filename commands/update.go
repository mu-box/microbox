package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/display"
)

var (
	// UpdateCmd updates the images
	UpdateCmd = &cobra.Command{
		Use:   "update-images",
		Short: "Updates docker images.",
		// Short:  "Updates docker images and checks to see if the microbox binary needs an update.",
		Long:   ``,
		PreRun: steps.Run("start"),
		Run:    updateFn,
	}
)

// updateFn ...
func updateFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Update())
}
