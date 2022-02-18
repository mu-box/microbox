package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/display"
)

var (

	// StopCmd ...
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop the Microbox virtual machine.",
		Long: `
Stops the Microbox virtual machine as well as
any running local or dry-run environments.
		`,
		Run: stopFn,
	}
)

// stopFn ...
func stopFn(ccmd *cobra.Command, args []string) {
	registry.Set("keep-share", true)
	display.CommandErr(processors.Stop())
}
