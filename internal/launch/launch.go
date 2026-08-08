package launch

import (
	"path/filepath"
	"strings"
)

// Launch starts Noita. When the game is a Steam install it goes through the
// Steam client so Steam tracks the run; otherwise it launches the executable
// directly. Pass direct to force a direct launch regardless of install type.
func Launch(installDir string, direct bool) error {
	if !direct && IsSteamInstall(installDir) {
		return LaunchSteam()
	}
	return LaunchDirect(installDir)
}

// IsSteamInstall reports whether installDir lives inside a Steam library,
// where the game sits under a steamapps/common folder.
func IsSteamInstall(installDir string) bool {
	for _, part := range strings.Split(filepath.ToSlash(installDir), "/") {
		if part == "steamapps" {
			return true
		}
	}
	return false
}
