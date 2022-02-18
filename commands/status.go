package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/display"
)

var (

	// StatusCmd ...
	StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Display the status of your Microbox VM & apps.",
		Long:  ``,
		Run:   statusFn,
	}
)

func statusFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Status())
}
