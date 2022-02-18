package processors

import (
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/code"
	"github.com/mu-box/microbox/processors/env"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/locker"
)

// Build sets up the environment and runs a code build
func Build(envModel *models.Env) error {
	// by acquiring a local lock we are only allowing
	// one build to happen at a time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// init docker client and env mounts
	if err := env.Setup(envModel); err != nil {
		return util.ErrorAppend(err, "failed to prepare environment")
	}

	// print a warning if this is the first build
	if envModel.BuiltBoxfile == "" {
		display.FirstBuild()
	}

	// build code
	if err := code.Build(envModel); err != nil {
		return util.ErrorAppend(err, "failed to build the code")
	}

	return nil
}
