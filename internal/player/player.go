package player

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Wand is a wand item the player has in their inventory, extracted from the
// AbilityComponent of a nested "wand" entity in player.xml.
type Wand struct {
	Name            string
	GunLevel        int
	Mana            int
	ManaMax         int
	ManaChargeSpeed int
	DeckCapacity    int
	ActionsPerRound int
	ReloadTime      int
	Shuffle         bool
	Spells          []string
}

// Player is the parsed state of player.xml.
type Player struct {
	Wands []Wand
}

// Load reads player.xml from savePath (the save slot folder) and returns the
// wands currently held by the player.
func Load(savePath string) (*Player, error) {
	data, err := os.ReadFile(filepath.Join(savePath, "player.xml"))
	if err != nil {
		return nil, err
	}
	var root element
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing player.xml: %w", err)
	}

	p := &Player{}
	walk(&root, func(e *element) {
		if e.XMLName.Local == "Entity" && hasTag(e, "wand") {
			if w, ok := parseWand(e); ok {
				p.Wands = append(p.Wands, w)
			}
		}
	})
	return p, nil
}

type element struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Kids    []element  `xml:",any"`
}

func parseWand(e *element) (Wand, bool) {
	var ac *element
	for i := range e.Kids {
		if e.Kids[i].XMLName.Local == "AbilityComponent" {
			ac = &e.Kids[i]
			break
		}
	}
	if ac == nil {
		return Wand{}, false
	}

	w := Wand{}
	a := attrs(ac)
	w.Name = a["ui_name"]
	w.GunLevel = intVal(a["gun_level"])
	w.Mana = intVal(a["mana"])
	w.ManaMax = intVal(a["mana_max"])
	w.ManaChargeSpeed = intVal(a["mana_charge_speed"])

	for i := range ac.Kids {
		k := &ac.Kids[i]
		ka := attrs(k)
		switch k.XMLName.Local {
		case "gun_config":
			w.DeckCapacity = intVal(ka["deck_capacity"])
			w.ActionsPerRound = intVal(ka["actions_per_round"])
			w.ReloadTime = intVal(ka["reload_time"])
			w.Shuffle = ka["shuffle_deck_when_empty"] == "1"
		case "gunaction_config":
			if id := ka["action_id"]; id != "" {
				w.Spells = append(w.Spells, id)
			}
		}
	}

	w.Spells = append(w.Spells, collectSpells(e)...)
	return w, true
}

// collectSpells finds the spells loaded into a wand. Each loaded spell is
// serialized as a nested card_action entity (an ItemActionComponent holding
// the spell's action_id) inside the wand entity.
func collectSpells(e *element) []string {
	var spells []string
	var collect func(*element)
	collect = func(n *element) {
		for i := range n.Kids {
			k := &n.Kids[i]
			if k.XMLName.Local == "Entity" && hasTag(k, "card_action") {
				for j := range k.Kids {
					if k.Kids[j].XMLName.Local == "ItemActionComponent" {
						if id := attrs(&k.Kids[j])["action_id"]; id != "" {
							spells = append(spells, id)
						}
					}
				}
			}
			collect(k)
		}
	}
	collect(e)
	return spells
}

func walk(e *element, fn func(*element)) {
	fn(e)
	for i := range e.Kids {
		walk(&e.Kids[i], fn)
	}
}

func hasTag(e *element, want string) bool {
	for _, t := range strings.Split(attrs(e)["tags"], ",") {
		if strings.TrimSpace(t) == want {
			return true
		}
	}
	return false
}

func attrs(e *element) map[string]string {
	m := make(map[string]string, len(e.Attrs))
	for _, a := range e.Attrs {
		m[a.Name.Local] = a.Value
	}
	return m
}

func intVal(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
