package service

import (
	"bytes"
	"fmt"
	"os/exec"
)

func serviceConfigFile(name string) string {
	fmtString := ""
	switch launchSystem() {
	case "systemd":
		fmtString = "/etc/systemd/system/%s.service"
	case "upstart":
		fmtString = "/etc/init/%s.conf"
	}
	return fmt.Sprintf(fmtString, name)
}

func launchSystem() string {
	_, err := exec.LookPath("systemctl")
	if err == nil {
		return "systemd"
	}

	_, err = exec.LookPath("initctl")
	if err == nil {
		return "upstart"
	}

	return ""
}

func startCmd(name string) []string {
	switch launchSystem() {
	case "systemd":
		// systemctl start microbox-openvpn.service
		return []string{"systemctl", "start", fmt.Sprintf("%s.service", name)}
	case "upstart":
		// initctl start microbox-openvpn
		return []string{"initctl", "start", name}
	}

	return nil
}

func Running(name string) bool {
	switch launchSystem() {
	case "systemd":
		out, err := exec.Command("systemctl", "--no-pager", "status", name).CombinedOutput()
		if err != nil {
			return false
		}

		if !bytes.Contains(out, []byte("running")) {
			return false
		}
		return true
	case "upstart":
		out, err := exec.Command("initctl", "status", name).CombinedOutput()
		if err != nil {
			return false
		}

		if !bytes.Contains(out, []byte("running")) {
			return false
		}
		return true
	}

	return false
}

func stopCmd(name string) []string {
	switch launchSystem() {
	case "systemd":
		// systemctl start microbox-openvpn.service
		return []string{"systemctl", "stop", fmt.Sprintf("%s.service", name)}
	case "upstart":
		// initctl start microbox-openvpn
		return []string{"initctl", "stop", name}
	}

	return nil
}

func removeCmd(name string) []string {
	return nil
}
