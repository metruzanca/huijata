package noitadata

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// writeWak builds a data.wak file in dir from the given entries and returns
// its path.
func writeWak(t *testing.T, dir string, entries []struct {
	name string
	data []byte
}) string {
	t.Helper()

	dirSize := 0
	for _, e := range entries {
		dirSize += 4 + len(e.name) + 8
	}
	dataOffset := uint32(24 + dirSize)

	offs := make([]uint32, len(entries))
	var data []byte
	for i, e := range entries {
		data = append(data, e.data...)
		offs[i] = dataOffset + uint32(len(data))
	}

	var dirBuf bytes.Buffer
	for i, e := range entries {
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], uint32(len(e.name)))
		dirBuf.Write(b[:])
		dirBuf.WriteString(e.name)
		binary.LittleEndian.PutUint32(b[:], offs[i])
		dirBuf.Write(b[:])
		dirBuf.Write([]byte{0, 0, 0, 0}) // a field, meaning unknown
	}

	var hdr [24]byte
	binary.LittleEndian.PutUint32(hdr[4:8], uint32(len(entries)))
	binary.LittleEndian.PutUint32(hdr[8:12], dataOffset)
	binary.LittleEndian.PutUint32(hdr[16:20], dataOffset)
	if len(entries) > 0 {
		binary.LittleEndian.PutUint32(hdr[20:24], uint32(len(entries[0].data)-1))
	}

	var out bytes.Buffer
	out.Write(hdr[:])
	out.Write(dirBuf.Bytes())
	out.Write(data)

	path := filepath.Join(dir, "data.wak")
	if err := os.WriteFile(path, out.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestWakFind(t *testing.T) {
	lua := []byte(`dofile_once("x")
local actions = {
	{ id = "RUBBER_BALL", name = "$action_rubber_ball" },
}`)
	dir := t.TempDir()
	path := writeWak(t, dir, []struct {
		name string
		data []byte
	}{
		{name: "data/credits.txt", data: []byte("credits")},
		{name: "data/scripts/gun/gun_actions.lua", data: lua},
	})

	w, err := OpenWak(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := w.Find("data/scripts/gun/gun_actions.lua")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, lua) {
		t.Fatalf("extracted %q, want %q", got, lua)
	}

	if _, err := w.Find("data/nope.lua"); err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestWakFindInflate(t *testing.T) {
	text := []byte("this is a zlib compressed file")
	var comp bytes.Buffer
	zw := zlib.NewWriter(&comp)
	if _, err := zw.Write(text); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	path := writeWak(t, dir, []struct {
		name string
		data []byte
	}{
		{name: "data/compressed.bin", data: comp.Bytes()},
	})
	w, err := OpenWak(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := w.Find("data/compressed.bin")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, text) {
		t.Fatalf("decompressed %q, want %q", got, text)
	}
}

func TestSpellNames(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data", "translations"), 0o755); err != nil {
		t.Fatal(err)
	}
	lua := []byte(`local actions = {
	{ id = "RUBBER_BALL", name = "$action_rubber_ball" },
	{ id = "LIGHTNING_BOLT", name = "$action_lightning" },
	{ id = "BOMB", name = "$action_bomb" },
	{ id = "FUNKY_SPELL", name = "???" },
	{ id = "NO_TRANSLATION", name = "$action_missing_key" },
}`)
	writeWak(t, filepath.Join(dir, "data"), []struct {
		name string
		data []byte
	}{
		{name: "data/scripts/gun/gun_actions.lua", data: lua},
	})
	csv := "key,en\n" +
		"action_rubber_ball,Bouncing burst\n" +
		"action_lightning,Lightning bolt\n" +
		"action_bomb,Bomb\n"
	if err := os.WriteFile(filepath.Join(dir, "data", "translations", "common.csv"), []byte(csv), 0o644); err != nil {
		t.Fatal(err)
	}

	names, err := SpellNames(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"RUBBER_BALL":   "Bouncing burst",
		"LIGHTNING_BOLT": "Lightning bolt",
		"BOMB":           "Bomb",
	}
	for id, name := range want {
		if got := names[id]; got != name {
			t.Errorf("%s = %q, want %q", id, got, name)
		}
	}
	for _, id := range []string{"FUNKY_SPELL", "NO_TRANSLATION"} {
		if _, ok := names[id]; ok {
			t.Errorf("expected %s to be omitted", id)
		}
	}
}

func TestSpellNamesBadInstall(t *testing.T) {
	if _, err := SpellNames(t.TempDir()); err == nil {
		t.Fatal("expected error for a dir without data.wak")
	}
}

func TestDetectInstallEnv(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "data", "data.wak"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HUIJATA_NOITA_INSTALL", dir)
	got, err := DetectInstall()
	if err != nil {
		t.Fatal(err)
	}
	if got != dir {
		t.Fatalf("DetectInstall() = %q, want %q", got, dir)
	}
}

func TestParseLibraryFolders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "libraryfolders.vdf")
	content := `"libraryfolders"
{
	"0" { "path" "D:\SteamLibrary" }
	"1" { "path" "D:\SteamLibrary" }
	"2" { "path" "E:\Games\Steam" }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	libs := parseLibraryFolders(path)
	want := []string{`D:\SteamLibrary`, `E:\Games\Steam`}
	if len(libs) != len(want) {
		t.Fatalf("got %v, want %v", libs, want)
	}
	for i := range want {
		if libs[i] != want[i] {
			t.Errorf("lib[%d] = %q, want %q", i, libs[i], want[i])
		}
	}
}
