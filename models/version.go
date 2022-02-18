package models

import "fmt"

var (
	// will be set with build flags, defaults for one-off `go build`
	microVersion string = "0.0.0"  // git tag
	microCommit  string = "custom" // commit id of build
	microBuild   string = "now"    // date of build
)

func VersionString() string {
	return fmt.Sprintf("Microbox Version %s-%s (%s)", microVersion, microBuild, microCommit)
}
