package processors

import (
	"fmt"
	"os"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/odin"
)

// Process ...
func Login(username, password, endpoint string) error {

	// request Username/Password if missing
	if username == "" && os.Getenv("MICROBOX_USERNAME") != "" {
		username = os.Getenv("MICROBOX_USERNAME")
	}

	if username == "" {
		user, err := display.ReadUsername()
		if err != nil {
			return util.ErrorAppend(err, "unable to retrieve username")
		}
		username = user
	}

	if password == "" && os.Getenv("MICROBOX_PASSWORD") != "" {
		password = os.Getenv("MICROBOX_PASSWORD")
	}

	if password == "" {
		// ReadPassword prints Password: already
		pass, err := display.ReadPassword("Microbox")
		if err != nil {
			return util.ErrorAppend(err, "failed to read password")
		}
		password = pass
	}

	if endpoint == "" && os.Getenv("MICROBOX_ENDPOINT") != "" {
		endpoint = os.Getenv("MICROBOX_ENDPOINT")
	}

	if endpoint == "" {
		endpoint = "microbox"
	}

	// set the odin endpoint
	odin.SetEndpoint(endpoint)

	// verify that the user exists
	token, err := odin.Auth(username, password)
	if err != nil {
		fmt.Print(`! The username/password was incorrect, but we're continuing on.
  To reattempt authentication, run 'microbox login'.
`)
		return nil
	}

	// store the user token
	auth := models.Auth{
		Endpoint: endpoint,
		Key:      token,
	}
	if auth.Save() != nil {
		return util.Errorf("unable to save user authentication")
	}

	display.LoginComplete()

	return nil
}
