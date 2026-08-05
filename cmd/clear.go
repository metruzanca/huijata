package cmd

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/metruzanca/huijata/internal/config"
	"github.com/metruzanca/huijata/internal/snapshots"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all saved snapshots",
	Long: `Removes every snapshot stored by huijata. Your save file itself is not
touched.`,
	RunE: runClear,
}

func runClear(cmd *cobra.Command, args []string) error {
	snapshotsDir, err := config.SnapshotsDir()
	if err != nil {
		return err
	}

	count, err := snapshots.Count(snapshotsDir)
	if err != nil {
		return err
	}
	if count == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No snapshots to clear.")
		return nil
	}

	confirm := true
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Delete all snapshots?").
				Description(fmt.Sprintf("This will permanently delete %d snapshot(s). Your save file won't be touched.", count)).
				Affirmative("Yes, delete them").
				Negative("Cancel").
				Value(&confirm),
		),
	).Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}
	if err != nil {
		return err
	}
	if !confirm {
		return nil
	}

	deleted, err := snapshots.Clear(snapshotsDir)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Deleted %d snapshot(s).\n", deleted)
	return nil
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
