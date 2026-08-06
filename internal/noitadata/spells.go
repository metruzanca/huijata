package noitadata

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// luaAction matches one entry of gun_actions.lua and captures its id and its
// translation key (e.g. id = "RUBBER_BALL", name = "$action_rubber_ball").
var luaAction = regexp.MustCompile(`id\s*=\s*"([^"]+)"\s*,\s*name\s*=\s*"([^"]+)"`)

// SpellNames returns a map from a spell's internal id (e.g. "RUBBER_BALL")
// to its English display name (e.g. "Bouncing burst"), by reading the game's
// own data: gun_actions.lua inside data.wak maps id -> translation key, and
// data/translations/common.csv maps the key -> English name. Spells whose
// name key is missing from common.csv (or is a placeholder) are omitted.
func SpellNames(installDir string) (map[string]string, error) {
	install, err := resolveInstall(installDir)
	if err != nil {
		return nil, err
	}

	wak, err := OpenWak(filepath.Join(install, "data", "data.wak"))
	if err != nil {
		return nil, err
	}
	lua, err := wak.Find("data/scripts/gun/gun_actions.lua")
	if err != nil {
		return nil, err
	}

	idToKey := make(map[string]string)
	for _, m := range luaAction.FindAllSubmatch(lua, -1) {
		id := string(m[1])
		key := strings.TrimPrefix(string(m[2]), "$")
		if key != "" {
			idToKey[id] = key
		}
	}

	keyToName, err := parseCommonCSV(filepath.Join(install, "data", "translations", "common.csv"))
	if err != nil {
		return nil, err
	}

	names := make(map[string]string, len(idToKey))
	for id, key := range idToKey {
		if name, ok := keyToName[key]; ok {
			names[id] = name
		}
	}
	return names, nil
}

// parseCommonCSV reads the English names for action_* keys from the game's
// translation file.
func parseCommonCSV(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	names := make(map[string]string)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 2 {
			continue
		}
		key := rec[0]
		name := strings.TrimSpace(rec[1])
		if strings.HasPrefix(key, "action_") && name != "" {
			names[key] = name
		}
	}
	return names, nil
}

func resolveInstall(installDir string) (string, error) {
	if installDir != "" {
		if fileExists(filepath.Join(installDir, "data", "data.wak")) {
			return installDir, nil
		}
		return "", fmt.Errorf("no data/data.wak in %q", installDir)
	}
	return DetectInstall()
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
