package seed

import (
	"testing"
)

type perkRowTest struct {
	seed     uint32
	expected [][]string
}

var perkRowTests = []perkRowTest{
	{
		seed: 2046892791,
		expected: [][]string{
			{"MEGA_BEAM_STONE", "NO_MORE_SHUFFLE", "PERKS_LOTTERY"},
			{"ALWAYS_CAST", "FREEZE_FIELD", "MANA_FROM_KILLS"},
			{"EXTRA_KNOCKBACK", "PROTECTION_MELEE", "IRON_STOMACH"},
			{"TELEKINESIS", "EXTRA_MONEY", "EXTRA_PERK"},
			{"ORBIT", "DEATH_GHOST", "DISSOLVE_POWDERS"},
			{"PROJECTILE_REPULSION_SECTOR", "REVENGE_BULLET", "INVISIBILITY"},
			{"HUNGRY_GHOST", "ANGRY_LEVITATION", "RISKY_CRITICAL"},
		},
	},
	{
		seed: 123,
		expected: [][]string{
			{"LOW_HP_DAMAGE_BOOST", "PROJECTILE_HOMING", "MOLD"},
			{"BLEED_GAS", "RADAR_ENEMY", "FASTER_LEVITATION"},
			{"GLASS_CANNON", "PROJECTILE_REPULSION_SECTOR", "ORBIT"},
			{"TELEPORTITIS", "REMOVE_FOG_OF_WAR", "SAVING_GRACE"},
			{"RESPAWN", "PERSONAL_LASER", "TELEKINESIS"},
			{"FAST_PROJECTILES", "BOUNCE", "BLEED_OIL"},
			{"HOMUNCULUS", "PEACE_WITH_GODS", "PROJECTILE_EATER_SECTOR"},
		},
	},
}

func TestPerkRows(t *testing.T) {
	perks, err := LoadPerks()
	if err != nil {
		t.Fatalf("LoadPerks failed: %v", err)
	}

	for _, tc := range perkRowTests {
		t.Run("", func(t *testing.T) {
			deck := GeneratePerkDeck(perks, tc.seed)
			rows := DealPerkRows(deck, 7, 3)

			if len(rows) != len(tc.expected) {
				t.Errorf("seed %d: expected %d rows, got %d", tc.seed, len(tc.expected), len(rows))
				return
			}

			for i, row := range rows {
				if len(row) != len(tc.expected[i]) {
					t.Errorf("seed %d row %d: expected %d perks, got %d", tc.seed, i, len(tc.expected[i]), len(row))
					continue
				}
				for j, perk := range row {
					if perk != tc.expected[i][j] {
						t.Errorf("seed %d row %d perk %d: expected %s, got %s", tc.seed, i, j, tc.expected[i][j], perk)
					}
				}
			}
		})
	}
}

func TestPerkDeckDeterminism(t *testing.T) {
	perks, err := LoadPerks()
	if err != nil {
		t.Fatalf("LoadPerks failed: %v", err)
	}

	deck1 := GeneratePerkDeck(perks, 12345)
	deck2 := GeneratePerkDeck(perks, 12345)

	if len(deck1.Deck) != len(deck2.Deck) {
		t.Errorf("deck size mismatch between runs: %d vs %d", len(deck1.Deck), len(deck2.Deck))
	}

	for i := range deck1.Deck {
		if deck1.Deck[i] != deck2.Deck[i] {
			t.Errorf("deck content mismatch at index %d: %s vs %s", i, deck1.Deck[i], deck2.Deck[i])
		}
	}
}