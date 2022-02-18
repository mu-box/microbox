package bridge

import (
	"github.com/mu-box/microbox/util/provider/bridge"
)

// ask the server to start the bridge
func Start() error {
	return bridge.Start(ConfigFile())
}
