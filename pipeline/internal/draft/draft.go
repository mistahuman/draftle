// Package draft holds the pro draft model shared by every pipeline step.
package draft

// Game is one completed pro draft. Team 1 is blue side, team 2 is red side.
// Picks and roles are in pick order: Picks1[0] is B1, Picks2[4] is R5.
type Game struct {
	GameID   string    `json:"gameId"`
	Patch    string    `json:"patch"`
	Overview string    `json:"overview"`
	Winner   int       `json:"winner"` // 1 or 2, 0 if unknown
	Picks1   [5]string `json:"picks1"`
	Picks2   [5]string `json:"picks2"`
	Roles1   [5]string `json:"roles1"`
	Roles2   [5]string `json:"roles2"`
	Bans1    [5]string `json:"bans1"`
	Bans2    [5]string `json:"bans2"`
}

// Complete reports whether every pick, role and the patch are present.
func (g Game) Complete() bool {
	if g.Patch == "" {
		return false
	}
	for i := 0; i < 5; i++ {
		if g.Picks1[i] == "" || g.Picks2[i] == "" || g.Roles1[i] == "" || g.Roles2[i] == "" {
			return false
		}
	}
	return true
}

// LaneOpponent returns the champion the given team's i-th pick faces in its role,
// or "" if the other team has nobody listed in that role.
func (g Game) LaneOpponent(team, i int) string {
	roles, otherRoles, otherPicks := g.Roles1, g.Roles2, g.Picks2
	if team == 2 {
		roles, otherRoles, otherPicks = g.Roles2, g.Roles1, g.Picks1
	}
	for j := 0; j < 5; j++ {
		if otherRoles[j] == roles[i] {
			return otherPicks[j]
		}
	}
	return ""
}

// TakenBeforeR5 is every champion banned or picked before the last red pick.
func (g Game) TakenBeforeR5() map[string]bool {
	taken := map[string]bool{}
	for i := 0; i < 5; i++ {
		for _, c := range []string{g.Bans1[i], g.Bans2[i], g.Picks1[i]} {
			if c != "" {
				taken[c] = true
			}
		}
		if i < 4 {
			taken[g.Picks2[i]] = true
		}
	}
	return taken
}
