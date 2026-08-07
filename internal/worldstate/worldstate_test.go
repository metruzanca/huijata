package worldstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	data := `<Entity name="unknown" tags="world_state" >
  <WorldStateComponent
    EVERYTHING_TO_GOLD="1"
    INFINITE_GOLD_HAPPENING="0"
    day_count="3"
    is_initialized="1" >
  </WorldStateComponent>
  <Entity name="player_stats" ></Entity>
</Entity>`
	if err := os.WriteFile(filepath.Join(dir, "world_state.xml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if s.Attrs["EVERYTHING_TO_GOLD"] != "1" {
		t.Errorf("EVERYTHING_TO_GOLD = %q, want 1", s.Attrs["EVERYTHING_TO_GOLD"])
	}
	if s.Attrs["day_count"] != "3" {
		t.Errorf("day_count = %q, want 3", s.Attrs["day_count"])
	}
	if s.Attrs["is_initialized"] != "1" {
		t.Errorf("is_initialized = %q, want 1", s.Attrs["is_initialized"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(t.TempDir()); err == nil {
		t.Fatal("expected an error for a missing world_state.xml")
	}
}

func TestLoadNoComponent(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "world_state.xml"), []byte("<Entity></Entity>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir); err == nil {
		t.Fatal("expected an error when there is no WorldStateComponent")
	}
}
