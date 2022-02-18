package config

import (
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestGlobalDir(t *testing.T) {
	dir := GlobalDir()
	if !strings.HasSuffix(dir, ".microbox") {
		t.Errorf("missing microbox suffix")
	}
}

func TestLocalDir(t *testing.T) {
	dir := LocalDir()
	// this is 'microbox', because the boxfile is at the root level. localdir returns
	// a parent boxfile if none is found in the current directory
	if !strings.HasSuffix(dir, "microbox") {
		t.Errorf("local dir mismatch '%s'", dir)
	}
}

func TestLocalDirName(t *testing.T) {
	dir := LocalDirName()
	if dir != "microbox" {
		t.Errorf("local dir name mismatch '%s'", dir)
	}
}

func TestSSHDir(t *testing.T) {
	homedir, _ := homedir.Dir()
	if SSHDir() != homedir+"/.ssh" {
		t.Errorf("incorrect ssh directory")
	}
}
