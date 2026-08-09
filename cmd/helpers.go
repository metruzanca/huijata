package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/metruzanca/huijata/internal/config"
	"github.com/metruzanca/huijata/internal/launch"
	"github.com/metruzanca/huijata/internal/noitadata"
	"github.com/spf13/cobra"
)

// loadConfig reads the config, failing with a hint to run init when missing.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("no config found. Run `huijata init` first.")
	}
	return cfg, nil
}

func validateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return errors.New("description must not be empty")
	}
	if strings.ContainsAny(description, `/\`) {
		return errors.New("description must not contain path separators")
	}
	return nil
}

// maybeStartNoita launches Noita unless the game is already running. When ask
// is true the user is prompted first; when false the game is launched without
// asking. Failures to detect or launch the game are reported but do not fail
// the operation that already completed.
func maybeStartNoita(cmd *cobra.Command, ask bool) error {
	if launch.IsRunning() {
		return nil
	}

	if ask {
		start, err := askToStart()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}
		if !start {
			return nil
		}
	}

	install, err := noitadata.DetectInstall()
	if err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "Could not find the Noita install to launch it: %v\n", err)
		return nil
	}
	if err := launch.Launch(install, false); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "Noita failed to launch: %v\n", err)
		return nil
	}
	return nil
}

// askToStart shows a confirm prompt asking whether to launch Noita now. The
// default is yes; esc answers no.
func askToStart() (bool, error) {
	start := true
	km := huh.NewDefaultKeyMap()
	km.Confirm.Reject.SetKeys("n", "N", "esc")
	km.Confirm.Reject.SetHelp("esc", "No")
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Start Noita now?").
				Description("Launch the game to continue playing?").
				Affirmative("Yes, start Noita").
				Negative("No").
				Value(&start),
		),
	).WithKeyMap(km).Run()
	return start, err
}
