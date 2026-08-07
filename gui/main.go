package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"github.com/metruzanca/huijata/internal/config"
	"github.com/metruzanca/huijata/internal/snapshots"
)

type ui struct {
	win          fyne.Window
	cfg          *config.Config
	snapshotsDir string
	snaps        []snapshots.Snapshot
	pathLabel    *widget.Label
	status       *widget.Label
	list         *widget.List
	selected     widget.ListItemID
}

func main() {
	// The CLI's .env lives in the repo root and sets HUIJATA_CONFIG_PATH=".",
	// which resolves against the working directory. When the GUI is launched
	// from gui/, step up to the repo root so relative config paths resolve the
	// same way they do for the CLI.
	if _, err := os.Stat(".env"); err != nil {
		if _, err2 := os.Stat("../.env"); err2 == nil {
			_ = os.Chdir("..")
		}
	}
	_ = godotenv.Load()

	a := app.New()
	w := a.NewWindow("huijata")
	w.Resize(fyne.NewSize(680, 480))

	u := &ui{win: w}
	u.build()
	u.refresh()

	w.ShowAndRun()
}

func (u *ui) build() {
	u.pathLabel = widget.NewLabel("")
	u.pathLabel.Wrapping = fyne.TextWrapWord

	u.status = widget.NewLabel("")
	u.status.Wrapping = fyne.TextWrapWord

	u.list = widget.NewList(
		func() int { return len(u.snaps) },
		func() fyne.CanvasObject {
			desc := widget.NewLabel("Snapshot description")
			desc.Truncation = fyne.TextTruncateEllipsis
			when := widget.NewLabel("2026-01-01 00:00:00")
			when.Truncation = fyne.TextTruncateEllipsis
			when.TextStyle = fyne.TextStyle{Monospace: true}
			bg := canvas.NewRectangle(theme.Color(theme.ColorNameButton))
			bg.CornerRadius = theme.Size(theme.SizeNameSelectionRadius)
			bg.SetMinSize(fyne.NewSize(1, desc.MinSize().Height+when.MinSize().Height+theme.Padding()*3))
			return container.NewStack(bg, container.NewPadded(container.NewVBox(desc, when)))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			s := u.snaps[id]
			box := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(shortDesc(s.Description, 60))
			box.Objects[1].(*widget.Label).SetText(s.CreatedAt.Format("2006-01-02 15:04:05"))
		},
	)
	u.selected = -1
	u.list.OnSelected = func(id widget.ListItemID) {
		u.selected = id
		u.setStatus("")
	}

	initBtn := widget.NewButton("Init config", u.onInit)
	saveBtn := widget.NewButton("Save snapshot", u.onSave)
	restoreBtn := widget.NewButton("Restore selected", u.onRestore)
	clearBtn := widget.NewButton("Clear all", u.onClear)
	refreshBtn := widget.NewButton("Refresh", u.refresh)

	u.win.SetContent(container.NewBorder(
		container.NewVBox(u.pathLabel, u.status),
		container.NewHBox(initBtn, saveBtn, restoreBtn, clearBtn, refreshBtn),
		nil,
		nil,
		u.list,
	))
}

func (u *ui) refresh() {
	cfg, err := config.Load()
	if err != nil {
		u.setError("Failed to load config: " + err.Error())
		u.list.Refresh()
		return
	}
	if cfg == nil {
		u.cfg = nil
		u.pathLabel.SetText("")
		u.setError("No config found. Use the 'Init config' button below.")
		u.list.Refresh()
		return
	}
	dir, err := config.SnapshotsDir()
	if err != nil {
		u.setError("Failed to resolve snapshots dir: " + err.Error())
		u.list.Refresh()
		return
	}
	u.cfg = cfg
	u.snapshotsDir = dir
	u.pathLabel.SetText(fmt.Sprintf("Save path: %s\nSnapshots: %s", cfg.SavePath(), dir))

	u.snaps, err = snapshots.List(dir)
	if err != nil {
		u.setError("Failed to list snapshots: " + err.Error())
		u.list.Refresh()
		return
	}
	u.list.Refresh()
	u.setStatus(fmt.Sprintf("%d snapshot(s)", len(u.snaps)))
}

func (u *ui) onSave() {
	if u.cfg == nil {
		return
	}
	entry := widget.NewEntry()
	entry.SetPlaceHolder("e.g. strong run before Hiisi base")
	dialog.ShowCustomConfirm("Save snapshot", "Save", "Cancel",
		container.NewVBox(widget.NewLabel("Description:"), entry),
		func(ok bool) {
			if !ok {
				return
			}
			desc := strings.TrimSpace(entry.Text)
			if desc == "" {
				u.setError("Description must not be empty.")
				return
			}
			snap, err := snapshots.Create(u.cfg.SavePath(), u.snapshotsDir, desc, time.Now())
			if err != nil {
				u.setError("Save failed: " + err.Error())
				return
			}
			u.refresh()
			u.setStatus(fmt.Sprintf("Saved snapshot %s: %s", snap.ID, snap.Description))
		},
		u.win)
}

func (u *ui) onRestore() {
	if u.cfg == nil {
		return
	}
	id := u.selected
	if id < 0 || id >= len(u.snaps) {
		u.setError("Select a snapshot to restore first.")
		return
	}
	snap := u.snaps[id]
	dialog.ShowConfirm("Restore snapshot?",
		fmt.Sprintf("Roll the run back to %q. Unlocks and stats are kept, but any progress since this snapshot will be lost.", snap.Description),
		func(ok bool) {
			if !ok {
				return
			}
			if err := snapshots.Restore(snap.Dir, u.cfg.SavePath()); err != nil {
				u.setError("Restore failed: " + err.Error())
				return
			}
			u.setStatus(fmt.Sprintf("Restored snapshot %s: %s", snap.ID, snap.Description))
		},
		u.win)
}

func (u *ui) onClear() {
	if u.cfg == nil {
		return
	}
	count, err := snapshots.Count(u.snapshotsDir)
	if err != nil {
		u.setError("Failed to count snapshots: " + err.Error())
		return
	}
	if count == 0 {
		u.setStatus("No snapshots to clear.")
		return
	}
	dialog.ShowConfirm("Delete all snapshots?",
		fmt.Sprintf("This will permanently delete %d snapshot(s). Your save file won't be touched.", count),
		func(ok bool) {
			if !ok {
				return
			}
			deleted, err := snapshots.Clear(u.snapshotsDir)
			if err != nil {
				u.setError("Clear failed: " + err.Error())
				return
			}
			u.refresh()
			u.setStatus(fmt.Sprintf("Deleted %d snapshot(s).", deleted))
		},
		u.win)
}

func (u *ui) onInit() {
	guess := config.GuessGamePath()
	if guess != "" && config.DirExists(guess) {
		dialog.ShowCustomConfirm("Found your Noita game folder",
			"Yes, save it",
			"Choose another",
			container.NewVBox(widget.NewLabel("Is this it?"), widget.NewLabel(guess)),
			func(ok bool) {
				if ok {
					u.saveConfig(guess)
					return
				}
				u.pickGameFolder()
			}, u.win)
		return
	}
	u.pickGameFolder()
}

func (u *ui) pickGameFolder() {
	dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			u.setError("Failed to browse for game folder: " + err.Error())
			return
		}
		if uri == nil {
			return
		}
		u.saveConfig(uri.Path())
	}, u.win).Show()
}

func (u *ui) saveConfig(path string) {
	if path == "" {
		return
	}
	if !config.DirExists(path) {
		u.setError(fmt.Sprintf("%q is not a directory", path))
		return
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		u.setError("Failed to resolve path: " + err.Error())
		return
	}
	cfg := &config.Config{GamePath: abs}
	if err := cfg.Save(); err != nil {
		u.setError("Failed to save config: " + err.Error())
		return
	}
	u.refresh()
	u.setStatus("Game folder saved: " + abs)
}

func (u *ui) setStatus(msg string) {
	u.status.SetText(msg)
}

func shortDesc(s string, max int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= max {
		return string(r)
	}
	return string(r[:max-1]) + "…"
}

func (u *ui) setError(msg string) {
	u.setStatus("Error: " + msg)
}
