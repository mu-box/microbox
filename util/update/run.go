package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mu-box/microbox/util/display"
)

var Server bool

func Run(path string) error {
	if path == "" {
		return fmt.Errorf("invalid path")
	}

	// create a temporary file
	tmpFileName := filepath.Join(filepath.Dir(path), TmpName)
	tmpFile, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if !strings.Contains(path, "microbox-update") {
		fmt.Printf("Current version: %s", getCurrentVersion(path))
	}

	// download the file and display the progress bar
	resp, err := http.Get(remotePath())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dp := display.DownloadPercent{Total: resp.ContentLength}
	if Server {
		// on the Server we dont really care to see this
		dp.Output = ioutil.Discard
	}
	dp.Copy(tmpFile, resp.Body)

	// close the tmp file
	tmpFile.Close()

	// replace binary
	if err := os.Rename(tmpFileName, path); err != nil {
		return err
	}

	// update the model
	update := newUpdate()

	if !strings.Contains(path, "microbox-update") {
		fmt.Printf("\nUpdated to version: %s\n\n", getCurrentVersion(path))
		fmt.Println("Check out the release notes here:")
		fmt.Println("https://github.com/mu-box/microbox/blob/master/CHANGELOG.md")
	}

	return update.Save()
}

func getCurrentVersion(path string) string {
	if path == "" {
		fmt.Println("invalid path")
		return ""
	}
	version, err := exec.Command(path, "version").Output()
	if err != nil {
		fmt.Println("Error while trying to get the microbox version")
		return ""
	}
	return string(version)
}
