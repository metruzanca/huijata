package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/metruzanca/huijata/internal/stats"
	"github.com/metruzanca/huijata/internal/worldstate"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Show the current run's state",
	Long: `Reads world_state.xml and the latest session stats from the save and
prints the current run's world seed, status and special conditions.`,
	RunE: runRun,
}

func runRun(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	ws, err := worldstate.Load(cfg.SavePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no world_state.xml found in the save — is there an active run?")
		}
		return err
	}

	sec, err := currentRun(cfg.SavePath())
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	printKV(out, "World seed", sec["world_seed"])
	printKV(out, "Status", runStatus(sec))
	printKV(out, "Playtime", sec["playtime_str"])
	printKV(out, "Gold", sec["gold"])
	printKV(out, "Enemies killed", sec["enemies_killed"])
	printKV(out, "HP", sec["hp"])
	printKV(out, "Places visited", sec["places_visited"])
	if killed := sec["killed_by"]; killed != "" {
		printKV(out, "Killed by", killed)
	}
	printKV(out, "World day", ws.Attrs["day_count"])
	printKV(out, "NG+", yesNo(flagSet(cfg.SavePath(), "progress_ngplus")))
	printKV(out, "Nightmare", yesNo(flagSet(cfg.SavePath(), "progress_nightmare")))
	printKV(out, "Everything to gold", yesNo(ws.Attrs["EVERYTHING_TO_GOLD"] == "1"))
	printKV(out, "Infinite gold", yesNo(ws.Attrs["INFINITE_GOLD_HAPPENING"] == "1"))
	printKV(out, "Sun ending", yesNo(ws.Attrs["ENDING_HAPPINESS"] == "1"))
	return nil
}

// currentRun returns the stats of the run currently being played: the newest
// entry in stats/sessions, falling back to the <session> block of
// stats/_stats.salakieli when no session file has been written yet.
func currentRun(savePath string) (stats.Section, error) {
	sessions, err := stats.Sessions(savePath)
	if err != nil {
		return nil, err
	}
	if len(sessions) > 0 {
		return sessions[0].Stats, nil
	}
	st, err := stats.Load(savePath)
	if err != nil {
		return nil, err
	}
	return st.Session, nil
}

func runStatus(sec stats.Section) string {
	if sec["dead"] == "1" {
		return "dead"
	}
	return "alive"
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func flagSet(savePath, name string) bool {
	_, err := os.Stat(filepath.Join(savePath, "persistent", "flags", name))
	return err == nil
}

func printKV(out io.Writer, k, v string) {
	if v == "" {
		v = "-"
	}
	fmt.Fprintf(out, "%-18s  %s\n", k, v)
}

func init() {
	rootCmd.AddCommand(runCmd)
}
