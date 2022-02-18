package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/display"
)

var (

	// ImplodeCmd ...
	ImplodeCmd = &cobra.Command{
		Use:   "implode",
		Short: "Remove all Microbox-created containers, files, & data.",
		Long: `
Removes the Microbox container, all projects, filesystem mounts,
& local data. All that will remain is microbox binaries.
		`,
		Run: implodeFn,
	}
)

// implodeFn ...
func implodeFn(ccmd *cobra.Command, args []string) {
	registry.Set("full-implode", true)
	display.CommandErr(processors.Implode())
}
