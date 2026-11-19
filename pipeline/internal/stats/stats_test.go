package stats

import (
	"testing"

	"github.com/mistahuman/draftle/pipeline/internal/draft"
)

func game(r5, oppJungle string) draft.Game {
	return draft.Game{
		Patch:  "26.05",
		Picks1: [5]string{"Orianna", oppJungle, "Rumble", "Kai'Sa", "Rell"},
		Roles1: [5]string{"Mid", "Jungle", "Top", "Bot", "Support"},
		Picks2: [5]string{"Sion", "Syndra", "Lucian", "Braum", r5},
		Roles2: [5]string{"Top", "Mid", "Bot", "Support", "Jungle"},
	}
}

func TestR5PoolLeaveOneOut(t *testing.T) {
	games := []draft.Game{game("Poppy", "Vi"), game("Poppy", "Vi"), game("Sejuani", "Vi")}
	pool := Build(games).R5Pool(games[0])

	got := map[string]Candidate{}
	for _, c := range pool {
		got[c.Champ] = c
	}
	// Poppy is picked twice; this game is one of them, so only one remains.
	if p := got["Poppy"]; p.Freq != 1 || p.VsOpp != 1 {
		t.Errorf("Poppy = %+v, want Freq 1 VsOpp 1", p)
	}
	if s := got["Sejuani"]; s.Freq != 1 || s.VsOpp != 1 {
		t.Errorf("Sejuani = %+v, want Freq 1 VsOpp 1", s)
	}
	// Vi is blue's jungler: picked in the role, but taken in this draft.
	if _, ok := got["Vi"]; ok {
		t.Error("Vi is already picked and must not be a candidate")
	}
}

func TestR5PoolDropsUnseenAnswer(t *testing.T) {
	games := []draft.Game{game("Poppy", "Vi"), game("Sejuani", "Vi")}
	for _, c := range Build(games).R5Pool(games[0]) {
		if c.Champ == "Poppy" {
			t.Error("Poppy only appears in the game itself and must drop out")
		}
	}
}

func TestR5PoolIsPerPatch(t *testing.T) {
	other := game("Poppy", "Vi")
	other.Patch = "26.04"
	games := []draft.Game{game("Sejuani", "Vi"), game("Poppy", "Vi"), other}
	for _, c := range Build(games).R5Pool(games[0]) {
		if c.Champ == "Poppy" && (c.Freq != 1 || c.VsOpp != 1) {
			t.Errorf("Poppy = %+v: the 26.04 game must not count in 26.05", c)
		}
	}
}
