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

var AppSetup func(envModel *models.Env, appModel *models.App, name string) error

// Add adds a dns entry to the local hosts file
func Add(envModel *models.Env, appModel *models.App, name string) error {

	if err := AppSetup(envModel, appModel, appModel.Name); err != nil {
		return util.ErrorAppend(err, "failed to setup app")
	}

	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := appModel.LocalIPs["env"]

	// generate the dns entry
	entry := dns.Entry(envIP, name, appModel.ID)

	// short-circuit if this entry already exists
	if dns.Exists(entry) {
		return nil
	}

	// make sure the server is running since it will do the dns addition
	if err := server.Setup(); err != nil {
		return util.ErrorAppend(err, "failed to setup server")
	}

	// issue a warning about `.dev` being unusable with Chrome
	if name[len(name)-4:] == ".dev" {
		tld := appModel.DisplayName()
		if tld == "dry-run" {
			tld = "test"
		}

		return util.Errorf("Google has been locking down use of the .dev TLD, and your app may not be accessible with this domain.\nTry using %s.%s instead.", name[0:len(name)-4], tld)
	}

	// add the entry
	if err := dns.Add(entry); err != nil {
		lumber.Error("dns:Add:dns.Add(%s): %s", entry, err.Error())
		return util.ErrorAppend(err, "unable to add dns entry")
	}

	display.Info("\n%s %s added\n", display.TaskComplete, name)

	return nil
}
