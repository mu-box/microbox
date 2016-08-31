package dev

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	"github.com/nanobox-io/nanobox/models"
)

// Run ...
func Run(appModel *models.App) error {

	// load the start commands from the boxfile
	starts := loadStarts(appModel)

	// run the start commands in from the boxfile
	// in the dev container
	if err := runStarts(starts); err != nil {
		return err
	}

	// catch signals and stop the run commands on signal
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt)
	defer signal.Stop(sigs)

	for range sigs {
		// if we get a interupt we will jut return here
		// causing the container to be destroyed and our
		// exec processes to die
		return nil
	}

	return nil
}

func loadStarts(appModel *models.App) map[string]string {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	starts := map[string]string{}

	// loop through the nodes and get there start commands
	for _, node := range boxfile.Nodes("code") {

		values := boxfile.Node(node).Value("start")

		switch values.(type) {
		case string:
			starts[node] = values.(string)
		case []interface{}:
			// if it is an array we need the keys to be
			// web.site.2 where 2 is the index of the element
			for index, iFace := range values.([]interface{}) {
				if str, ok := iFace.(string); ok {
					starts[fmt.Sprintf("%s.%d", node, index)] = str
				}
			}
		case map[interface{}]interface{}:
			for key, val := range values.(map[interface{}]interface{}) {
				k, keyOk := key.(string)
				v, valOk := val.(string)
				if keyOk && valOk {
					starts[fmt.Sprintf("%s.%s", node, k)] = v
				}
			}
		}
	}
	return starts
}

func runStarts(starts map[string]string) error {
	// loop through the starts and run them in go routines
	for key, start := range starts {
		go runStart(key, start)
	}
	return nil
}

func runStart(name, command string) error {

	// create the docker command
	cmd := []string{
		"docker",
		"exec",
		"-u",
		"gonano",
		container_generator.DevName(),
		"/bin/bash",
		"-lc",
		fmt.Sprintf("cd /app/; %s", command),
	}

	lumber.Debug("run:runstarts: %+v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)

	// TODO: these will be replaced with something from the
	// new print library
	// we will also want to use 'name' to create some prefix
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	// run the process
	if err := process.Run(); err != nil && err.Error() != "exit status 137" {
		return err
	}

	return nil
}