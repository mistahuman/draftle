package draft

import "testing"

func TestMajor(t *testing.T) {
	cases := map[string]bool{
		"LCK/2026 Season/Rounds 1-2":                                 true,
		"LPL/2026 Season/Split 3":                                    true,
		"LEC/2026 Season/Versus Season":                              true,
		"LCS/2026 Season/Lock-In":                                    true,
		"LCP/2026 Season/Split 1":                                    true,
		"CBLOL/2026 Season/Cup":                                      true,
		"2026 Season World Championship/Main Event":                  true,
		"2026 Season World Championship/Play-In":                     true,
		"2026 Mid-Season Invitational":                               true,
		"2026 First Stand":                                           true,
		"Esports World Cup 2026":                                     true,
		"Esports World Cup 2026/Online Qualifiers/EMEA":              false,
		"LCK CL/2026 Season/Rounds 1-2":                              false,
		"LCK Academy Series/2026 Season/2nd Championship":            false,
		"EMEA Masters/2026 Season/Summer":                            false,
		"Hitpoint 2nd Division Challengers/2026 Season/Summer Split": false,
		"2026 Asian Games":                                           false,
		"2026 World Star Challengers Invitational":                   false,
		"VCS/2026 Season/Summer Season":                              false,
	}
	for ov, want := range cases {
		if got := Major(ov); got != want {
			t.Errorf("Major(%q) = %v, want %v", ov, got, want)
		}
	}
}
