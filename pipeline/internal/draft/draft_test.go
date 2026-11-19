package draft

import "testing"

func sample() Game {
	return Game{
		Patch:  "26.05",
		Picks1: [5]string{"Orianna", "Vi", "Rumble", "Kai'Sa", "Rell"},
		Roles1: [5]string{"Mid", "Jungle", "Top", "Bot", "Support"},
		Picks2: [5]string{"Sion", "Syndra", "Lucian", "Braum", "Poppy"},
		Roles2: [5]string{"Top", "Mid", "Bot", "Support", "Jungle"},
		Bans1:  [5]string{"Malphite", "Neeko", "Azir", "", ""},
		Bans2:  [5]string{"Ahri", "Jax", "Nautilus", "Gnar", "Ezreal"},
	}
}

func TestLaneOpponent(t *testing.T) {
	g := sample()
	tests := []struct {
		team, i int
		want    string
	}{
		{2, 4, "Vi"},     // R5 jungle faces blue jungle
		{1, 0, "Syndra"}, // B1 mid faces red mid
		{2, 0, "Rumble"}, // R1 top faces blue top
	}
	for _, tt := range tests {
		if got := g.LaneOpponent(tt.team, tt.i); got != tt.want {
			t.Errorf("LaneOpponent(%d, %d) = %q, want %q", tt.team, tt.i, got, tt.want)
		}
	}
}

func TestTakenBeforeR5(t *testing.T) {
	g := sample()
	taken := g.TakenBeforeR5()
	if taken["Poppy"] {
		t.Error("R5 itself must not count as taken")
	}
	for _, c := range []string{"Orianna", "Braum", "Malphite", "Ezreal"} {
		if !taken[c] {
			t.Errorf("%s should be taken", c)
		}
	}
	if taken[""] {
		t.Error("empty ban slots must be skipped")
	}
	if len(taken) != 17 { // 9 picks + 8 real bans
		t.Errorf("len(taken) = %d, want 17", len(taken))
	}
}

func TestComplete(t *testing.T) {
	g := sample()
	if !g.Complete() {
		t.Fatal("sample should be complete")
	}
	g.Roles2[4] = ""
	if g.Complete() {
		t.Error("missing role should make the game incomplete")
	}
}
