package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/netfs"
)

var (

	// DevNetfsCmd ...
	DevNetfsCmd = &cobra.Command{
		Use:    "netfs",
		Short:  "add or remove netfs directories",
		Long:   ``,
		Hidden: true,
	}

	// DevNetfsAddCmd ...
	DevNetfsAddCmd = &cobra.Command{
		Use:    "add",
		Short:  "add a netfs export",
		Long:   ``,
		Hidden: true,
		Run:    devNetfsAddFunc,
	}

	// DevNetfsRmCmd ...
	DevNetfsRmCmd = &cobra.Command{
		Use:    "rm",
		Short:  "remove a netfs export",
		Long:   ``,
		Hidden: true,
		Run:    devNetfsRmFunc,
	}
)

func init() {
	DevNetfsCmd.AddCommand(DevNetfsAddCmd)
	DevNetfsCmd.AddCommand(DevNetfsRmCmd)
}

// devNetfsAddFunc will run the netfs function for adding a netfs export
func devNetfsAddFunc(ccmd *cobra.Command, args []string) {
	// validate that a path was provided
	host := args[0]
	path := args[1]

	// todo: error if path is nil

	netfs.Add(host, path)
}

// devNetfsRmFunc will run the netfs function for removing a netfs export
func devNetfsRmFunc(ccmd *cobra.Command, args []string) {
	// validate that a path was provided
	host := args[0]
	path := args[1]

	// todo: error if path is nil

	netfs.Remove(host, path)
}