package util

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
)

func OsDetect() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getDarwin()
	case "windows":
		return "windows", nil
	case "linux":
		return "linux", nil
	}

	return "", fmt.Errorf("Unsupported operating system (%s). Please contact support.", runtime.GOOS)
}

func getDarwin() (string, error) {
	out, err := exec.Command("/usr/bin/sw_vers", "-productVersion").Output()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve version - %s", err.Error())
	}

	r, _ := regexp.Compile(`10\.([0-9]+).*`)
	match := r.FindStringSubmatch(string(out))
	if len(match) == 2 {
		return tenToDarwin(match[1])
	}

	r, _ = regexp.Compile(`([0-9]+)\..*`)
	match = r.FindStringSubmatch(string(out))
	if len(match) == 2 {
		return laterToDarwin(match[1])
	}

	return "", fmt.Errorf("Failed to parse version from %s. Please contact support.", out)
}

func tenToDarwin(v string) (string, error) {
	switch v {
	case "12":
		return "sierra", nil
	case "13":
		return "high sierra", nil
	case "14":
		return "mojave", nil
	case "15":
		return "catalina", nil
	}

	return "", fmt.Errorf("Unsupported macOS version 10.%s. Please contact support.", v)
}

func laterToDarwin(v string) (string, error) {
	switch v {
	case "11":
		return "big sur", nil
	case "12":
		return "monterey", nil
	}

	return "", fmt.Errorf("Unsupported macOS version %s.x. Please contact support.", v)
}

func ArchDetect() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "amd64", nil
	case "arm64":
		return "arm64", nil
	case "arm":
		return "arm32", nil
	case "s390x":
		return "s390x", nil
	}

	return "", fmt.Errorf("Unsupported processor type %s. Please contact support.", runtime.GOARCH)
}

func OsArchDetect() (string, error) {
	var os, arch string

	if v, err := OsDetect(); err != nil {
		return "", err
	} else {
		os = v
	}
	if v, err := ArchDetect(); err != nil {
		return "", err
	} else {
		arch = v
	}

	osArch := fmt.Sprintf("%s/%s", os, arch)
	switch osArch {
	case "sierra/amd64":
		return "sierra/amd64", nil
	case "high sierra/amd64":
		return "high sierra/amd64", nil
	case "mojave/amd64":
		return "mojave/amd64", nil
	case "catalina/amd64":
		return "catalina/amd64", nil
	case "big sur/amd64":
		return "big sur/amd64", nil
	case "big sur/arm64":
		return "big sur/arm64", nil
	case "monterey/amd64":
		return "monterey/amd64", nil
	case "monterey/arm64":
		return "monterey/arm64", nil
	case "windows/amd64":
		return "windows/amd64", nil
	case "linux/amd64":
		return "linux/amd64", nil
	case "linux/arm64":
		return "linux/arm64", nil
	case "linux/arm32":
		return "linux/arm32", nil
	case "linux/s390x":
		return "linux/s390x", nil
	}

	return "", fmt.Errorf("Unsupported os/arch combo %s. Please contact support.", osArch)
}
