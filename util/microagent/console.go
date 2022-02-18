package microagent

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	syscall "github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"

	"github.com/mu-box/microbox/util/display"
)

// establishes a remote docker client on a production app
func Console(key, location string) error {
	// establish connection to microagent
	path := fmt.Sprintf("/exec?key=%s", key)
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return fmt.Errorf("failed to generate a request for microagent: %s", err.Error())
	}

	// printMOTD and warning
	display.MOTD()
	display.InfoProductionContainer()

	// connect to remote machine
	remoteConn, err := connect(req, location)
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	stdIn, stdOut, _ := term.StdStreams()

	// establish file descriptors
	stdInFD, isTerminal := term.GetFdInfo(stdIn)
	stdOutFD, _ := term.GetFdInfo(stdOut)

	// if we are using a term, lets upgrade it to RawMode
	if isTerminal {
		// handle all incoming os signals and act accordingly; default behavior is to
		// forward all signals to microbox server
		go monitor(stdOutFD, location, key)

		oldInState, err := term.SetRawTerminal(stdInFD)
		if err == nil {
			defer term.RestoreTerminal(stdInFD, oldInState)
		}

		oldOutState, err := term.SetRawTerminalOutput(stdOutFD)
		if err == nil {
			defer term.RestoreTerminal(stdOutFD, oldOutState)
		}
	}

	go io.Copy(remoteConn, os.Stdin)
	io.Copy(os.Stdout, remoteConn)

	return nil
}

// monitor ...
func monitor(stdOutFD uintptr, location, key string) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGWINCH)
	defer signal.Stop(sigs)

	// inform the server what the starting size is
	resizeTTY(stdOutFD, location, key)

	// resize the tty for any signals received
	for range sigs {
		resizeTTY(stdOutFD, location, key)
	}
}

func resizeTTY(fd uintptr, location, key string) error {

	ws, err := term.GetWinsize(fd)
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %s", err.Error())
	}

	// extract height and width
	w := int(ws.Width)
	h := int(ws.Height)

	// resize by posting to a path on microagent
	url := fmt.Sprintf("https://%s/resizeexec?key=%s&w=%d&h=%d", location, key, w, h)

	// POST
	if _, err := http.Post(url, "plain/text", nil); err != nil {
		return fmt.Errorf("Failed to resize the terminal: %s", err.Error())
	}

	return nil
}
