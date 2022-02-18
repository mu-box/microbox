package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util/display"
)

var (

	// LogoutCmd ...
	LogoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Remove your microbox.cloud api token from your local microbox client.",
		Long:  ``,
		Run:   logoutFn,
	}

	// loginCmdFlags ...
	logoutCmdFlags = struct {
		endpoint string
	}{}
)

func init() {
	LogoutCmd.Flags().StringVarP(&logoutCmdFlags.endpoint, "endpoint", "e", "", "endpoint")
}

// logoutFn ...
func logoutFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Logout(logoutCmdFlags.endpoint))
}
