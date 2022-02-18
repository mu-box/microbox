// Package dhcp ...
package dhcp

import (
	"errors"
	"net"
	"sync"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/locker"
	"github.com/mu-box/microbox/util/provider"
)

var (
	errIPNotFound = errors.New("Ip Not Found")
	mutex         = sync.Mutex{}
)

// IPSpace ...
type IPSpace struct {
	GlobalIP  net.IP
	GlobalNet net.IPNet
	LocalIP   net.IP
	LocalNet  net.IPNet
	NativeIP  net.IP
	NativeNet net.IPNet
}

// ReserveGlobal ...
func ReserveGlobal() (net.IP, error) {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	//
	ipSpace, err := getIPSpace()
	if err != nil {
		return nil, err
	}

	//
	reservedIPs, err := getReserved()
	if err != nil {
		return nil, err
	}

	// dump the first ip because it is the gateway
	ip := ipSpace.GlobalIP
	inc(ip)

	//
	for ; ipSpace.GlobalNet.Contains(ip); inc(ip) {
		if !contains(reservedIPs, ip) {
			err := setReserved(append(reservedIPs, ip))
			if err != nil {
				return nil, err
			}
			return ip, nil
		}
	}

	return nil, errIPNotFound
}

// Flush ...
func Flush() {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	// remove all the ip models
	ips := models.IPs{}
	ips.Delete()
}

func LocalNet() (*net.IPNet, error) {
	ipSpace, err := getIPSpace()
	if err != nil {
		return nil, err
	}
	// switch based on what provider we are using
	config, _ := models.LoadConfig()
	switch config.Provider {
	case "docker-machine":
		return &ipSpace.LocalNet, nil
	case "native":
		return &ipSpace.NativeNet, nil
	}
	return nil, errors.New("no network found")
}

// ReserveLocal ...
func ReserveLocal() (net.IP, error) {

	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	//
	ipSpace, err := getIPSpace()
	if err != nil {
		return nil, err
	}

	//
	reservedIPs, err := getReserved()
	if err != nil {
		return nil, err
	}

	// switch based on what provider we are using
	config, _ := models.LoadConfig()

	switch config.Provider {
	case "docker-machine":

		// dump the first ip because it is the gateway
		ip := ipSpace.LocalIP
		inc(ip)

		// get dockers local ipspace
		for ; ipSpace.LocalNet.Contains(ip); inc(ip) {
			if !contains(reservedIPs, ip) {
				err := setReserved(append(reservedIPs, ip))
				if err != nil {
					return nil, err
				}
				return ip, nil
			}
		}

	case "native":

		// dump the first ip because it is the gateway
		ip := ipSpace.NativeIP
		inc(ip)

		// get the native ipspace
		for ; ipSpace.NativeNet.Contains(ip); inc(ip) {
			if !contains(reservedIPs, ip) {
				err := setReserved(append(reservedIPs, ip))
				if err != nil {
					return nil, err
				}
				return ip, nil
			}
		}

	}

	return nil, errIPNotFound
}

// ReturnIP ...
func ReturnIP(ip net.IP) error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()
	mutex.Lock()
	defer mutex.Unlock()

	//
	reservedIPs, err := getReserved()
	if err != nil {
		return err
	}

	//
	for i, reservedIP := range reservedIPs {
		if reservedIP.Equal(ip) {
			return setReserved(append(reservedIPs[:i], reservedIPs[i+1:]...))
		}
	}

	return nil
}

// getIPSpace do not store the space on the disk.
func getIPSpace() (IPSpace, error) {
	ipSpace := IPSpace{}
	config, _ := models.LoadConfig()

	// there was no data stored for ip space so we need to populate it
	ip, ipNet, err := net.ParseCIDR(config.ExternalNetworkSpace)
	if err != nil {
		return ipSpace, err
	}
	ipSpace.GlobalIP = ip
	ipSpace.GlobalNet = *ipNet

	//
	ip, ipNet, err = net.ParseCIDR(config.DockerMachineNetworkSpace)
	if err != nil {
		return ipSpace, err
	}
	ipSpace.LocalIP = ip
	ipSpace.LocalNet = *ipNet

	//
	ip, ipNet, err = net.ParseCIDR(config.NativeNetworkSpace)
	if err != nil {
		return ipSpace, err
	}
	ipSpace.NativeIP = ip
	ipSpace.NativeNet = *ipNet
	return ipSpace, nil
}

// contains ...
func contains(ips []net.IP, ip net.IP) bool {
	// check against the ips in the data set
	for _, setIP := range ips {
		if setIP.Equal(ip) {
			return true
		}
	}

	// check against the ips the provider needs
	for _, providerIP := range provider.ReservedIPs() {
		if ip.String() == providerIP {
			return true
		}
	}

	return false
}

// getReserved ...
func getReserved() ([]net.IP, error) {
	ips, _ := models.LoadIPs()
	return []net.IP(ips), nil
}

// setReserved ...
func setReserved(ips []net.IP) error {
	mIPs := models.IPs(ips)
	return mIPs.Save()
}

// inc ...
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
