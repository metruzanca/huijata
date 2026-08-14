package cmd

import (
	"fmt"
	"os"

	"github.com/metruzanca/huijata/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "huijata",
	Short:        "Huijata is a tool for the game noita",
	SilenceUsage: true,
}

// updateHelpLong fills in the root command's help with the resolved locations
// of huijata's config file and the Noita save folder, so `huijata --help`
// always shows where things live.
func updateHelpLong() {
	cfgPath, err := config.Path()
	if err != nil {
		cfgPath = "(could not determine)"
	}

	savePath := ""
	cfg, err := config.Load()
	if err == nil && cfg != nil {
		savePath = cfg.SavePath()
	} else {
		savePath = config.GuessGamePath()
		if savePath != "" {
			savePath += string(os.PathSeparator) + "save00"
		} else {
			savePath = "(not configured; run `huijata init`)"
		}
	}

	rootCmd.Long = fmt.Sprintf(`Huijata is a tool for the game noita.

Config file: %s
Save folder: %s
`, cfgPath, savePath)
}

func Execute() {
	updateHelpLong()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
