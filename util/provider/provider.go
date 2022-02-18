package provider

import (
	"errors"
	"fmt"

	"github.com/mu-box/microbox/models"
)

// Provider ...
type Provider interface {
	BridgeRequired() bool
	Status() string
	IsReady() bool
	HostShareDir() string
	HostMntDir() string
	HostIP() (string, error)
	ReservedIPs() []string
	Valid() (error, []string)
	Create() error
	Reboot() error
	Stop() error
	Implode() error
	Destroy() error
	Start() error
	DockerEnv() error
	// we might be able to remove ip stuff as well
	AddIP(ip string) error
	RemoveIP(ip string) error
	SetDefaultIP(ip string) error
	// AddNat(host, container string) error
	// RemoveNat(host, container string) error
	RequiresMount() bool
	HasMount(mount string) bool
	AddMount(local, host string) error
	RemoveMount(local, host string) error
	RemoveEnvDir(id string) error
	Run(command []string) ([]byte, error)
}

var (
	providers = map[string]Provider{}
	verbose   = true
)

// Register ...
func Register(name string, p Provider) {
	providers[name] = p
}

// Display ...
func Display(verb bool) {
	verbose = verb
}

// Valid ...
func Valid() (error, []string) {

	p, err := fetchProvider()
	if err != nil {
		return fmt.Errorf("invalid provider - %s", err.Error()), []string{"invalid provider"}
	}

	return p.Valid()
}

func ValidReady() error {
	if !IsReady() {
		return errors.New("the provider is not ready try running 'microbox start' first")
	}
	return nil
}

// Status ...
func Status() string {

	p, err := fetchProvider()
	if err != nil {
		return "err: " + err.Error()
	}

	return p.Status()
}

// Create ...
func Create() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Create()
}

// Reboot ...
func Reboot() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Reboot()
}

// Stop ...
func Stop() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Stop()
}

// Implode ..
func Implode() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Implode()
}

// Destroy ..
func Destroy() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Destroy()
}

// Start ..
func Start() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.Start()
}

// HostShareDir ...
func HostShareDir() string {

	p, err := fetchProvider()
	if err != nil {
		return ""
	}

	return p.HostShareDir()
}

// HostMntDir ..
func HostMntDir() string {

	p, err := fetchProvider()
	if err != nil {
		return ""
	}

	return p.HostMntDir()
}

// HostIP ..
func HostIP() (string, error) {

	p, err := fetchProvider()
	if err != nil {
		return "", err
	}

	return p.HostIP()
}

// ReservedIPs ..
func ReservedIPs() []string {

	p, err := fetchProvider()
	if err != nil {
		return []string{}
	}

	return p.ReservedIPs()
}

// DockerEnv ..
func DockerEnv() error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.DockerEnv()
}

// AddIP ..
func AddIP(ip string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddIP(ip)
}

// RemoveIP ...
func RemoveIP(ip string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveIP(ip)
}

// SetDefaultIP ...
func SetDefaultIP(ip string) error {
	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.SetDefaultIP(ip)
}

// // AddNat ..
// func AddNat(host, container string) error {

// 	p, err := fetchProvider()
// 	if err != nil {
// 		return err
// 	}

// 	return p.AddNat(host, container)
// }

// // RemoveNat ..
// func RemoveNat(host, container string) error {

// 	p, err := fetchProvider()
// 	if err != nil {
// 		return err
// 	}

// 	return p.RemoveNat(host, container)
// }

// RequiresMount ...
func RequiresMount() bool {

	p, err := fetchProvider()
	if err != nil {
		return false
	}

	return p.RequiresMount()
}

// HasMount ...
func HasMount(path string) bool {

	p, err := fetchProvider()
	if err != nil {
		return false
	}

	return p.HasMount(path)
}

// AddMount ...
func AddMount(local, host string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.AddMount(local, host)
}

// RemoveMount ...
func RemoveMount(local, host string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveMount(local, host)
}

// RemoveEnvDir ...
func RemoveEnvDir(id string) error {

	p, err := fetchProvider()
	if err != nil {
		return err
	}

	return p.RemoveEnvDir(id)
}

// Run a command inside of the provider context
func Run(command []string) ([]byte, error) {

	p, err := fetchProvider()
	if err != nil {
		return nil, err
	}

	return p.Run(command)
}

func IsReady() bool {

	p, err := fetchProvider()
	if err != nil {
		return false
	}

	return p.IsReady()
}

func BridgeRequired() bool {

	p, err := fetchProvider()
	if err != nil {
		return false
	}

	return p.BridgeRequired()
}

// fetchProvider fetches the registered provider from the configured name
func fetchProvider() (Provider, error) {
	p, ok := providers[Name()]
	if !ok {
		return nil, errors.New("invalid provider")
	}

	return p, nil
}

func Name() string {
	config, _ := models.LoadConfig()

	prov := config.Provider
	if prov == "docker_machine" {
		prov = "docker-machine"
	}
	// set the provider to the default if it is a bad input
	if prov != "docker-machine" && prov != "native" {
		prov = "docker-machine"
	}

	return prov
}
