//go:build windows

package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
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

// isRunning reports whether a Noita process is currently running, by walking
// the process snapshot for an executable named noita.exe.
func isRunning() bool {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	for err := windows.Process32First(snapshot, &entry); err == nil; err = windows.Process32Next(snapshot, &entry) {
		if strings.EqualFold(windows.UTF16ToString(entry.ExeFile[:]), "noita.exe") {
			return true
		}
	}
	return false
}
