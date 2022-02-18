package config

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	boxfile "github.com/mu-box/microbox-boxfile"

	"github.com/mu-box/microbox/util/fileutil"
)

// AppName ...
func AppName() string {

	// if no name is given use localDirName
	app := LocalDirName()

	// read boxfile and look for dev:name
	box := boxfile.NewFromPath(Boxfile())
	devName := box.Node("dev").StringValue("name")

	// set the app name
	if devName != "" {
		app = devName
	}

	return app
}

// EnvID ...
func EnvID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(LocalDir())))
}

// MicroboxPath ...
func MicroboxPath() string {

	programName, err := os.Executable()
	if err == nil {
		return programName
	}

	// lookup the full path to microbox
	path, err := exec.LookPath(os.Args[0])
	if err == nil {
		return path
	}

	// if args[0] was a path to microbox already
	if fileutil.Exists(programName) {
		return programName
	}

	// unable to find the full path, just return what was called
	return programName
}

// the path where the vpn is located
func VpnPath() string {
	bridgeClient := "microbox-vpn"

	// lookup the full path to microbox
	path, err := exec.LookPath(bridgeClient)
	if err == nil {
		return path
	}

	cmd := filepath.Join(BinDir(), bridgeClient)

	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf(`%s\%s.exe`, BinDir(), bridgeClient)
	}

	return cmd
}
