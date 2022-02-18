package mixpanel

import (
	"runtime"

	"github.com/jcelliott/lumber"
	mp "github.com/timehop/go-mixpanel"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
)

var token string

func Report(args string) {

	go func() {
		config, _ := models.LoadConfig()
		// do not include ci in the mixpanel results
		if config.CIMode {
			return
		}
		mx := mp.NewMixpanel(token)
		id := util.UniqueID()

		err := mx.Track(id, "command", mp.Properties{
			"os":         runtime.GOOS,
			"provider":   config.Provider,
			"mount-type": config.MountType,
			"args":       args,
			"cpus":       runtime.NumCPU(),
		})

		if err != nil {
			lumber.Error("mixpanel(%s).Report(%s): %s", token, args, err)
		}
	}()
}
