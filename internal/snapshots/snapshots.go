package snapshots

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RunFiles are the files and folders that make up a Noita run's state. They
// appear when a new run starts and get wiped when the run ends, so they are the
// only things a snapshot needs to capture. Everything else (persistent/,
// stats/, mod configs, steam_autocloud.vdf) is lifetime progression. Snapshots
// leave that alone.
var RunFiles = []string{
	"player.xml",
	"session_numbers.salakieli",
	"world_state.xml",
	"world",
	"stats/_streaks.salakieli",
}

// Snapshot is a saved copy of a run's state.
type Snapshot struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Dir         string    `json:"-"`
}

// Create copies the run-state files from savePath into a new snapshot folder
// inside snapshotsDir. The folder name is a generated unique id.
func Create(savePath, snapshotsDir, description string, now time.Time) (*Snapshot, error) {
	id, err := newID(now)
	if err != nil {
		return nil, err
	}

	snap := &Snapshot{
		ID:          id,
		Description: description,
		CreatedAt:   now,
		Dir:         filepath.Join(snapshotsDir, id),
	}

	if err := os.MkdirAll(snap.Dir, 0o755); err != nil {
		return nil, err
	}

	for _, rel := range RunFiles {
		src := filepath.Join(savePath, rel)
		if !pathExists(src) {
			continue
		}
		if err := copyAny(src, filepath.Join(snap.Dir, rel)); err != nil {
			return nil, fmt.Errorf("copying %s: %w", rel, err)
		}
	}

	if err := writeMeta(snap); err != nil {
		return nil, err
	}
	return snap, nil
}

// List returns the snapshots in snapshotsDir, newest first. Folders without a
// readable meta.json are ignored.
func List(snapshotsDir string) ([]Snapshot, error) {
	entries, err := os.ReadDir(snapshotsDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var snaps []Snapshot
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(snapshotsDir, e.Name())
		snap, err := readMeta(dir)
		if err != nil {
			continue
		}
		snap.Dir = dir
		snaps = append(snaps, *snap)
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].CreatedAt.After(snaps[j].CreatedAt)
	})
	return snaps, nil
}

// Restore makes savePath match the snapshot's run-state files. Run files in
// the snapshot are copied over. Run files missing from the snapshot are
// removed from savePath, so a dead-state snapshot correctly wipes the run.
// Everything outside RunFiles is left untouched.
func Restore(snapDir, savePath string) error {
	for _, rel := range RunFiles {
		src := filepath.Join(snapDir, rel)
		dst := filepath.Join(savePath, rel)
		if pathExists(src) {
			if err := os.RemoveAll(dst); err != nil {
				return err
			}
			if err := copyAny(src, dst); err != nil {
				return fmt.Errorf("restoring %s: %w", rel, err)
			}
		} else if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}
	return nil
}

// Clear removes every snapshot and leaves the player's save file untouched.
// It returns the number of snapshot folders deleted.
func Clear(snapshotsDir string) (int, error) {
	count, err := Count(snapshotsDir)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, nil
	}
	if err := os.RemoveAll(snapshotsDir); err != nil {
		return 0, err
	}
	return count, nil
}

// Count returns how many snapshot folders exist in snapshotsDir.
func Count(snapshotsDir string) (int, error) {
	entries, err := os.ReadDir(snapshotsDir)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			count++
		}
	}
	return count, nil
}

func newID(now time.Time) (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return now.Format("20060102-150405") + "-" + hex.EncodeToString(b[:]), nil
}

func writeMeta(snap *Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(snap.Dir, "meta.json"), data, 0o644)
}

func readMeta(dir string) (*Snapshot, error) {
	data, err := os.ReadFile(filepath.Join(dir, "meta.json"))
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyAny(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst, info.Mode())
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
