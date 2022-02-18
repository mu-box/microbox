// Package commands defines the commands that microbox can run
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/commands/server"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/update"
)

var (
	// debug mode
	debugMode bool

	// display level debug
	displayDebugMode bool

	// display level trace
	displayTraceMode bool

	internalCommand bool
	showVersion     bool
	endpoint        string

	// MicroboxCmd ...
	MicroboxCmd = &cobra.Command{
		Use:   "microbox",
		Short: "",
		Long:  ``,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// report the command to microbox
			processors.SubmitLog(strings.Replace(ccmd.CommandPath(), "microbox ", "", 1))
			// mixpanel.Report(strings.Replace(ccmd.CommandPath(), "microbox ", "", 1))

			registry.Set("debug", debugMode)

			// setup the display output
			if displayDebugMode {
				lumber.Level(lumber.DEBUG)
				display.Summary = false
				display.Level = "debug"
			}

			if displayTraceMode {
				lumber.Level(lumber.TRACE)
				display.Summary = false
				display.Level = "trace"
			}

			// alert the user if an update is needed
			update.Check()

			if _, err := util.OsArchDetect(); err != nil {
				fmt.Println(err)
				os.Exit(126) // Command invoked cannot execute
			}

			configModel, _ := models.LoadConfig()

			// TODO: look into global messaging
			if internalCommand {
				registry.Set("internal", internalCommand)
				// setup a file logger, this will be replaced in verbose mode.
				fileLogger, _ := lumber.NewAppendLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "microbox.log")))
				lumber.SetLogger(fileLogger)

			} else {
				// We should only allow admin in 3 cases
				// 1 cimode
				// 2 server is running
				// 3 configuring
				fullCmd := strings.Join(os.Args, " ")
				if util.IsPrivileged() &&
					!configModel.CIMode &&
					!strings.Contains(fullCmd, "set ci") &&
					!strings.Contains(ccmd.CommandPath(), "server") {
					// if it is not an internal command (starting the server requires privileges)
					// we wont run microbox as privilege
					display.UnexpectedPrivilege()
					os.Exit(1)
				}
			}

			if endpoint != "" {
				registry.Set("endpoint", endpoint)
			}

			if configModel.CIMode {
				lumber.Level(lumber.INFO)
				display.Summary = false
				display.Level = "info"
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if displayDebugMode || showVersion {
				fmt.Println(models.VersionString())
				return
			}
			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

// init creates the list of available microbox commands and sub commands
func init() {

	// persistent flags
	MicroboxCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "production endpoint")
	MicroboxCmd.PersistentFlags().MarkHidden("endpoint")
	MicroboxCmd.PersistentFlags().BoolVarP(&internalCommand, "internal", "", false, "Skip pre-requisite checks")
	MicroboxCmd.PersistentFlags().MarkHidden("internal")
	MicroboxCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "", false, "In the event of a failure, drop into debug context")
	MicroboxCmd.PersistentFlags().BoolVarP(&displayDebugMode, "verbose", "v", false, "Increases display output and sets level to debug")
	MicroboxCmd.PersistentFlags().BoolVarP(&showVersion, "version", "", false, "Print version information and exit")
	MicroboxCmd.PersistentFlags().BoolVarP(&displayTraceMode, "trace", "t", false, "Increases display output and sets level to trace")

	// log specific flags
	LogCmd.Flags().BoolVarP(&logRaw, "raw", "r", false, "Print raw log timestamps instead")
	LogCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "Follow logs (live feed)")
	LogCmd.Flags().IntVarP(&logNumber, "number", "n", 0, "Number of historic logs to print")
	// todo:
	// LogCmd.Flags().StringVarP(&logStart, "start", "", "", "Timestamp of oldest historic log to print")
	// LogCmd.Flags().StringVarP(&logEnd, "end", "", "", "Timestamp of newest historic log to print")
	// LogCmd.Flags().StringVarP(&logLimit, "limit", "", "", "Time to limit amount of historic logs to print")

	// subcommands
	MicroboxCmd.AddCommand(ConfigureCmd)
	MicroboxCmd.AddCommand(RunCmd)
	MicroboxCmd.AddCommand(BuildCmd)
	MicroboxCmd.AddCommand(CompileCmd)
	MicroboxCmd.AddCommand(DeployCmd)
	MicroboxCmd.AddCommand(ConsoleCmd)
	MicroboxCmd.AddCommand(RemoteCmd)
	MicroboxCmd.AddCommand(StatusCmd)
	MicroboxCmd.AddCommand(LoginCmd)
	MicroboxCmd.AddCommand(LogoutCmd)
	MicroboxCmd.AddCommand(CleanCmd)
	MicroboxCmd.AddCommand(InfoCmd)
	MicroboxCmd.AddCommand(TunnelCmd)
	MicroboxCmd.AddCommand(ImplodeCmd)
	MicroboxCmd.AddCommand(DestroyCmd)
	MicroboxCmd.AddCommand(StartCmd)
	MicroboxCmd.AddCommand(StopCmd)
	MicroboxCmd.AddCommand(UpdateCmd)
	MicroboxCmd.AddCommand(EvarCmd)
	MicroboxCmd.AddCommand(DnsCmd)
	MicroboxCmd.AddCommand(LogCmd)
	MicroboxCmd.AddCommand(VersionCmd)
	MicroboxCmd.AddCommand(server.ServerCmd)

	// hidden subcommands
	MicroboxCmd.AddCommand(EnvCmd)
	MicroboxCmd.AddCommand(InspectCmd)
}
