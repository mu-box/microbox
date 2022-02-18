package evar

import (
	"fmt"

	"github.com/mu-box/microbox/models"
)

func List(appModel *models.App) error {

	// print the header
	fmt.Printf("\nEnvironment Variables\n")

	// iterate
	for key, val := range appModel.Evars {
		fmt.Printf("  %s = %s\n", key, val)
	}

	fmt.Println()

	return nil
}
