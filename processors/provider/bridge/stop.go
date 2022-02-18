package bridge

import (
	"github.com/mu-box/microbox/util/provider/bridge"
)

// ask the server to stop the bridge
func Stop() error {
	return bridge.Stop()
}
