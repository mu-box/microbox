package containers

import (
	"testing"

	docker "github.com/mu-box/golang-docker-client"
)

func TestSetProxyVars(t *testing.T) {
	config := docker.ContainerConfig{
		Name: "test-container",
	}

	config.Env = []string{"thing=thang"}

	setProxyVars(&config)

	if config.Env[0] != "thing=thang" {
		t.Errorf("Failed to preserve prior envs!")
	}
}
