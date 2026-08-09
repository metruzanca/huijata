package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/metruzanca/huijata/internal/config"
	"github.com/metruzanca/huijata/internal/snapshots"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore the game to a saved snapshot",
	Long: `Rolls the current run back to a previously saved snapshot. Any progress
made since that snapshot is lost; lifetime progression (unlocks, stats) is
kept.`,
	RunE: runRestore,
}

var restoreYes bool

func runRestore(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	snapshotsDir, err := config.SnapshotsDir()
	if err != nil {
		return err
	}

	snaps, err := snapshots.List(snapshotsDir)
	if err != nil {
		return err
	}
	if len(snaps) == 0 {
		return errors.New("no snapshots found, run `huijata save <description>` first")
	}

	var chosen snapshots.Snapshot
	if restoreYes {
		chosen = snaps[0]
	} else {
		proceed := false
		chosen, proceed, err = selectSnapshot(snaps)
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}
	}

	if err := snapshots.Restore(chosen.Dir, cfg.SavePath()); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Restored snapshot %s: %s\n", chosen.ID, chosen.Description)
	return maybeStartNoita(cmd, !restoreYes)
}

// selectSnapshot prompts the user to pick a snapshot and confirm the restore.
// The bool reports whether the restore should proceed; false means the user
// cancelled (esc or ctrl+C), in which case nothing is restored.
func selectSnapshot(snaps []snapshots.Snapshot) (snapshots.Snapshot, bool, error) {
	options := make([]huh.Option[string], 0, len(snaps))
	for _, s := range snaps {
		options = append(options, huh.NewOption(
			fmt.Sprintf("%s — %s", s.Description, s.CreatedAt.Format("2006-01-02 15:04:05")),
			s.ID,
		))
	}

	var chosenID string
	selectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Restore which snapshot?").
				Description("Your current run's progress will be replaced.").
				Options(options...).
				Value(&chosenID),
		),
	)
	if err := selectForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return snapshots.Snapshot{}, false, nil
		}
		return snapshots.Snapshot{}, false, err
	}

	var chosen snapshots.Snapshot
	for _, s := range snaps {
		if s.ID == chosenID {
			chosen = s
			break
		}
	}

	confirm := true
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Restore this snapshot?").
				Description(fmt.Sprintf("Roll the run back to %q. Unlocks and stats are kept, but any progress since this snapshot will be lost.", chosen.Description)).
				Affirmative("Yes, restore").
				Negative("Cancel").
				Value(&confirm),
		),
	)
	if err := confirmForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return snapshots.Snapshot{}, false, nil
		}
		return snapshots.Snapshot{}, false, err
	}
	if !confirm {
		return snapshots.Snapshot{}, false, nil
	}
	return chosen, true, nil
}

func init() {
	restoreCmd.Flags().BoolVarP(&restoreYes, "yes", "y", false, "restore the latest snapshot and start Noita without asking")
	rootCmd.AddCommand(restoreCmd)
}
