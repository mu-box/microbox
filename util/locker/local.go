package locker

import (
	"fmt"
	"net"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/config"
)

var (
	lln    net.Listener // local locking network
	lCount int
)

// LocalLock locks on port
func LocalLock() error {

	//
	for {
		if success, _ := LocalTryLock(); success {
			break
		}
		lumber.Trace("local lock waiting...")
		<-time.After(time.Second)
	}

	mutex.Lock()
	lCount++
	lumber.Trace("local lock acquired (%d)", lCount)
	mutex.Unlock()

	return nil
}

// LocalTryLock ...
func LocalTryLock() (bool, error) {

	var err error

	//
	if lln != nil {
		return true, nil
	}

	//
	config, _ := models.LoadConfig()
	port := config.LockPort
	if port == 0 {
		port = 12345
	}
	port = port + localPort()

	//
	if lln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
		return true, nil
	}

	return false, nil
}

// LocalUnlock ...
func LocalUnlock() (err error) {

	mutex.Lock()
	lCount--
	lumber.Trace("local lock released (%d)", lCount)
	mutex.Unlock()

	// if im not the last guy to release my lock quit immediately instead of closing
	// the connection
	if lCount > 0 || lln == nil {
		return nil
	}

	err = lln.Close()
	lln = nil

	return
}

// localPort ...
func localPort() (num int) {

	b := []byte(config.EnvID())

	//
	for i := 0; i < len(b); i++ {
		num = num + int(b[i])
	}

	return num
}
