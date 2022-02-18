package dns

import (
	// "fmt"

	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/server"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/dns"
)

// RemoveAll removes all dns entries for an app
func RemoveAll(a *models.App) error {

	// shortcut if we dont have any entries for this app
	if len(dns.List(a.ID)) == 0 {
		return nil
	}

	// make sure the server is running since it will do the dns work
	if err := server.Setup(); err != nil {
		return util.ErrorAppend(err, "failed to setup server")
	}

	if err := dns.Remove(a.ID); err != nil {
		lumber.Error("dns:RemoveAll:dns.Remove(%s): %s", a.ID, err.Error())
		return util.ErrorAppend(err, "failed to remove all dns entries")
	}

	display.Info("\n%s removed all\n", display.TaskComplete)
	return nil
}
