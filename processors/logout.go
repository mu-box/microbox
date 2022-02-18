package processors

import (
	"fmt"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

// Process ...
func Logout(endpoint string) error {

	if endpoint == "" {
		endpoint = "microbox"
	}

	// lookup the auth by the endpoint
	auth, _ := models.LoadAuthByEndpoint(endpoint)

	// short-circuit if the auth is already deleted
	if auth.IsNew() {
		fmt.Printf("%s Already logged out\n", display.TaskComplete)
		return nil
	}

	// remove token from database
	if err := auth.Delete(); err != nil {
		return util.ErrorAppend(err, "failed to delete user authentication")
	}

	fmt.Printf("%s You've logged out\n", display.TaskComplete)

	return nil
}
