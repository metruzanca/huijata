package player

import (
	"os"
	"path/filepath"
	"testing"
)

const samplePlayerXML = `<Entity 
  name="DEBUG_NAME:player" 
  tags="mortal,player_unit" >
  <Inventory2Component _enabled="1" >
  </Inventory2Component>
  <Entity 
    name="" 
    tags="teleportable_NOT,item,wand" >
    <AbilityComponent 
      ui_name="Bolt staff" 
      gun_level="1" 
      mana="104" 
      mana_max="104" 
      mana_charge_speed="26" >
      <gun_config 
        actions_per_round="1" 
        deck_capacity="3" 
        reload_time="23" 
        shuffle_deck_when_empty="0" >
      </gun_config>
      <gunaction_config action_id="" >
      </gunaction_config>
    </AbilityComponent>
    <Entity 
      name="" 
      tags="card_action" >
      <ItemActionComponent action_id="LIGHTNING_BOLT" >
      </ItemActionComponent>
    </Entity>
    <Entity 
      name="" 
      tags="card_action" >
      <ItemActionComponent action_id="BURST_2" >
      </ItemActionComponent>
    </Entity>
  </Entity>
  <Entity 
    name="" 
    tags="card_action" >
    <AbilityComponent ui_name="$action_bomb" >
    </AbilityComponent>
  </Entity>
</Entity>`

func writePlayer(t *testing.T, content string) string {
	t.Helper()
	save := t.TempDir()
	if err := os.WriteFile(filepath.Join(save, "player.xml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return save
}

func TestLoadReturnsWands(t *testing.T) {
	p, err := Load(writePlayer(t, samplePlayerXML))
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Wands) != 1 {
		t.Fatalf("len(Wands) = %d, want 1", len(p.Wands))
	}

	w := p.Wands[0]
	if w.Name != "Bolt staff" {
		t.Errorf("Name = %q, want %q", w.Name, "Bolt staff")
	}
	if w.GunLevel != 1 || w.Mana != 104 || w.ManaMax != 104 || w.ManaChargeSpeed != 26 {
		t.Errorf("unexpected ability stats: %+v", w)
	}
	if w.DeckCapacity != 3 || w.ActionsPerRound != 1 || w.ReloadTime != 23 || w.Shuffle {
		t.Errorf("unexpected gun_config: %+v", w)
	}
	if len(w.Spells) != 2 || w.Spells[0] != "LIGHTNING_BOLT" || w.Spells[1] != "BURST_2" {
		t.Errorf("Spells = %v, want [LIGHTNING_BOLT BURST_2]", w.Spells)
	}
}

func TestLoadIgnoresNonWandItems(t *testing.T) {
	p, err := Load(writePlayer(t, samplePlayerXML))
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range p.Wands {
		if w.Name == "$action_bomb" {
			t.Error("card_action item must not be treated as a wand")
		}
	}
}

func TestLoadEmptyWandDeck(t *testing.T) {
	xml := `<Entity name="player" tags="mortal">
  <Entity name="" tags="item,wand">
    <AbilityComponent ui_name="Bomb wand" gun_level="1">
      <gun_config actions_per_round="1" deck_capacity="1" reload_time="8" shuffle_deck_when_empty="1">
      </gun_config>
    </AbilityComponent>
    <Entity name="" tags="card_action">
      <ItemActionComponent action_id="">
      </ItemActionComponent>
    </Entity>
  </Entity>
</Entity>`

	p, err := Load(writePlayer(t, xml))
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Wands) != 1 {
		t.Fatalf("len(Wands) = %d, want 1", len(p.Wands))
	}
	w := p.Wands[0]
	if w.Name != "Bomb wand" || !w.Shuffle {
		t.Errorf("unexpected wand: %+v", w)
	}
	if len(w.Spells) != 0 {
		t.Errorf("Spells = %v, want empty", w.Spells)
	}
}

func TestLoadMissingPlayerXML(t *testing.T) {
	_, err := Load(t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing player.xml")
	}
}
