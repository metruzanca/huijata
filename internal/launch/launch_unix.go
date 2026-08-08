//go:build !windows

package launch

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

const appID = "881100"

// LaunchSteam asks the Steam client to run Noita via a steam:// URI, using
// the platform's default opener so Steam starts if it is not running.
func LaunchSteam() error {
	uri := "steam://rungameid/" + appID
	if runtime.GOOS == "darwin" {
		return exec.Command("open", uri).Start()
	}
	return exec.Command("xdg-open", uri).Start()
}

// LaunchDirect runs the Noita executable in installDir.
func LaunchDirect(installDir string) error {
	return exec.Command(filepath.Join(installDir, "noita.exe")).Start()
}
