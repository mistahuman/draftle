package leaguepedia

import "testing"

func TestParse(t *testing.T) {
	body := []byte(`{"cargoquery":[{"title":{"gid":"G1","ov":"LCK/2026 Season/Cup","win":"2","patch":"26.05","t1":"T1","t2":"Gen.G","date":"2026-08-14 09:00:00",
		"p11":"Orianna","r11":"Mid","b11":"Malphite","p25":"Poppy","r25":"Jungle","b25":"Ezreal"}}]}`)
	games, err := Parse(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("got %d games, want 1", len(games))
	}
	g := games[0]
	if g.GameID != "G1" || g.Patch != "26.05" || g.Winner != 2 ||
		g.Team1 != "T1" || g.Team2 != "Gen.G" || g.Date != "2026-08-14 09:00:00" {
		t.Errorf("header fields wrong: %+v", g)
	}
	if g.Picks1[0] != "Orianna" || g.Roles1[0] != "Mid" || g.Bans1[0] != "Malphite" {
		t.Errorf("B1 wrong: %q %q %q", g.Picks1[0], g.Roles1[0], g.Bans1[0])
	}
	if g.Picks2[4] != "Poppy" || g.Roles2[4] != "Jungle" || g.Bans2[4] != "Ezreal" {
		t.Errorf("R5 wrong: %q %q %q", g.Picks2[4], g.Roles2[4], g.Bans2[4])
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	if _, err := Parse([]byte("not json")); err == nil {
		t.Error("expected an error")
	}
}
