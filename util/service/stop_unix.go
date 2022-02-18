//go:build !windows
// +build !windows

package service

import (
	"fmt"
	"os/exec"
)

func Stop(name string) error {
	cmd := stopCmd(name)
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("out: %s, err: %s", out, err)
	}
	return nil
}
