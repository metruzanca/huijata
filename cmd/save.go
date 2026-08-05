package cmd

import (
	"fmt"
	"time"

	"github.com/metruzanca/huijata/internal/config"
	"github.com/metruzanca/huijata/internal/snapshots"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save <description>",
	Short: "Snapshot the current run so it can be restored later",
	Long: `Copies the files that make up the current run (player, world, streaks)
into a snapshot. Lifetime progression such as unlocked spells and stats is
left alone.`,
	Args: cobra.ExactArgs(1),
	RunE: runSave,
}

func runSave(cmd *cobra.Command, args []string) error {
	description := args[0]
	if err := validateDescription(description); err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	snapshotsDir, err := config.SnapshotsDir()
	if err != nil {
		return err
	}

	snap, err := snapshots.Create(cfg.SavePath(), snapshotsDir, description, time.Now())
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Saved snapshot %s: %s\n", snap.ID, snap.Description)
	return nil
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
