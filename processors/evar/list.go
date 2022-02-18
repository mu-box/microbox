package evar

import (
	"fmt"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/odin"
)

func List(envModel *models.Env, appID string) error {
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

	// print the header
	fmt.Printf("\nEnvironment Variables\n")

	// iterate
	for _, evar := range evars {
		fmt.Printf("  %s = %s\n", evar.Key, evar.Value)
	}

	fmt.Println()

	return nil
}
