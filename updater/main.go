// Package main ...
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/service"
	"github.com/mu-box/microbox/util/update"
)

// microbox-update's version
var VERSION = "0.8.0" // todo: maybe we'll ldflag populate this too

// main ...
func main() {
	path := ""
	var err error

	// this will write a new binary at location provided `microbox-update newbinary`
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Printf("microbox-update %s\n", VERSION)
			return
		}
		path = os.Args[1]
	} else {
		// get the location of the current microbox
		path, err = exec.LookPath(update.Name)
		if err != nil {
			fmt.Printf("Cannot find %s: %s\n", update.Name, err)
			os.Exit(1)
		}
	}

	if !util.IsPrivileged() {

		if runtime.GOOS == "windows" {
			// re-run this command as the administrative user
			fmt.Println()
			fmt.Println("The update process requires Administrator privileges.")
			fmt.Println("Another window will be opened as the Administrator to continue this process.")

			// block here until the user hits enter. It's not ideal, but we need to make
			// sure they see the new window open.
			fmt.Println("Enter to continue:")
			var input string
			fmt.Scanln(&input)

		}

		// todo: make sure removing this doesn't break things
		// // make sure the .microbox folder is created by our user
		// models.LoadUpdate()

		// determine the full path to the executable in case
		// it isn't in the path when run with sudo
		cmdPath, err := os.Executable()
		if err != nil {
			fmt.Printf("Cannot find the path for microbox-update\n")
			os.Exit(1)
		}
		cmd := fmt.Sprintf("%s \"%s\"", cmdPath, path)
		if err := util.PrivilegeExec(cmd); err != nil {
			os.Exit(1)
		}

		// we're done
		return
	}

	// stop the microbox service so we can replace the microbox binary
	service.Stop("microbox-server")

	// run the update
	err = update.Run(path)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	if runtime.GOOS == "windows" {
		// The update process was spawned in a separate window, which will
		// close as soon as this command is finished. To ensure they see the
		// message, we need to hold open the process until they hit enter.
		fmt.Println()
		fmt.Println("Enter to continue:")
		var input string
		fmt.Scanln(&input)
	}
}
