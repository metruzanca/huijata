package stats

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// Session is one run recorded in stats/sessions/. Time comes from the session
// file name; Stats holds the attributes of the run's <stats> block.
type Session struct {
	Time  time.Time
	Stats Section
}

var sessionFileRE = regexp.MustCompile(`^(\d{8}-\d{6})_stats\.xml$`)

// Sessions reads the per-run stats from stats/sessions/ in the save slot
// folder, newest first. A missing sessions folder yields an empty list.
func Sessions(savePath string) ([]Session, error) {
	dir := filepath.Join(savePath, "stats", "sessions")
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if s, ok := parseSessionFile(filepath.Join(dir, e.Name())); ok {
			sessions = append(sessions, s)
		}
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Time.After(sessions[j].Time)
	})
	return sessions, nil
}

// parseSessionFile reads one <timestamp>_stats.xml file. Files that do not
// match the session naming scheme, or that cannot be parsed, are skipped.
func parseSessionFile(path string) (Session, bool) {
	m := sessionFileRE.FindStringSubmatch(filepath.Base(path))
	if m == nil {
		return Session{}, false
	}
	t, err := time.Parse("20060102-150405", m[1])
	if err != nil {
		return Session{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Session{}, false
	}
	var root element
	if err := xml.Unmarshal(data, &root); err != nil {
		return Session{}, false
	}
	s := Session{Time: t, Stats: Section{}}
	for i := range root.Kids {
		if root.Kids[i].XMLName.Local == "stats" {
			s.Stats = section(&root.Kids[i])
			break
		}
	}
	return s, true
}
