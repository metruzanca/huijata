package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

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
		func() fyne.CanvasObject { return widget.NewLabel("snapshot") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			s := u.snaps[id]
			label.SetText(fmt.Sprintf("%s\n%s", s.Description, s.CreatedAt.Format("2006-01-02 15:04:05")))
		},
	)
	u.selected = -1
	u.list.OnSelected = func(id widget.ListItemID) {
		u.selected = id
		u.setStatus("")
	}

	saveBtn := widget.NewButton("Save snapshot", u.onSave)
	restoreBtn := widget.NewButton("Restore selected", u.onRestore)
	clearBtn := widget.NewButton("Clear all", u.onClear)
	refreshBtn := widget.NewButton("Refresh", u.refresh)

	u.win.SetContent(container.NewBorder(
		container.NewVBox(u.pathLabel, u.status),
		container.NewHBox(saveBtn, restoreBtn, clearBtn, refreshBtn),
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
		u.setError("No config found. Run `huijata init` from the CLI first, then restart the app.")
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

func (u *ui) setStatus(msg string) {
	u.status.SetText(msg)
}

func (u *ui) setError(msg string) {
	u.setStatus("Error: " + msg)
}
