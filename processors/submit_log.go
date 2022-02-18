package processors

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
	"github.com/mu-box/microbox/util/odin"
)

func SubmitLog(args string) error {
	// if we are running as privilege we dont submit
	if util.IsPrivileged() {
		return nil
	}

	auth, _ := models.LoadAuth()
	conf, _ := models.LoadConfig()

	// if we are in ci mode or we are setting a configuration
	// leave here
	if strings.Contains(args, "login") || strings.Contains(args, "config") || conf.CIMode {
		return nil
	}

	if auth.Key == "" && !conf.Anonymous {
		display.LoginRequired()
		err := Login("", "", "")
		if err != nil {
			return err
		}
	}

	app := ""

	env, err := models.FindEnvByID(config.EnvID())
	if strings.Contains(args, "deploy") || strings.Contains(args, "tunnel") || strings.Contains(args, "console") {
		if err == nil {
			remote, ok := env.Remotes["default"]
			if ok {
				app = remote.ID
			}
		}
	}

	// tell microbox
	go odin.SubmitEvent(
		fmt.Sprintf("desktop/%s", args),
		fmt.Sprintf("desktop command: microbox %s", args),
		app,
		map[string]interface{}{
			"os":         runtime.GOOS,
			"provider":   conf.Provider,
			"mount-type": conf.MountType,
			"boxfile":    env.UserBoxfile,
		},
	)

	return nil
}
