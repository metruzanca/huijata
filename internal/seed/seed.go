package seed

import (
	"strconv"
)

type SeedInfo struct {
	WorldSeed uint32
	LC        []string
	AP        []string
	Fungal    []FungalShift
	PerkRows  [][]string
}

func Calculate(worldSeedStr string) (*SeedInfo, error) {
	ws, err := strconv.ParseUint(worldSeedStr, 10, 32)
	if err != nil {
		return nil, err
	}
	worldSeed := uint32(ws)

	lc, ap := AlchemyRecipe(worldSeed)
	fungal := FungalShifts(worldSeed, 20)
	perkRows, err := PerkRowsForSeed(worldSeed)
	if err != nil {
		return nil, err
	}

	return &SeedInfo{
		WorldSeed: worldSeed,
		LC:        lc,
		AP:        ap,
		Fungal:    fungal,
		PerkRows:  perkRows,
	}, nil
}
