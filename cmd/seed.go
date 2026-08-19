package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/metruzanca/huijata/internal/noitadata"
	"github.com/metruzanca/huijata/internal/seed"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Show seed predictions for the current run",
	Long: `Shows seed-based predictions for the current active run:
- Current world seed
- Alchemic Precursor and Lively Concoction recipes
- Fungal shifts table (first 5)
- Mountain perk tables for all 7 holy mountains`,
	RunE: runSeed,
}

var shiftsCmd = &cobra.Command{
	Use:   "shifts",
	Short: "Show all 20 fungal shifts for the current run",
	Long:  `Displays all 20 possible fungal shifts for the current run's seed.`,
	RunE:  runShifts,
}

var holyMountainNames = []string{
	"Coal Pits",
	"Snowy Depths",
	"Frozen Floor / Secret",
	"Mine",
	"Dark Cave",
	"Temple of the Art",
	"Pyramid",
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.AddCommand(shiftsCmd)
}

func runSeed(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	sec, err := currentRun(cfg.SavePath())
	if err != nil {
		return err
	}

	worldSeed := sec["world_seed"]
	if worldSeed == "" {
		return fmt.Errorf("no world seed found in current run")
	}

	info, err := seed.Calculate(worldSeed)
	if err != nil {
		return fmt.Errorf("failed to calculate seed info: %w", err)
	}

	translations, err := getTranslations()
	out := cmd.OutOrStdout()
	if err != nil {
		fmt.Fprintf(out, "Note: could not load translations (%v); using raw material/perk IDs\n", err)
		translations = make(map[string]string)
	}

	translateMaterial := makeMaterialTranslator(translations)

	fmt.Fprintf(out, "Seed: %s\n\n", worldSeed)

	fmt.Fprintf(out, "=== Alchemic Precursor ===\n")
	fmt.Fprintf(out, "%s\n\n", strings.Join(translateMaterials(info.AP, translateMaterial), " • "))

	fmt.Fprintf(out, "=== Lively Concoction ===\n")
	fmt.Fprintf(out, "%s\n\n", strings.Join(translateMaterials(info.LC, translateMaterial), " • "))

	fmt.Fprintf(out, "=== Fungal Shifts (first 5 of %d) ===\n", len(info.Fungal))
	printFungalShifts(out, info.Fungal, 5, translateMaterial)
	fmt.Fprintln(out)

	fmt.Fprintf(out, "=== Mountain Perks ===\n")
	printPerkRows(out, info.PerkRows, translations)

	fmt.Fprintf(out, "\nHint: run 'huijata seed shifts' to see all %d fungal shifts\n", len(info.Fungal))

	return nil
}

func runShifts(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	sec, err := currentRun(cfg.SavePath())
	if err != nil {
		return err
	}

	worldSeed := sec["world_seed"]
	if worldSeed == "" {
		return fmt.Errorf("no world seed found in current run")
	}

	info, err := seed.Calculate(worldSeed)
	if err != nil {
		return fmt.Errorf("failed to calculate seed info: %w", err)
	}

	translations, err := getTranslations()
	out := cmd.OutOrStdout()
	if err != nil {
		fmt.Fprintf(out, "Note: could not load translations (%v); using raw material IDs\n", err)
		translations = make(map[string]string)
	}

	translateMaterial := makeMaterialTranslator(translations)

	fmt.Fprintf(out, "Seed: %s\n\n", worldSeed)
	fmt.Fprintf(out, "=== All Fungal Shifts (%d) ===\n", len(info.Fungal))
	printFungalShifts(out, info.Fungal, len(info.Fungal), translateMaterial)

	return nil
}

func printFungalShifts(out io.Writer, shifts []seed.FungalShift, limit int, translate func(string) string) {
	if limit > len(shifts) {
		limit = len(shifts)
	}
	fmt.Fprintf(out, "%-4s %-30s %-15s %-6s %-12s %-12s\n", "#", "From", "To", "Flask", "Gold→X", "Grass→X")
	for i := 0; i < limit; i++ {
		s := shifts[i]
		flask := ""
		if s.FlaskFrom {
			flask = "from"
		} else if s.FlaskTo {
			flask = "to"
		}
		from := strings.Join(translateMaterials(s.From, translate), ",")
		if len(from) > 28 {
			from = from[:28] + "…"
		}
		fmt.Fprintf(out, "%-4d %-30s %-15s %-6s %-12s %-12s\n",
			i+1, from, translate(s.To), flask, translate(s.GoldToX), translate(s.GrassToX))
	}
}

func printPerkRows(out io.Writer, rows [][]string, translations map[string]string) {
	names := loadPerkDisplayNames(translations)
	for i, row := range rows {
		hmName := "Holy Mountain"
		if i < len(holyMountainNames) {
			hmName = holyMountainNames[i]
		}
		displayPerks := make([]string, len(row))
		for j, perkID := range row {
			displayName := names[perkID]
			if displayName == "" {
				displayName = perkID
			}
			displayPerks[j] = displayName
		}
		fmt.Fprintf(out, "%s (%d): %s\n", hmName, i+1, strings.Join(displayPerks, " • "))
	}
}

func loadPerkDisplayNames(translations map[string]string) map[string]string {
	perks, err := seed.LoadPerks()
	if err != nil {
		return make(map[string]string)
	}
	names := make(map[string]string, len(perks))
	for _, p := range perks {
		key := strings.TrimPrefix(p.UIName, "$")
		if displayName, ok := translations[key]; ok {
			names[p.ID] = displayName
		} else {
			names[p.ID] = p.ID
		}
	}
	return names
}

func getTranslations() (map[string]string, error) {
	install, err := noitadata.DetectInstall()
	if err != nil {
		return nil, err
	}
	return noitadata.TranslationNames(install)
}

func makeMaterialTranslator(translations map[string]string) func(string) string {
	return func(raw string) string {
		if raw == "" {
			return raw
		}
		candidates := []string{
			"$material_" + raw,
			"material_" + raw,
			"$liquid_" + raw,
			"liquid_" + raw,
			"$" + raw,
			raw,
		}
		for _, key := range candidates {
			if name, ok := translations[key]; ok && name != "" {
				return name
			}
		}
		return raw
	}
}

func translateMaterials(list []string, translate func(string) string) []string {
	result := make([]string, len(list))
	for i, s := range list {
		result[i] = translate(s)
	}
	return result
}
