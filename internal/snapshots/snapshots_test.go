package snapshots

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newTestSave(t *testing.T) string {
	t.Helper()
	save := t.TempDir()
	writeTestFile(t, filepath.Join(save, "player.xml"), "<player/>")
	writeTestFile(t, filepath.Join(save, "session_numbers.salakieli"), "1")
	writeTestFile(t, filepath.Join(save, "world_state.xml"), "<world/>")
	writeTestFile(t, filepath.Join(save, "world", "00", "cell00.xml"), "cell")
	writeTestFile(t, filepath.Join(save, "stats", "_streaks.salakieli"), "5")
	writeTestFile(t, filepath.Join(save, "stats", "worldmap.bin"), "keep me")
	writeTestFile(t, filepath.Join(save, "persistent", "flags", "action_wand"), "keep me")
	return save
}

func TestCreateSnapshotsRunFilesOnly(t *testing.T) {
	save := newTestSave(t)
	snapshotsDir := t.TempDir()
	now := time.Date(2026, 8, 5, 10, 0, 0, 0, time.UTC)

	snap, err := Create(save, snapshotsDir, "first run", now)
	if err != nil {
		t.Fatal(err)
	}

	if snap.ID == "" {
		t.Fatal("expected non-empty id")
	}
	if snap.Description != "first run" {
		t.Fatalf("Description = %q, want %q", snap.Description, "first run")
	}
	if !snap.CreatedAt.Equal(now) {
		t.Fatalf("CreatedAt = %v, want %v", snap.CreatedAt, now)
	}

	for _, rel := range RunFiles {
		dst := filepath.Join(snap.Dir, rel)
		if _, err := os.Stat(dst); err != nil {
			t.Errorf("missing run file %s: %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(snap.Dir, "stats", "worldmap.bin")); !os.IsNotExist(err) {
		t.Error("snapshot must not include stats/worldmap.bin")
	}
	if _, err := os.Stat(filepath.Join(snap.Dir, "persistent")); !os.IsNotExist(err) {
		t.Error("snapshot must not include persistent/")
	}
	if _, err := os.Stat(filepath.Join(snap.Dir, "meta.json")); err != nil {
		t.Errorf("missing meta.json: %v", err)
	}
}

func TestCreateSkipsMissingRunFiles(t *testing.T) {
	save := t.TempDir()
	writeTestFile(t, filepath.Join(save, "player.xml"), "<player/>")

	snap, err := Create(save, t.TempDir(), "partial", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(snap.Dir, "world")); !os.IsNotExist(err) {
		t.Error("world/ must not be created when absent from save")
	}
}

func TestDuplicateDescriptionsGetDistinctIDs(t *testing.T) {
	save := newTestSave(t)
	snapshotsDir := t.TempDir()

	a, err := Create(save, snapshotsDir, "same", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	b, err := Create(save, snapshotsDir, "same", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID {
		t.Fatalf("duplicate descriptions produced the same id %q", a.ID)
	}
}

func TestListNewestFirst(t *testing.T) {
	save := newTestSave(t)
	snapshotsDir := t.TempDir()

	older := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)
	if _, err := Create(save, snapshotsDir, "older", older); err != nil {
		t.Fatal(err)
	}
	if _, err := Create(save, snapshotsDir, "newer", newer); err != nil {
		t.Fatal(err)
	}

	snaps, err := List(snapshotsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 2 {
		t.Fatalf("len = %d, want 2", len(snaps))
	}
	if snaps[0].Description != "newer" || snaps[1].Description != "older" {
		t.Fatalf("order = [%s, %s], want [newer, older]", snaps[0].Description, snaps[1].Description)
	}
	if snaps[0].Dir == "" {
		t.Error("Dir must be populated by List")
	}
}

func TestListEmpty(t *testing.T) {
	snaps, err := List(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 0 {
		t.Fatalf("len = %d, want 0", len(snaps))
	}
}

func TestListIgnoresFoldersWithoutMeta(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "junk"), 0o755); err != nil {
		t.Fatal(err)
	}
	snaps, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 0 {
		t.Fatalf("len = %d, want 0", len(snaps))
	}
}

func TestRestoreMirrorsRunFiles(t *testing.T) {
	save := newTestSave(t)
	snapshotsDir := t.TempDir()
	snap, err := Create(save, snapshotsDir, "run", time.Now())
	if err != nil {
		t.Fatal(err)
	}

	// Wipe the run but keep lifetime files, then restore.
	for _, rel := range RunFiles {
		if err := os.RemoveAll(filepath.Join(save, rel)); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(save, "stats", "_streaks.salakieli"), []byte("99"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Restore(snap.Dir, save); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(save, "player.xml"))
	if err != nil || string(got) != "<player/>" {
		t.Errorf("player.xml not restored: %v, %q", err, got)
	}
	got, err = os.ReadFile(filepath.Join(save, "world", "00", "cell00.xml"))
	if err != nil || string(got) != "cell" {
		t.Errorf("world file not restored: %v, %q", err, got)
	}
	got, err = os.ReadFile(filepath.Join(save, "stats", "_streaks.salakieli"))
	if err != nil || string(got) != "5" {
		t.Errorf("_streaks.salakieli not restored: %v, %q", err, got)
	}
	got, err = os.ReadFile(filepath.Join(save, "persistent", "flags", "action_wand"))
	if err != nil || string(got) != "keep me" {
		t.Errorf("persistent file was touched: %v, %q", err, got)
	}
	got, err = os.ReadFile(filepath.Join(save, "stats", "worldmap.bin"))
	if err != nil || string(got) != "keep me" {
		t.Errorf("stats/worldmap.bin was touched: %v, %q", err, got)
	}
}

func TestRestoreDeadStateWipesRun(t *testing.T) {
	save := newTestSave(t)
	snapshotsDir := t.TempDir()

	// Dead state: only meta.json is stored, no run files.
	if err := os.MkdirAll(snapshotsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	id := "dead-snap"
	snap := &Snapshot{ID: id, Description: "dead", CreatedAt: time.Now(), Dir: filepath.Join(snapshotsDir, id)}
	if err := os.MkdirAll(snap.Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeMeta(snap); err != nil {
		t.Fatal(err)
	}

	if err := Restore(snap.Dir, save); err != nil {
		t.Fatal(err)
	}

	for _, rel := range RunFiles {
		if _, err := os.Stat(filepath.Join(save, rel)); !os.IsNotExist(err) {
			t.Errorf("run file %s must be removed, err = %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(save, "persistent", "flags", "action_wand")); err != nil {
		t.Errorf("persistent must survive restore: %v", err)
	}
}

func TestClear(t *testing.T) {
	save := newTestSave(t)
	snapshotsDir := t.TempDir()

	for _, desc := range []string{"a", "b", "c"} {
		if _, err := Create(save, snapshotsDir, desc, time.Now()); err != nil {
			t.Fatal(err)
		}
	}

	count, err := Clear(snapshotsDir)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("count = %d, want 3", count)
	}
	if _, err := os.Stat(snapshotsDir); !os.IsNotExist(err) {
		t.Errorf("snapshots dir must be removed, err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(save, "player.xml")); err != nil {
		t.Errorf("save file must survive clear: %v", err)
	}
}

func TestClearEmpty(t *testing.T) {
	count, err := Clear(filepath.Join(t.TempDir(), "missing"))
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("count = %d, want 0", count)
	}
}


