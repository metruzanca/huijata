package stats

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSessions(t *testing.T) {
	dir := t.TempDir()
	sessDir := filepath.Join(dir, "stats", "sessions")
	if err := os.MkdirAll(sessDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"20260807-014746_stats.xml": `<Stats BUILD_NAME="x" >
  <stats dead="1" gold="10" enemies_killed="5" playtime_str="0:00:02" world_seed="111" ></stats>
</Stats>`,
		"20260806-100000_stats.xml": `<Stats BUILD_NAME="x" >
  <stats dead="0" gold="99" enemies_killed="40" playtime_str="1:00:00" world_seed="222" ></stats>
</Stats>`,
		"not_a_session.txt": `junk`,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(sessDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	sessions, err := Sessions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("Sessions = %d entries, want 2", len(sessions))
	}
	// newest first
	if sessions[0].Time.Format("20060102-150405") != "20260807-014746" {
		t.Errorf("Sessions[0] = %v, want 20260807-014746", sessions[0].Time)
	}
	if sessions[1].Time.Format("20060102-150405") != "20260806-100000" {
		t.Errorf("Sessions[1] = %v, want 20260806-100000", sessions[1].Time)
	}
	if sessions[0].Stats["gold"] != "10" || sessions[0].Stats["world_seed"] != "111" {
		t.Errorf("Sessions[0].Stats = %v", sessions[0].Stats)
	}
	if sessions[1].Stats["dead"] != "0" || sessions[1].Stats["enemies_killed"] != "40" {
		t.Errorf("Sessions[1].Stats = %v", sessions[1].Stats)
	}
}

func TestSessionsMissingFolder(t *testing.T) {
	sessions, err := Sessions(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 0 {
		t.Errorf("Sessions = %d entries, want 0", len(sessions))
	}
}

func TestSessionTimestamp(t *testing.T) {
	want, err := time.Parse("20060102-150405", "20260807-014746")
	if err != nil {
		t.Fatal(err)
	}
	if want.Day() != 7 || want.Month() != time.August || want.Year() != 2026 {
		t.Errorf("parsed time = %v", want)
	}
}
