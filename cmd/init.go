package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/metruzanca/huijata/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Point huijata at your Noita game folder",
	Long: `Finds the Nolla_Games_Noita folder that holds Noita's saves and stores
it in huijata's config file so later commands know where to look.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	guess := config.GuessGamePath()

	path := guess
	if guess != "" && config.DirExists(guess) {
		useGuess := true
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Found your Noita game folder").
					Description(fmt.Sprintf("Is this it?\n\n  %s", guess)).
					Affirmative("Yes, save it").
					Negative("No, I'll type it").
					Value(&useGuess),
			),
		).Run()
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		if err != nil {
			return err
		}
		if !useGuess {
			path = ""
		}
	}

	if path == "" {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Noita game folder location").
					Description("Enter the full path to the Nolla_Games_Noita folder that contains your saves (e.g. ...\\Nolla_Games_Noita).").
					Placeholder(guess).
					Value(&path).
					Validate(func(s string) error {
						if !config.DirExists(s) {
							return fmt.Errorf("%q is not a directory", s)
						}
						return nil
					}),
			),
		).Run()
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		if err != nil {
			return err
		}
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cfg := &config.Config{GamePath: abs}
	if err := cfg.Save(); err != nil {
		return err
	}

	cfgPath, err := config.Path()
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Game folder saved to %s\n", cfgPath)
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
