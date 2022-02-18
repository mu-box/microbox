package server

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/service"
)

func Stop() error {
	// run as admin
	if !util.IsPrivileged() {
		return reExecPrivilegeStop()
	}

	return service.Stop("microbox-server")
}

// reExecPrivilegeStop re-execs the current process with a privileged user
func reExecPrivilegeStop() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to stop the server")

	cmd := fmt.Sprintf("\"%s\" env server stop", config.MicroboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("server:reExecPrivilegeStop:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
