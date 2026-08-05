package cmd

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultSavePathWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows only")
	}
	t.Setenv("USERPROFILE", `C:\Users\test`)
	want := filepath.Join(`C:\Users\test\AppData`, "LocalLow", "Nolla_Games_Noita")
	if got := defaultSavePath(); got != want {
		t.Fatalf("defaultSavePath() = %q, want %q", got, want)
	}
}
