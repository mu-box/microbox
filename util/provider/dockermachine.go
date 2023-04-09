package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"

	// "path/filepath"
	"io"
	"regexp"
	"runtime"
	"strings"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"

	// "github.com/mu-box/microbox/util/fileutil"
	"github.com/mu-box/microbox/util/vbox"
)

var (
	vboxManageCmd    = vbox.DetectVBoxManageCmd()
	dockerMachineCmd = "microbox-machine"
)

// DockerMachine ...
type DockerMachine struct{}

// init ...
func init() {
	Register("docker-machine", DockerMachine{})
}

// Valid ensures docker-machine is installed and available
func (machine DockerMachine) Valid() (error, []string) {
	var masterErr error // will hold a collection of any errors hit

	missingParts := []string{}

	// we install our own docker-machine so we dont need to check

	// do you have vbox manage?
	out, err := exec.Command(vboxManageCmd, "-v").CombinedOutput()
	if err != nil {
		// append rather than return here so that we can check everything before erroring.
		// otherwise a user could error multiple times, learning about one problem at a time.
		missingParts = append(missingParts, "vboxmanage")
		masterErr = fmt.Errorf("missing vboxmanage - %q", err.Error())
	}
	// check output for errors/instruction
	if bytes.Contains(out, []byte("WARNING")) {
		if masterErr == nil {
			masterErr = fmt.Errorf("vboxmanage gave warning - \n%s", out)
		} else {
			masterErr = fmt.Errorf("%s - vboxmanage gave warning - \n%s", masterErr.Error(), out)
		}
	}

	// return early if we are running native mounts
	config, _ := models.LoadConfig()

	if config.MountType == "native" {
		if len(missingParts) == 0 && masterErr == nil {
			return nil, missingParts
		}
		return masterErr, missingParts
	}

	// net share checking
	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("netstat", "-tln").CombinedOutput()
		if err != nil {
			missingParts = append(missingParts, "nfs-kernel-server")
			if masterErr == nil {
				masterErr = fmt.Errorf("failed to check for netfs - %s", err.Error())
			} else {
				masterErr = fmt.Errorf("%s - failed to check for netfs - %s", masterErr.Error(), err.Error())
			}
		} else if !bytes.Contains(out, []byte("2049")) {
			missingParts = append(missingParts, "nfs-kernel-server")
			if masterErr == nil {
				masterErr = fmt.Errorf("missing or not running netfs")
			} else {
				masterErr = fmt.Errorf("%s - missing or not running netfs", masterErr.Error())
			}
		}

		if err := exec.Command("exportfs").Run(); err != nil {
			missingParts = append(missingParts, "exportfs")
			masterErr = fmt.Errorf("%s - missing exportfs - %q", masterErr.Error(), err.Error())
		}

	case "darwin":
		// unixCheck()
		// if err := exec.Command("nfsd", "status").Run(); err != nil {
		// 	missingParts = append(missingParts, "nfsd")
		// }

	case "windows":

	}

	return masterErr, missingParts
}

func (machine DockerMachine) Status() string {
	cmd := exec.Command(dockerMachineCmd, "status", "microbox")
	output, _ := cmd.CombinedOutput()
	return strings.TrimSpace(string(output))
}

func (machine DockerMachine) BridgeRequired() bool {
	return true
}

// Create creates the docker-machine vm
func (machine DockerMachine) Create() error {

	//
	if machine.isCreated() {
		return nil
	}

	// make sure dockermachine is clean
	exec.Command(dockerMachineCmd, "rm", "-f", "microbox").Run()

	display.ProviderSetup()

	// load the configuration for docker-machine
	conf, _ := models.LoadConfig()

	// load the cpus setting
	cpus := conf.CPUs
	if cpus < 1 {
		cpus = 1
	}

	// load the ram setting
	ram := conf.RAM
	if ram < 1 {
		ram = 1
	}

	// load in the disk size
	disk := conf.Disk

	cmd := []string{
		dockerMachineCmd,
		"create",
		"--driver",
		"virtualbox",
		"--virtualbox-boot2docker-url",
		"https://s3.amazonaws.com/tools.microbox.cloud/boot2docker/v1/boot2docker.iso",
		"--engine-env",
		fmt.Sprintf("HTTP_PROXY=%s", os.Getenv("HTTP_PROXY")),
		"--engine-env",
		fmt.Sprintf("HTTPS_PROXY=%s", os.Getenv("HTTPS_PROXY")),
		"--engine-env",
		fmt.Sprintf("NO_PROXY=%s", os.Getenv("NO_PROXY")),
		"--engine-env",
		fmt.Sprintf("http_proxy=%s", os.Getenv("http_proxy")),
		"--engine-env",
		fmt.Sprintf("https_proxy=%s", os.Getenv("https_proxy")),
		"--engine-env",
		fmt.Sprintf("no_proxy=%s", os.Getenv("no_proxy")),
		"--virtualbox-cpu-count",
		fmt.Sprintf("%d", cpus),
		"--virtualbox-memory",
		fmt.Sprintf("%d", ram*1024),
	}

	// append the disk if they set it big enough
	if disk >= 15360 {
		cmd = append(cmd, "--virtualbox-disk-size", fmt.Sprintf("%d", disk))
	}

	cmd = append(cmd, "microbox")

	process := exec.Command(cmd[0], cmd[1:]...)

	// lets try getting the extra bytes from creating
	fullBuffer := &bytes.Buffer{}
	writer := io.MultiWriter(display.NewStreamer("info"), fullBuffer)

	process.Stdout = writer
	process.Stderr = writer

	display.StartTask("Launching VM")

	if err := process.Run(); err != nil {
		// if its complaining about the world writable error issue
		if strings.Contains(fullBuffer.String(), "VERR_SUPLIB_WORLD_WRITABLE") {
			display.WorldWritable()
			return util.ErrorfQuiet("[USER] %s: %s", fullBuffer.String(), err)
		}
		display.ErrorTask()
		return util.Errorf("%s: %s", fullBuffer.String(), err)
	}

	display.StopTask()

	return nil
}

// Reboot reboots the docker-machine vm
func (machine DockerMachine) Reboot() error {

	if err := machine.Stop(); err != nil {
		return err
	}

	return machine.Start()
}

// Stop stops the docker-machine vm
func (machine DockerMachine) Stop() error {

	//
	if !machine.IsReady() {
		return nil
	}

	cmd := []string{
		dockerMachineCmd,
		"stop",
		"microbox",
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	process.Stdout = display.NewStreamer("info")
	process.Stderr = display.NewStreamer("info")

	display.StartTask("Shutting down VM")

	if err := process.Run(); err != nil {
		display.ErrorTask()
		return nil
	}

	display.StopTask()

	return nil
}

// imploding the docker-machine provider
// is the same as destroying it
func (machine DockerMachine) Implode() error {
	return Destroy()
}

// Destroy destroys the docker-machine vm
func (machine DockerMachine) Destroy() error {

	if !machine.isCreated() {
		return nil
	}

	cmd := []string{
		dockerMachineCmd,
		"rm",
		"-f",
		"microbox",
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	process.Stdout = display.NewStreamer("info")
	process.Stderr = display.NewStreamer("info")

	display.StartTask("Destroying VM")

	if err := process.Run(); err != nil {
		display.ErrorTask()
		return nil
	}

	display.StopTask()

	return nil
}

// Start starts and bootstraps docker-machine vm
func (machine DockerMachine) Start() error {

	// start the docker-machine vm
	if !machine.IsReady() {

		cmd := []string{
			dockerMachineCmd,
			"start",
			"microbox",
		}

		process := exec.Command(cmd[0], cmd[1:]...)

		// lets try getting the extra bytes from creating
		fullBuffer := &bytes.Buffer{}
		writer := io.MultiWriter(display.NewStreamer("info"), fullBuffer)

		process.Stdout = writer
		process.Stderr = writer

		display.StartTask("Booting VM")

		if err := process.Run(); err != nil {
			display.ErrorTask()
			return util.Errorf("%s: %s", fullBuffer.String(), err)
		}

		display.StopTask()
	}

	// create custom microbox docker network
	if !machine.hasNetwork() {
		config, _ := models.LoadConfig()
		ip, ipNet, err := net.ParseCIDR(config.DockerMachineNetworkSpace)
		if err != nil {
			return err
		}

		cmd := []string{
			dockerMachineCmd,
			"ssh",
			"microbox",
			"docker",
			"network",
			"create",
			"--driver=bridge",
			fmt.Sprintf("--subnet=%s", ipNet.String()),
			"--opt='com.docker.network.driver.mtu=1450'",
			"--opt='com.docker.network.bridge.name=redd0'",
			fmt.Sprintf("--gateway=%s", ip.String()),
			"microbox",
		}

		process := exec.Command(cmd[0], cmd[1:]...)

		// lets try getting the extra bytes from creating
		fullBuffer := &bytes.Buffer{}
		writer := io.MultiWriter(display.NewStreamer("info"), fullBuffer)

		process.Stdout = writer
		process.Stderr = writer

		display.StartTask("Configuring Network")

		if err := process.Run(); err != nil {
			display.ErrorTask()
			return util.Errorf("%s: %s", fullBuffer.String(), err)
		}

		display.StopTask()
	}

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"modprobe",
		"ip_vs",
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	// lets try getting the extra bytes from creating
	fullBuffer := &bytes.Buffer{}
	writer := io.MultiWriter(display.NewStreamer("info"), fullBuffer)

	process.Stdout = writer
	process.Stderr = writer

	display.StartTask("Loading kernel modules")

	if err := process.Run(); err != nil {
		display.ErrorTask()
		return util.Errorf("%s: %s", fullBuffer.String(), err)
	}

	display.StopTask()

	display.StartTask("Cleaning DHCP")

	// kill dhcp
	cmd = []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"pkill",
		"udhcpc",
	}

	process = exec.Command(cmd[0], cmd[1:]...)

	process.Stdout = display.NewStreamer("info")
	process.Stderr = display.NewStreamer("info")
	process.Run()

	display.StopTask()

	// check to see if the vm can talk to the host
	cmd = []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"curl",
		"--connect-timeout",
		"5",
		"192.168.99.1:65000",
	}

	process = exec.Command(cmd[0], cmd[1:]...)
	if err := process.Run(); err != nil {
		display.ErrorTask()
		display.VMCommunicationError()
		return util.ErrorfQuiet("[MICROBOX] VM cannot communicate with host")
	}

	if machine.changedIP() {
		return machine.regenerateCert()
	}

	return nil
}

// HostShareDir ...
func (machine DockerMachine) HostShareDir() string {
	return "/share/"
}

// HostMntDir ...
func (machine DockerMachine) HostMntDir() string {
	return "/mnt/sda1/"
}

// DockerEnv exports the docker connection information to the running process
func (machine DockerMachine) DockerEnv() error {
	// docker-machine env microbox
	// export DOCKER_TLS_VERIFY="1"
	// export DOCKER_HOST="tcp://192.168.99.102:2376"
	// export DOCKER_CERT_PATH="/Users/lyon/.docker/machine/machines/microbox"
	// export DOCKER_MACHINE_NAME="microbox"

	// create an anonymous struct that we will populate after running inspect
	inspect := struct {
		Driver struct {
			IPAddress string
		}
		HostOptions struct {
			EngineOptions struct {
				TLSVerify bool
			}
			AuthOptions struct {
				StorePath string
			}
		}
	}{}

	// fetch the docker-machine endpoint information
	cmd := exec.Command(dockerMachineCmd, "inspect", "microbox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// marshal the json output into the anonymous struct as defined above
	err = json.Unmarshal(b, &inspect)
	if err != nil {
		return err
	}

	// export TLS verify if set
	if inspect.HostOptions.EngineOptions.TLSVerify {
		os.Setenv("DOCKER_TLS_VERIFY", "1")
	}

	if inspect.Driver.IPAddress == "" {
		return util.Err{Message: "docker-machine didnt start docker properly", Suggest: "run `microbox-machine stop microbox` and try again"}
	}
	// set docker environment variables for client connections
	os.Setenv("DOCKER_MACHINE_NAME", "microbox")
	os.Setenv("DOCKER_HOST", fmt.Sprintf("tcp://%s:2376", inspect.Driver.IPAddress))
	os.Setenv("DOCKER_CERT_PATH", inspect.HostOptions.AuthOptions.StorePath)

	return nil
}

// AddIP adds an IP into the docker-machine vm for host access
func (machine DockerMachine) AddIP(ip string) error {
	if machine.hasIP(ip) {
		return nil
	}

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"ip",
		"addr",
		"add",
		ip,
		"dev",
		"eth1",
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

// RemoveIP removes an IP from the docker-machine vm
func (machine DockerMachine) RemoveIP(ip string) error {

	if !machine.hasIP(ip) {
		return nil
	}

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"ip",
		"addr",
		"del",
		ip,
		"dev",
		"eth1",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

func (dockermachine DockerMachine) SetDefaultIP(ip string) error {
	// check to see if the route is already set as the default
	checkOut, _ := exec.Command(dockerMachineCmd, "ssh", "microbox", "ip", "route").CombinedOutput()
	if strings.Contains(string(checkOut), ip) {
		// if the output has this ip already we can return without changing anything
		return nil
	}

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"ip",
		"route",
		"change",
		"192.168.99.0/24",
		"dev",
		"eth1",
		"src",
		ip,
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

// AddNat adds a nat to make an container accessible to the host network stack
func (machine DockerMachine) AddNat(ip, containerIP string) error {

	// add iptables prerouting rule
	if !machine.hasNatPreroute(ip, containerIP) {

		cmd := []string{
			dockerMachineCmd,
			"ssh",
			"microbox",
			"sudo",
			"/usr/local/sbin/iptables",
			"-t",
			"nat",
			"-A",
			"PREROUTING",
			"-d",
			ip,
			"-j",
			"DNAT",
			"--to-destination",
			containerIP,
		}

		process := exec.Command(cmd[0], cmd[1:]...)
		b, err := process.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", b, err)
		}

		// todo: check output for failures

	}

	// add iptables postrouting rule
	if !machine.hasNatPostroute(ip, containerIP) {

		cmd := []string{
			dockerMachineCmd,
			"ssh",
			"microbox",
			"sudo",
			"/usr/local/sbin/iptables",
			"-t",
			"nat",
			"-A",
			"POSTROUTING",
			"-s",
			containerIP,
			"-j",
			"SNAT",
			"--to-source",
			ip,
		}

		process := exec.Command(cmd[0], cmd[1:]...)
		b, err := process.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", b, err)
		}

		// todo: check output for failures

	}

	return nil
}

// RemoveNat removes nat from making a container inaccessible to the host network stack
func (machine DockerMachine) RemoveNat(ip, containerIP string) error {

	// remove iptables prerouting rule
	if machine.hasNatPreroute(ip, containerIP) {

		cmd := []string{
			dockerMachineCmd,
			"ssh",
			"microbox",
			"sudo",
			"/usr/local/sbin/iptables",
			"-t",
			"nat",
			"-D",
			"PREROUTING",
			"-d",
			ip,
			"-j",
			"DNAT",
			"--to-destination",
			containerIP,
		}

		process := exec.Command(cmd[0], cmd[1:]...)
		b, err := process.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", b, err)
		}

		// todo: check output for failures

	}

	// remove iptables postrouting rule
	if machine.hasNatPostroute(ip, containerIP) {

		cmd := []string{
			dockerMachineCmd,
			"ssh",
			"microbox",
			"sudo",
			"/usr/local/sbin/iptables",
			"-t",
			"nat",
			"-D",
			"POSTROUTING",
			"-s",
			containerIP,
			"-j",
			"SNAT",
			"--to-source",
			ip,
		}

		process := exec.Command(cmd[0], cmd[1:]...)
		b, err := process.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", b, err)
		}

		// todo: check output for failures

	}

	return nil
}

//
func (machine DockerMachine) RemoveEnvDir(id string) error {
	if id == "" {
		return fmt.Errorf("invalid env id")
	}

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"rm",
		"-rf",
		machine.HostMntDir() + id,
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	return nil
}

// HostIP inspects docker-machine to return the IP address of the vm
func (machine DockerMachine) HostIP() (string, error) {
	// create an anonymous struct that we will populate after running inspect
	inspect := struct {
		Driver struct {
			IPAddress string
		}
	}{}

	// fetch the docker-machine endpoint information
	cmd := exec.Command(dockerMachineCmd, "inspect", "microbox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", b, err)
	}

	// marshal the json output into the anonymous struct as defined above
	err = json.Unmarshal(b, &inspect)
	if err != nil {
		return "", err
	}

	return inspect.Driver.IPAddress, nil
}

func (machine DockerMachine) ReservedIPs() (rtn []string) {
	for i := 0; i < 20; i++ {
		rtn = append(rtn, fmt.Sprintf("192.168.99.%d", 100+i))
	}
	return
}

// Run a command in the vm
func (machine DockerMachine) Run(command []string) ([]byte, error) {

	// All commands need to be run in the docker machine, so we create a prefix
	context := []string{dockerMachineCmd, "ssh", "microbox"}

	// now we can generate a run command combining the context with the command
	run := append(context, command...)

	// when we actually run the command, we need to pop off the first item
	cmd := exec.Command(run[0], run[1:]...)

	// run the command and return the output
	return cmd.CombinedOutput()
}

// regenerate the certificates
// this should be used when the machine starts up with an ip that
// is different then last time
func (machine DockerMachine) regenerateCert() error {
	// fetch the docker-machine endpoint information

	display.StartTask("Regenerating Docker certs")

	cmd := exec.Command(dockerMachineCmd, "regenerate-certs", "-f", "microbox")
	b, err := cmd.CombinedOutput()
	if err != nil {
		display.ErrorTask()
		return fmt.Errorf("%s: %s", b, err)
	}

	display.StopTask()
	return nil
}

// isCreated ...
func (machine DockerMachine) isCreated() bool {
	// docker-machine status microbox
	cmd := exec.Command(dockerMachineCmd, "status", "microbox")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false
	}

	if bytes.Contains(output, []byte("Host does not exist: \"microbox\"")) ||
		bytes.Contains(output, []byte("error")) {
		return false
	}

	return true
}

// hasNetwork ...
func (machine DockerMachine) hasNetwork() bool {

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"docker",
		"network",
		"inspect",
		"microbox",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()

	if err != nil {
		return false
	}

	if bytes.Contains(output, []byte("Error: No such network: microbox")) {
		return false
	}

	return true
}

// IsReady ...
func (machine DockerMachine) IsReady() bool {

	// docker-machine status microbox
	cmd := exec.Command(dockerMachineCmd, "status", "microbox")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	matched, regerr := regexp.Match("Running", output)
	if regerr != nil {
		return false
	}

	return matched
}

// hasIP ...
func (machine DockerMachine) hasIP(ip string) bool {

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"ip",
		"addr",
		"show",
		"dev",
		"eth1",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	if err != nil {
		return false
	}

	matched, regerr := regexp.Match(ip, output)
	if regerr != nil {
		return false
	}

	return matched
}

// hasNatPreroute ...
func (machine DockerMachine) hasNatPreroute(hostIP, containerIP string) bool {

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"/usr/local/sbin/iptables",
		"-t",
		"nat",
		"-C",
		"PREROUTING",
		"-d",
		hostIP,
		"-j",
		"DNAT",
		"--to-destination",
		containerIP,
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()

	if err != nil {
		return false
	}

	if bytes.Contains(output, []byte("No chain/target/match by that name.")) {
		return false
	}

	return true
}

// hasNatPostroute
func (machine DockerMachine) hasNatPostroute(hostIP, containerIP string) bool {

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"microbox",
		"sudo",
		"/usr/local/sbin/iptables",
		"-t",
		"nat",
		"-C",
		"POSTROUTING",
		"-s",
		containerIP,
		"-j",
		"SNAT",
		"--to-source",
		hostIP,
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()

	if err != nil {
		return false
	}

	if bytes.Contains(output, []byte("No chain/target/match by that name.")) {
		return false
	}

	return true
}

func (machine DockerMachine) changedIP() bool {
	// get the previous host ip
	provider, err := models.LoadProvider()
	if err != nil {
		return true
	}

	// if it was never set the it cant have changed
	if provider.HostIP == "" {
		return false
	}

	// get the new host ip
	newIP, err := machine.HostIP()
	if err != nil {
		return true
	}

	// if the host ip has changed i need to update the database
	defer func() {
		if provider.HostIP != newIP {
			provider.HostIP = newIP
			provider.Save()
		}
	}()

	return provider.HostIP != newIP
}

func dockerMachineURL() string {

	download := "https://github.com/docker/machine/releases/download/v0.8.1"

	switch runtime.GOOS {
	case "darwin":
		// temporarily replace the docker-machine version with a custom one until
		// docker fixes the issues created by Sierra
		// download = fmt.Sprintf("%s/docker-machine-Darwin-x86_64", download)
		download = "https://s3.amazonaws.com/tools.microbox.cloud/docker-machine/darwin/docker-machine"
	case "linux":
		download = fmt.Sprintf("%s/docker-machine-Linux-x86_64", download)
	case "windows":
		download = fmt.Sprintf("%s/docker-machine-Windows-x86_64.exe", download)
	}

	return download
}
