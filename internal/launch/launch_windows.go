//go:build windows

package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const appID = "881100"

// LaunchSteam asks the Steam client to run Noita. Steam is started if it is
// not running yet.
func LaunchSteam() error {
	return exec.Command("cmd", "/c", "start", "", "steam://rungameid/"+appID).Start()
}

// LaunchDirect runs the Noita executable in installDir.
func LaunchDirect(installDir string) error {
	exe := filepath.Join(installDir, "noita.exe")
	if _, err := os.Stat(exe); err != nil {
		return fmt.Errorf("Noita executable not found at %s", exe)
	}
	return exec.Command(exe).Start()
}
