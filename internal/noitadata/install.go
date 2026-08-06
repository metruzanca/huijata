package noitadata

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

// DetectInstall finds the Noita install folder (the one containing
// data/data.wak). It checks, in order: the HUIJATA_NOITA_INSTALL env var,
// the default Steam install locations, then any Steam library folders listed
// in libraryfolders.vdf.
func DetectInstall() (string, error) {
	candidates := []string{}
	if p := os.Getenv("HUIJATA_NOITA_INSTALL"); p != "" {
		candidates = append(candidates, p)
	}
	candidates = append(candidates,
		filepath.Join("C:\\", "Program Files (x86)", "Steam", "steamapps", "common", "Noita"),
		filepath.Join("C:\\", "Program Files", "Steam", "steamapps", "common", "Noita"),
	)
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".local", "share", "Steam", "steamapps", "common", "Noita"),
			filepath.Join(home, ".steam", "steam", "steamapps", "common", "Noita"),
			filepath.Join(home, "Library", "Application Support", "Steam", "steamapps", "common", "Noita"),
		)
	}
	for _, lib := range steamLibraryFolders() {
		candidates = append(candidates, filepath.Join(lib, "steamapps", "common", "Noita"))
	}

	for _, p := range candidates {
		if fileExists(filepath.Join(p, "data", "data.wak")) {
			return p, nil
		}
	}
	return "", errors.New("could not find the Noita install folder (set HUIJATA_NOITA_INSTALL to point at it)")
}

// steamLibraryFolders returns the library paths parsed from the Steam
// libraryfolders.vdf files of the default Steam install locations.
func steamLibraryFolders() []string {
	var libs []string
	for _, base := range []string{
		filepath.Join("C:\\", "Program Files (x86)", "Steam"),
		filepath.Join("C:\\", "Program Files", "Steam"),
	} {
		libs = append(libs, parseLibraryFolders(filepath.Join(base, "steamapps", "libraryfolders.vdf"))...)
	}
	return libs
}

var vdfPath = regexp.MustCompile(`"path"\s+"([^"]+)"`)

func parseLibraryFolders(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	var libs []string
	for _, m := range vdfPath.FindAllSubmatch(data, -1) {
		lib := string(m[1])
		if !seen[lib] {
			seen[lib] = true
			libs = append(libs, lib)
		}
	}
	return libs
}
