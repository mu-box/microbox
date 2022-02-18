package processors

import (
	"fmt"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/microagent"
	"github.com/mu-box/microbox/util/odin"
)

func Tunnel(envModel *models.Env, tunnelConfig models.TunnelConfig) error {
	// fetch the remote
	remote, ok := envModel.Remotes[tunnelConfig.AppName]
	if ok {
		// set the odin endpoint
		odin.SetEndpoint(remote.Endpoint)
		// set the app id
		tunnelConfig.AppName = remote.Name
	}

	// set the app id to the directory name if it's default
	if tunnelConfig.AppName == "default" {
		tunnelConfig.AppName = config.AppName()
	}

	// set odins endpoint if the argument is passed
	if endpoint := registry.GetString("endpoint"); endpoint != "" {
		odin.SetEndpoint(endpoint)
	}

	// validate access to the app
	if err := helpers.ValidateOdinApp(tunnelConfig.AppName); err != nil {
		return util.ErrorAppend(err, "unable to validate app")
	}

	// initiate a tunnel session with odin
	tunInfo, err := odin.EstablishTunnel(tunnelConfig)
	if err != nil {
		return util.ErrorAppend(err, "failed to initiate a remote tunnel session")
	}

	// set a default port if the user didn't specify
	if tunnelConfig.ListenPort == 0 {
		tunnelConfig.ListenPort = tunInfo.Port
	}

	// connect up to the session
	if err := microagent.Tunnel(tunInfo.Token, tunInfo.URL, fmt.Sprint(tunnelConfig.ListenPort), tunnelConfig.Component); err != nil {
		return util.ErrorAppend(err, "failed to connect to remote tunnel session")
	}

	return nil
}
