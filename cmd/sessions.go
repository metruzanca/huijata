package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/metruzanca/huijata/internal/stats"
	"github.com/spf13/cobra"
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Show the run history from stats/sessions",
	RunE:  runSessionsList,
}

var sessionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recorded runs, newest first",
	RunE:  runSessionsList,
}

var sessionsByDayCmd = &cobra.Command{
	Use:   "by-day",
	Short: "Count how many runs happened on each day",
	RunE:  runSessionsByDay,
}

var sessionsLimit int

func init() {
	rootCmd.AddCommand(sessionsCmd)
	sessionsCmd.PersistentFlags().IntVarP(&sessionsLimit, "limit", "n", 20, "number of runs to show (0 = all)")
	sessionsCmd.AddCommand(sessionsListCmd)
	sessionsCmd.AddCommand(sessionsByDayCmd)
}

func runSessionsList(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	sessions, err := stats.Sessions(cfg.SavePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no stats/sessions folder found in the save")
		}
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%d runs recorded\n\n", len(sessions))
	if len(sessions) == 0 {
		return nil
	}
	if sessionsLimit > 0 && len(sessions) > sessionsLimit {
		sessions = sessions[:sessionsLimit]
	}

	fmt.Fprintf(out, "%-16s  %-6s  %-8s  %6s  %8s  %4s  %10s\n",
		"When", "Result", "Playtime", "Kills", "Gold", "Places", "Seed")
	for _, s := range sessions {
		result := "alive"
		if s.Stats["dead"] == "1" {
			result = "dead"
		}
		fmt.Fprintf(out, "%-16s  %-6s  %-8s  %6s  %8s  %4s  %10s\n",
			s.Time.Format("2006-01-02 15:04"),
			result,
			s.Stats["playtime_str"],
			s.Stats["enemies_killed"],
			s.Stats["gold"],
			s.Stats["places_visited"],
			s.Stats["world_seed"],
		)
	}
	return nil
}

func runSessionsByDay(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	sessions, err := stats.Sessions(cfg.SavePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no stats/sessions folder found in the save")
		}
		return err
	}

	counts := map[string]int{}
	for _, s := range sessions {
		counts[s.Time.Format("2006-01-02")]++
	}

	type dayCount struct {
		day   string
		count int
	}
	days := make([]dayCount, 0, len(counts))
	for day, count := range counts {
		days = append(days, dayCount{day, count})
	}
	sort.Slice(days, func(i, j int) bool {
		if days[i].count != days[j].count {
			return days[i].count > days[j].count
		}
		return days[i].day > days[j].day
	})

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%d runs across %d days\n\n", len(sessions), len(days))
	for _, d := range days {
		fmt.Fprintf(out, "%-10s  %5d\n", d.day, d.count)
	}
	return nil
}
