package provider

import (
	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/processors/provider/bridge"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/locker"
	"github.com/mu-box/microbox/util/provider"
)

// Stop stops the provider (stops the VM)
func Stop() error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Stopping Microbox")
	defer display.CloseContext()

	// stop the vpn
	if err := bridge.Stop(); err != nil {
		// do nothing about the error since provider stop happens
		// then we shut down the server (killing the bridge)
		// return util.ErrorAppend(err, "failed to stop vpn")
	}

	// stop the provider (VM)
	if err := provider.Stop(); err != nil {
		lumber.Error("provider:Stop:provider.Stop(): %s", err.Error())
		return util.ErrorAppend(err, "failed to stop the provider")
	}

	return nil
}
