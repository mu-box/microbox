package processors

import (
	"github.com/mu-box/microbox/processors/provider"
	"github.com/mu-box/microbox/processors/server"
)

// Start starts the provider (VM)
func Start() error {
	// start the microbox server
	if err := server.Setup(); err != nil {
		return err
	}

	// run a provider setup
	return provider.Setup()
}
