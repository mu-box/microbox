package server

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/service"
)

func Setup() error {
	if service.Running("microbox-server") {
		return nil
	}

	// run as admin
	if !util.IsPrivileged() {
		return reExecPrivilegeStart()
	}

	// TEMP: we need to remove the old microbox-vpn just incase it is left over
	// we will not catch errors here because if it doesnt exist or it breaks it
	// should not stop us from creating the new microbox-server
	service.Stop("microbox-vpn")
	service.Remove("microbox-vpn")

	// create the service this call is idempotent so we shouldnt need to check
	if err := service.Create("microbox-server", []string{config.MicroboxPath(), "server"}); err != nil {
		return err
	}

	// start the service
	return Start()

}

// reExecPrivilegeStart re-execs the current process with a privileged user
func reExecPrivilegeStart() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to start the server")

	cmd := fmt.Sprintf("\"%s\" env server start", config.MicroboxPath())
	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("server:reExecPrivilegeStart:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
