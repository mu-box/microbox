package server

import (
	"time"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/service"
)

func Start() error {
	// run as admin
	// the reExecPrivilegeStart function is defined in the setup
	// since the service create is idempotent it is fine to only have one
	// start command for the server
	if !util.IsPrivileged() {
		return reExecPrivilegeStart()
	}

	fn := func() error {
		return service.Start("microbox-server")
	}

	return util.Retry(fn, 3, time.Second)
}
