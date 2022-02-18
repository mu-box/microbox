package evar

import (
	"fmt"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/odin"
)

func Add(envModel *models.Env, appID string, evars map[string]string) error {

	// fetch the remote
	remote, ok := envModel.Remotes[appID]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		appID = remote.ID
	}

	// set odins endpoint if the argument is passed
	if endpoint := registry.GetString("endpoint"); endpoint != "" {
		odin.SetEndpoint(endpoint)
	}

	// iterate through the evars and add them to the app
	for key, val := range evars {
		err := odin.AddEvar(appID, key, val)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s added\n", display.TaskComplete, key)
	}

	return nil
}
