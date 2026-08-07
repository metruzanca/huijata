package stats

import "testing"

func TestParse(t *testing.T) {
	data := `<GameStats STATS_VERSION="4" session_dead="0" >
  <KEY_VALUE_STATS>
    <E key="rat" value="12" ></E>
    <E key="action_burst_2" value="999" ></E>
  </KEY_VALUE_STATS>
  <session dead="1" gold="10" killed_by="lava" playtime_str="0:00:05" world_seed="7" ></session>
  <highest gold="500" enemies_killed="100" playtime_str="1:00:00" world_seed="42" ></highest>
  <global gold="1000" death_count="3" playtime_str="5:00:00" ></global>
  <prev_best gold="500" enemies_killed="100" ></prev_best>
</GameStats>`

	s, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if s.Version != 4 {
		t.Errorf("Version = %d, want 4", s.Version)
	}
	if s.SessionDead {
		t.Error("SessionDead = true, want false")
	}
	if len(s.Entries) != 2 {
		t.Fatalf("Entries len = %d, want 2", len(s.Entries))
	}
	if s.Entries[0] != (Entry{Key: "rat", Value: "12"}) {
		t.Errorf("Entries[0] = %+v, want {rat 12}", s.Entries[0])
	}
	if s.Session["gold"] != "10" || s.Session["killed_by"] != "lava" || s.Session["world_seed"] != "7" {
		t.Errorf("Session = %v", s.Session)
	}
	if s.Highest["world_seed"] != "42" {
		t.Errorf("Highest = %v", s.Highest)
	}
	if s.Global["death_count"] != "3" {
		t.Errorf("Global = %v", s.Global)
	}
	if s.PrevBest["enemies_killed"] != "100" {
		t.Errorf("PrevBest = %v", s.PrevBest)
	}
}

func TestParseNonGameStatsRoot(t *testing.T) {
	if _, err := Parse([]byte("<SomethingElse />")); err == nil {
		t.Fatal("expected an error for a non GameStats root")
	}
}
