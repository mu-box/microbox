package processors

import (
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/env"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/locker"
)

//
func Clean(envModels []*models.Env) error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	display.OpenContext("Cleaning stale environments")
	defer display.CloseContext()

	// if any of the apps are stale, we'll mark this to true
	stale := false

	for _, envModel := range envModels {
		// check to see if the app folder still exists
		if !util.FolderExists(envModel.Directory) {

			if err := env.Destroy(envModel); err != nil {
				return util.ErrorAppend(err, "unable to destroy environment(%s)", envModel.Name)
			}
		}
	}

	if !stale {
		display.StartTask("Skipping (none detected)")
		display.StopTask()
	}

	return nil
}
