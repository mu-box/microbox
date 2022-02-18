//go:build !windows
// +build !windows

package service

import (
	"fmt"
	"os"
	"os/exec"
)

func Remove(name string) error {
	if len(removeCmd(name)) != 0 {
		cmd := removeCmd(name)
		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("out: %s, err: %s", out, err)
		}

	}

	os.Remove(serviceConfigFile(name))
	return nil
}
