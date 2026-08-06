package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/metruzanca/huijata/internal/noitadata"
	"github.com/metruzanca/huijata/internal/player"
	"github.com/spf13/cobra"
)

var wandsCmd = &cobra.Command{
	Use:   "wands",
	Short: "Inspect the wands in the current save",
}

var wandsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the wands the player currently has",
	Long: `Reads player.xml from the current save and prints each wand the player
carries, including its deck and the spells loaded into it.`,
	RunE: runWandsShow,
}

func runWandsShow(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	p, err := player.Load(cfg.SavePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no player.xml found in the save — is there an active run?")
		}
		return err
	}

	names, err := noitadata.SpellNames("")
	if err != nil {
		names = nil
	}

	out := cmd.OutOrStdout()
	if len(p.Wands) == 0 {
		fmt.Fprintln(out, "No wands found.")
		return nil
	}

	for i, w := range p.Wands {
		fmt.Fprintf(out, "%d. %s\n", i+1, w.Name)
		fmt.Fprintf(out, "   gun level: %d\n", w.GunLevel)
		fmt.Fprintf(out, "   mana: %d/%d (+%d/s)\n", w.Mana, w.ManaMax, w.ManaChargeSpeed)
		fmt.Fprintf(out, "   deck: %d capacity, %d per round, %d reload, shuffle=%t\n",
			w.DeckCapacity, w.ActionsPerRound, w.ReloadTime, w.Shuffle)
		if len(w.Spells) == 0 {
			fmt.Fprintln(out, "   spells: (none)")
		} else {
			display := make([]string, len(w.Spells))
			for j, id := range w.Spells {
				if name, ok := names[id]; ok {
					display[j] = name
				} else {
					display[j] = id
				}
			}
			fmt.Fprintf(out, "   spells: %s\n", strings.Join(display, ", "))
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(wandsCmd)
	wandsCmd.AddCommand(wandsShowCmd)
}
