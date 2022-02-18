package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/processors/env/share"
)

func (machine DockerMachine) RequiresMount() bool {
	return true
}

// HasMount checks to see if the mount exists in the vm
func (machine DockerMachine) HasMount(mount string) bool {

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"cat",
		"/proc/mounts",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	if err != nil {
		return false

	}
	return strings.Contains(string(output), mount) && !machine.staleMount(mount)
}

func (machine DockerMachine) staleMount(mount string) bool {

	// removed the ensure statement because its slow and already installed
	// // ensure stat is installed
	// cmd := []string{"sh", "-c", setupCoreUtilsScript()}
	// if b, err := Run(cmd); err != nil {
	// 	lumber.Error("stat install output: %s", b)
	// 	return true
	// }

	output, err := Run([]string{"stat", mount})
	if err != nil {
		if strings.Contains(string(output), "No such file or directory") {
			return false
		}
		return true
	}

	return strings.Contains(string(output), "Stale") || strings.Contains(string(output), "stale")
}

// AddMount adds a virtualbox mount into the docker-machine vm
func (machine DockerMachine) AddMount(local, host string) error {

	// stop early if already mounted
	if machine.HasMount(host) {
		lumber.Info("mount exists for %s", host)
		return nil
	}

	if machine.staleMount(host) {
		lumber.Debug("Removing stale mount: %s", host)
		if err := machine.removeNativeMount(local, host); err != nil {
			return fmt.Errorf("failed to clean up stale mount: %s", err)
		}
	}

	config, _ := models.LoadConfig()
	switch config.MountType {

	case "netfs":

		// add netfs share
		// here we use the processor so we can do privilege exec
		lumber.Info("adding share in netfs for %s", local)
		if err := share.Add(local); err != nil {
			return err
		}
		// add netfs mount

		lumber.Info("adding mount in for netfs %s -> %s", local, host)
		if err := machine.addNetfsMount(local, host); err != nil {
			return err
		}

	default:

		// add share
		lumber.Info("adding share for native %s", local)
		if err := machine.addShare(local, host); err != nil {
			return err
		}

		// add mount
		lumber.Info("adding mount for  native %s -> %s", local, host)
		if err := machine.addNativeMount(local, host); err != nil {
			return err
		}
	}

	return nil
}

// RemoveMount removes a mount from the docker-machine vm
func (machine DockerMachine) RemoveMount(local, host string) error {
	if machine.HasMount(host) {
		// all mounts are removed as if they are native
		if err := machine.removeNativeMount(local, host); err != nil {
			return err
		}
	}

	// if we are supposed to keep the shares return here
	if registry.GetBool("keep-share") {
		return nil
	}

	// remove any netfs shares
	if err := share.Remove(local); err != nil {
		return err
	}
	// remove any native shares
	if err := machine.removeShare(local, host); err != nil {
		return err
	}

	return nil
}

// hasShare checks to see if the share exists
func (machine DockerMachine) hasShare(local, host string) bool {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	cmd := []string{
		vboxManageCmd,
		"showvminfo",
		"microbox",
		"--machinereadable",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), name)
}

// AddShare adds the provided path as a shareable filesystem
func (machine DockerMachine) addShare(local, host string) error {

	if machine.hasShare(local, host) {
		return nil
	}

	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	cmd := []string{
		vboxManageCmd,
		"sharedfolder",
		"add",
		"microbox",
		"--name",
		name,
		"--hostpath",
		local,
		"--transient",
	}

	lumber.Info("add share native cmd: %v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures
	return nil
}

// RemoveShare removes the provided path as a shareable filesystem; we don't care
// what the user has configured, we need to remove any shares that may have been
// setup previously
func (machine DockerMachine) removeShare(local, host string) error {

	if !machine.hasShare(local, host) {
		return nil
	}

	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	cmd := []string{
		vboxManageCmd,
		"sharedfolder",
		"remove",
		"microbox",
		"--name",
		name,
		"--transient",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

func (machine DockerMachine) addNativeMount(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	// create folder

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"mkdir",
		"-p",
		host,
	}

	lumber.Info("add mount native cmd (mkdir): %v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	lumber.Info("mkdir data: %s", b)
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// mount
	cmd = []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"mount",
		"-t",
		"vboxsf",
		"-o",
		"uid=1000,gid=1000",
		name,
		host,
	}

	lumber.Info("add mount native cmd: %v", cmd)
	process = exec.Command(cmd[0], cmd[1:]...)
	b, err = process.CombinedOutput()
	lumber.Info("mount data: %s", b)
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

func (machine DockerMachine) removeNativeMount(local, host string) error {
	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"umount",
		"-f",
		host,
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

// setupCoreUtilsScript returns a string containing the script to setup cifs
func setupCoreUtilsScript() string {
	script := `
		if [ ! -f /usr/local/bin/stat ]; then
			wget -O /mnt/sda1/tmp/tce/optional/coreutils.tcz http://repo.tinycorelinux.net/7.x/x86_64/tcz/coreutils.tcz &&

			tce-load -i coreutils;
		fi
	`

	return strings.Replace(script, "\n", "", -1)
}
