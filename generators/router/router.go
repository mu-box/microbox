package router

import (
	boxfile "github.com/mu-box/microbox-boxfile"
	"github.com/mu-box/microbox/models"
)

//
func loadBoxfile(appModel *models.App) boxfile.Boxfile {
	return boxfile.New([]byte(appModel.DeployedBoxfile))
}
