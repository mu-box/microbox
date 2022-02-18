package evar

import (
	"fmt"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/odin"
)

func Remove(envModel *models.Env, appID string, keys []string) error {

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

	evars, err := odin.ListEvars(appID)
	if err != nil {
		return err
	}

	// delete the evars
	for _, key := range keys {
		removed := false
		for _, evar := range evars {
			if evar.Key == key {
				if err := odin.RemoveEvar(appID, evar.ID); err != nil {
					return err
				}
				removed = true
				fmt.Printf("%s %s removed\n", display.TaskComplete, key)
			}
		}
		if !removed {
			fmt.Printf("%s %s not found\n", display.TaskPause, key)
		}
	}
	fmt.Println()

	return nil
}
