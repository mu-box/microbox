package server

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/service"
)

func Teardown() error {
	// run as admin
	if !util.IsPrivileged() {
		return reExecPrivilegeRemove()
	}

	// make sure its stopped first
	if err := Stop(); err != nil {
		return err
	}

	return service.Remove("microbox-server")
}

// reExecPrivilegeRemove re-execs the current process with a privileged user
func reExecPrivilegeRemove() error {
	display.PauseTask()
	defer display.ResumeTask()

	display.PrintRequiresPrivilege("to remove the server")

	cmd := fmt.Sprintf("\"%s\" env server teardown", config.MicroboxPath())

	if err := util.PrivilegeExec(cmd); err != nil {
		lumber.Error("server:reExecPrivilegeRemove:util.PrivilegeExec(%s): %s", cmd, err)
		return err
	}

	return nil
}
