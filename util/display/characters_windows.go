//go:build windows
// +build windows

package display

var (
	TaskSpinner  = []string{`\`, `|`, `/`, `-`}
	TaskComplete = `√`
	TaskPause    = `*`
)
