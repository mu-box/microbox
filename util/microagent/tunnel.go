package microagent

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"

	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/display"
)

func Tunnel(key, location, port, name string) error {
	if port == "0" {
		return fmt.Errorf("port 0 is not a valid port - that component has nothing to tunnel to")
	}

	// establish a connection and just leave it open.
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		return fmt.Errorf("failed to generate a request for microagent: %s", err.Error())
	}

	// set noproxy because this connection allows more multiple connections
	// to use the tunnel
	req.Header.Set("X-NOPROXY", "true")
	conn, err := connect(req, location)
	if err != nil {
		return err
	}
	defer conn.Close()

	if port == "22" {
		port = "2022"
	}

	// setup a tcp listener
	serv, err := net.Listen("tcp4", fmt.Sprintf(":%s", port))
	if err != nil {
		err2 := util.Err{
			Code:    "TUNNEL",
			Message: err.Error(),
		}
		if strings.Contains(err.Error(), "address already in use") || err == syscall.EADDRINUSE {
			display.PortInUse(port)
			err2.Code = "USER"
			err2.Suggest = fmt.Sprintf("It appears your local port (%s) is in use. Please specify a different port with '-p'.", port)
		}
		if strings.Contains(err.Error(), "bind: permission denied") || err == syscall.EACCES {
			display.PortPrivileged(port)
			err2.Code = "USER"
			err2.Suggest = fmt.Sprintf("It appears you don't have permission to use port '%s'. Please specify a different port with '-p'.", port)
		}

		return util.ErrorAppend(err2, "failed to setup tcp listener")
	}

	display.TunnelEstablished(name, port)

	// handle connections
	for {
		conn, err := serv.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept client connection: %s", err.Error())
		}

		go handleConnection(conn, key, location)
	}

	// return nil
}

func handleConnection(conn net.Conn, key, location string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("/tunnel?key=%s", key), nil)
	if err != nil {
		return
	}

	remoteConn, err := connect(req, location)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	go io.Copy(conn, remoteConn)
	_, err = io.Copy(remoteConn, conn)
	if err != nil {
		return
	}
}
