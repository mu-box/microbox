package hookit

import (
	"fmt"
	"strings"

	"github.com/mu-box/microbox/commands/registry"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/console"
	"github.com/mu-box/microbox/util/display"
)

// since `--debug` runs the command twice, combined prevents duplication of the error output.
var combined bool

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {

	// display.Streamer is an io.Writer and can be passed to DockerExec
	stream := &display.Streamer{}

	if !combined {
		stream = display.NewStreamer(displayLevel)
	}

	stream.CaptureOutput(true)

	out, err := util.DockerExec(container, "root", "/opt/microbox/hooks/"+hook, []string{payload}, stream)
	if err != nil && (strings.Contains(string(out), "such file or directory") && strings.Contains(err.Error(), "bad exit code(126)")) {
		// if its a 126 the hook didnt exist
		return "", nil
	}

	outs := stream.Output()
	// these hooks depend on the output, we shouldn't append anything.
	if hook != "boxfile" && hook != "keys" {
		if out == "" {
			out = outs
		} else if outs != "" {
			out = fmt.Sprintf("%s --- %s", out, outs)
		}
	}

	if err != nil {
		// todo: add errorf's for other errors the hooks may stream
		if strings.Contains(outs, "INVALID BOXFILE") {
			return out, util.Errorf("[USER] invalid node in boxfile (see output for more detail)")
		}
		return out, util.Errorf("[HOOKS] failed to execute hook (%s) on %s: %s", hook, container, err)
	}
	return out, nil
}

func DebugExec(container, hook, payload, displayLevel string) (string, error) {
	res, err := Exec(container, hook, payload, displayLevel)

	// leave early if no error
	if err != nil {
		display.ErrorTask()
	}
	if err == nil || !registry.GetBool("debug") {
		return res, err
	}

	combined = true
	// todo: why run again if we are going to let them run it?
	res, err = Exec(container, hook, payload, displayLevel)

	fmt.Println()
	fmt.Printf("Failed to execute %s hook: %s\n", hook, err)
	fmt.Println("Entering Debug Mode")
	fmt.Printf("  container: %s\n", container)
	fmt.Printf("  hook:      %s\n", hook)
	fmt.Printf("  payload:   %s\n", payload)
	fmt.Println()

	err = console.Run(container, console.ConsoleConfig{})
	if err != nil {
		return res, fmt.Errorf("failed to establish a debug session: %s", err.Error())
	}
	combined = false

	// try running the exec one more time.
	return Exec(container, hook, payload, displayLevel)
}
