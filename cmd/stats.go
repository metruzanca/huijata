package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/metruzanca/huijata/internal/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show the player's saved stats",
	Long: `Decrypts stats/_stats.salakieli from the current save and prints the
player's lifetime, best-run and last-session stats, plus their top tracked
kills.`,
	RunE: runStats,
}

func runStats(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	st, err := stats.Load(cfg.SavePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no stats/_stats.salakieli found in the save")
		}
		return err
	}

	out := cmd.OutOrStdout()
	printSection(out, "Global", st.Global)
	printSection(out, "Best run", st.Highest)
	printSection(out, "Last session", st.Session)
	printTopKills(out, st.Entries)
	return nil
}

// statField is one line of a stats section: the display label and the
// attribute key the game writes.
type statField struct {
	label string
	key   string
}

var statFields = []statField{
	{"playtime", "playtime_str"},
	{"deaths", "death_count"},
	{"enemies killed", "enemies_killed"},
	{"gold", "gold"},
	{"items", "items"},
	{"places visited", "places_visited"},
	{"damage taken", "damage_taken"},
	{"healed", "healed"},
	{"heart containers", "heart_containers"},
	{"hp", "hp"},
	{"projectiles shot", "projectiles_shot"},
	{"killed by", "killed_by"},
	{"world seed", "world_seed"},
}

func printSection(out io.Writer, title string, sec stats.Section) {
	fmt.Fprintf(out, "%s\n", title)
	width := 0
	for _, f := range statFields {
		if len(f.label) > width {
			width = len(f.label)
		}
	}
	for _, f := range statFields {
		v, ok := sec[f.key]
		if !ok || v == "" {
			continue
		}
		fmt.Fprintf(out, "  %-*s  %s\n", width, f.label, v)
	}
	fmt.Fprintln(out)
}

func printTopKills(out io.Writer, entries []stats.Entry) {
	type kill struct {
		name  string
		count int64
	}
	var kills []kill
	for _, e := range entries {
		if strings.HasPrefix(e.Key, "action_") {
			continue
		}
		n, err := strconv.ParseInt(e.Value, 10, 64)
		if err != nil {
			continue
		}
		kills = append(kills, kill{e.Key, n})
	}
	if len(kills) == 0 {
		return
	}
	sort.Slice(kills, func(i, j int) bool { return kills[i].count > kills[j].count })
	if len(kills) > 10 {
		kills = kills[:10]
	}
	fmt.Fprintf(out, "Top kills\n")
	for _, k := range kills {
		fmt.Fprintf(out, "  %-30s %d\n", k.name, k.count)
	}
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
