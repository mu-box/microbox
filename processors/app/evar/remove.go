package evar

import (
	"fmt"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

func Remove(appModel *models.App, keys []string) error {

	// delete the evars
	for _, key := range keys {
		delete(appModel.Evars, key)
	}

	// persist the app model
	if err := appModel.Save(); err != nil {
		return util.ErrorAppend(err, "failed to delete evars")
	}

	// print the deleted keys
	fmt.Println()
	for _, key := range keys {
		fmt.Printf("%s %s removed\n", display.TaskComplete, key)
	}
	fmt.Println()

	return nil
}
