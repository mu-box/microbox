package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/processors/provider/bridge"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/provider"
	"github.com/mu-box/microbox/util/service"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the Microbox virtual machine.",
		Long:  ``,
		Run:   startFn,
	}
)

func init() {
	steps.Build("start", startCheck, startFn)
}

// startFn ...
func startFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Start())
}

func startCheck() bool {
	bridgeReady := true
	if provider.BridgeRequired() {
		bridgeReady = bridge.Connected()
	}
	return provider.IsReady() && service.Running("microbox-server") && bridgeReady
}
