package watch

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNotifyFiles(t *testing.T) {
	os.MkdirAll("/tmp/microbox/", 0777)
	notifyWatcher, err := newRecursiveWatcher("/tmp/microbox/")
	defer notifyWatcher.close()
	if err != nil {
		t.Fatalf("failed to watch: %s", err)
	}
	notifyWatcher.watch()

	<-time.After(time.Second)
	ioutil.WriteFile("/tmp/microbox/notify.tmp", []byte("hi"), 0777)

	// pull the first event off the channel
	ev := <-notifyWatcher.eventChan()

	if ev.file != "/tmp/microbox/notify.tmp" {
		t.Errorf("the wrong file path came out %s", ev.file)
	}
	if ev.error != nil {
		t.Errorf("an error occurred %s", ev.error)
	}
}
