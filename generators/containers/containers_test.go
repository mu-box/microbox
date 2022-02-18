package containers_test

import (
	"net"
	"testing"

	"github.com/mu-box/microbox/generators/containers"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/dhcp"
)

func TestBuildConfig(t *testing.T) {
	result := containers.BuildConfig("imagename")
	if result.Image != "imagename" ||
		result.Name != containers.BuildName() {
		// TODO: add checks for the binds
		t.Errorf("bad credentials")
	}
}

func TestCompileConfig(t *testing.T) {
	result := containers.CompileConfig("imagename")
	if result.Image != "imagename" ||
		result.Name != containers.CompileName() {
		// TODO: add checks for the binds
		t.Errorf("bad results")
	}
}

func TestComponentConfig(t *testing.T) {
	componentModel := &models.Component{
		Image:      "imagename",
		InternalIP: "1.2.3.4",
		AppID:      "2",
		Name:       "name",
	}

	result := containers.ComponentConfig(componentModel)
	if result.Image != "imagename" ||
		result.IP != "1.2.3.4" ||
		result.Name != "microbox_2_name" {
		t.Errorf("bad results")
	}
}

func TestPublishConfig(t *testing.T) {
	result := containers.PublishConfig("imagename")
	if result.Image != "imagename" ||
		result.Name != containers.PublishName() {
		// TODO: add checks for the binds
		t.Errorf("bad results")
	}
}

func TestDevConfig(t *testing.T) {
	appModel := &models.App{EnvID: "1", ID: "2"}
	result := containers.DevConfig(appModel)
	if result.Image != "mubox/build" ||
		result.Name != "microbox_2" {
		// TODO: add checks for the binds
		// TODO: add lib dir check
		t.Errorf("bad results")
	}
	dhcp.ReturnIP(net.ParseIP(result.IP))
}
