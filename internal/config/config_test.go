package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissing(t *testing.T) {
	t.Setenv("HUIJATA_CONFIG_PATH", "")
	t.Setenv("APPDATA", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg != nil {
		t.Fatalf("Load() = %v, want nil for missing config", cfg)
	}
}

func TestSavePathAppendsSave00(t *testing.T) {
	cfg := &Config{GamePath: `C:\Games\Noita`}
	want := filepath.Join(`C:\Games\Noita`, "save00")
	if got := cfg.SavePath(); got != want {
		t.Fatalf("SavePath() = %q, want %q", got, want)
	}
}

func TestPathEnvOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HUIJATA_CONFIG_PATH", dir)

	want := filepath.Join(dir, "config.toml")
	got, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
}

func TestSnapshotsDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HUIJATA_CONFIG_PATH", dir)

	want := filepath.Join(dir, "snapshots")
	got, err := SnapshotsDir()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("SnapshotsDir() = %q, want %q", got, want)
	}
}

func TestSaveAndLoad(t *testing.T) {
	t.Setenv("HUIJATA_CONFIG_PATH", "")
	t.Setenv("APPDATA", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg := &Config{GamePath: `C:\Games\Noita`}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded == nil {
		t.Fatal("expected config to load")
	}
	if loaded.GamePath != cfg.GamePath {
		t.Fatalf("GamePath = %q, want %q", loaded.GamePath, cfg.GamePath)
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "huijata", "config.toml")
	got, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
}
