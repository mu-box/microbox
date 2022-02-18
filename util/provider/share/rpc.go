package share

import (
	"github.com/mu-box/microbox/commands/server"
)

type ShareRPC struct{}

type Response struct {
	Message string
	Success bool
}

func init() {
	server.Register(&ShareRPC{})
}
