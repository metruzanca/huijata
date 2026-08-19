package seed

import (
	_ "embed"
	"encoding/json"
)

//go:embed perks.json
var perksJSON []byte

type PerkData struct {
	ID                       string  `json:"id"`
	UIName                   string  `json:"ui_name"`
	UIDescription            string  `json:"ui_description"`
	Stackable                bool    `json:"stackable"`
	StackableIsRare          bool    `json:"stackable_is_rare"`
	StackableMaximum         int     `json:"stackable_maximum"`
	MaxInPerkPool            int     `json:"max_in_perk_pool"`
	StackableHowOftenReappears int   `json:"stackable_how_often_reappears"`
	NotInDefaultPerkPool     bool    `json:"not_in_default_perk_pool"`
	RemoveOtherPerks         []string `json:"remove_other_perks,omitempty"`
}

type PerkDeck struct {
	Deck           []string
	StackableCount map[string]int
}

func LoadPerks() ([]PerkData, error) {
	var perks []PerkData
	if err := json.Unmarshal(perksJSON, &perks); err != nil {
		return nil, err
	}
	return perks, nil
}

func GeneratePerkDeck(perks []PerkData, worldSeed uint32) PerkDeck {
	rng := &NollaPrng{}
	rng.SetWorldSeed(worldSeed)
	rng.SetRandomSeed(1, 2)

	const minDistance = 4
	const defaultMaxStackable = 128

	perkDeck := []string{}
	stackableDistances := map[string]int{}
	stackableCount := map[string]int{}

	for _, perkData := range perks {
		if perkData.NotInDefaultPerkPool {
			continue
		}

		perkName := perkData.ID
		howManyTimes := 1
		stackableDistances[perkName] = -1
		stackableCount[perkName] = -1

		if perkData.Stackable {
			maxPerks := rng.Random(1, 2)
			if perkData.MaxInPerkPool > 0 {
				maxPerks = rng.Random(1, int32(perkData.MaxInPerkPool))
			}

			if perkData.StackableMaximum > 0 {
				stackableCount[perkName] = perkData.StackableMaximum
			} else {
				stackableCount[perkName] = defaultMaxStackable
			}

			if perkData.StackableIsRare {
				maxPerks = 1
			}

			if perkData.StackableHowOftenReappears > 0 {
				stackableDistances[perkName] = perkData.StackableHowOftenReappears
			} else {
				stackableDistances[perkName] = minDistance
			}

			howManyTimes = int(rng.Random(1, maxPerks))
		}

		for j := 0; j < howManyTimes; j++ {
			perkDeck = append(perkDeck, perkName)
		}
	}

	shuffleTable(rng, perkDeck)

	for i := len(perkDeck) - 1; i >= 1; i-- {
		perk := perkDeck[i]
		minDist := stackableDistances[perk]
		if minDist != -1 {
			removeMe := false
			for ri := i - minDist; ri < i; ri++ {
				if ri >= 0 && perkDeck[ri] == perk {
					removeMe = true
					break
				}
			}
			if removeMe {
				perkDeck = append(perkDeck[:i], perkDeck[i+1:]...)
			}
		}
	}

	return PerkDeck{Deck: perkDeck, StackableCount: stackableCount}
}

func shuffleTable(rng Rng, deck []string) {
	iterations := len(deck) - 1
	for i := iterations; i >= 1; i-- {
		j := int(rng.Random(0, int32(i)))
		tmp := deck[i]
		deck[i] = deck[j]
		deck[j] = tmp
	}
}

type PerkRow struct {
	Perks []string
}

func DealPerkRows(deck PerkDeck, numRows int, perksPerRow int) [][]string {
	result := [][]string{}
	nextIndex := 0
	deckCopy := make([]string, len(deck.Deck))
	copy(deckCopy, deck.Deck)

	for row := 0; row < numRows; row++ {
		rowPerks := []string{}
		for i := 0; i < perksPerRow; i++ {
			perkID := getNextPerk(deckCopy, &nextIndex)
			rowPerks = append(rowPerks, perkID)
		}
		result = append(result, rowPerks)
	}
	return result
}

func getNextPerk(deck []string, nextIndex *int) string {
	for {
		if *nextIndex >= len(deck) {
			*nextIndex = 0
		}
		perkID := deck[*nextIndex]
		if perkID != "" {
			deck[*nextIndex] = ""
			*nextIndex++
			if *nextIndex >= len(deck) {
				*nextIndex = 0
			}
			return perkID
		}
		*nextIndex++
		if *nextIndex >= len(deck) {
			*nextIndex = 0
		}
	}
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

func PerkRowsForSeed(worldSeed uint32) ([][]string, error) {
	perks, err := LoadPerks()
	if err != nil {
		return nil, err
	}

	deck := GeneratePerkDeck(perks, worldSeed)
	return DealPerkRows(deck, 7, 3), nil
}
