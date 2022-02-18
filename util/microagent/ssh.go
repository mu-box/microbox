package microagent

import (
	"fmt"
	"strconv"
	"strings"

	ssh "github.com/mu-box/golang-ssh"

	"github.com/mu-box/microbox/util/display"
)

func SSH(key, location string) error {

	// create the ssh client
	nanPass := ssh.Auth{Passwords: []string{key}}
	locationParts := strings.Split(location, ":")
	if len(locationParts) != 2 {
		return fmt.Errorf("location is not formatted properly (%s)", location)
	}

	// parse port
	port, err := strconv.Atoi(locationParts[1])
	if err != nil {
		return fmt.Errorf("unable to convert port (%s)", locationParts[1])
	}

	// establish connection
	client, err := ssh.NewNativeClient(key, locationParts[0], "SSH-2.0-microbox", port, &nanPass)
	if err != nil {
		return fmt.Errorf("Failed to create new client - %s", err)
	}

	// printMOTD and warning
	display.MOTD()
	display.InfoProductionHost()

	// establish the ssh client connection and shell
	err = client.Shell()
	if err != nil && err.Error() != "exit status 255" {
		return fmt.Errorf("Failed to request shell - %s", err)
	}

	return nil
}
