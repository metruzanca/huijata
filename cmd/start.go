package cmd

import (
	"fmt"

	"github.com/metruzanca/huijata/internal/launch"
	"github.com/metruzanca/huijata/internal/noitadata"
	"github.com/spf13/cobra"
)

var launchDirect bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch Noita",
	Long: `Launches Noita. Steam installs go through the Steam client so the run is
tracked; non-Steam installs launch the executable directly.`,
	RunE: runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	install, err := noitadata.DetectInstall()
	if err != nil {
		return err
	}

	viaSteam := !launchDirect && launch.IsSteamInstall(install)
	if viaSteam {
		fmt.Fprintln(cmd.OutOrStdout(), "Launching Noita through Steam...")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "Launching Noita directly...")
	}
	return launch.Launch(install, launchDirect)
}

func init() {
	startCmd.Flags().BoolVar(&launchDirect, "direct", false, "launch the Noita executable directly instead of through Steam")
	rootCmd.AddCommand(startCmd)
}
