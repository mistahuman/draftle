// Package stats counts pro picks so a draft can be ranked or hinted.
// Every lookup is leave-one-out: the game being asked about is removed from
// its own counts, so a hint never leaks the answer's own pick.
package stats

import (
	"sort"

	"github.com/mistahuman/draftle/pipeline/internal/draft"
)

type nested map[string]map[string]map[string]int

func (n nested) add(a, b, champ string) {
	if n[a] == nil {
		n[a] = map[string]map[string]int{}
	}
	if n[a][b] == nil {
		n[a][b] = map[string]int{}
	}
	n[a][b][champ]++
}

// Counts holds pick frequencies across a set of games.
type Counts struct {
	freq    nested // [patch][role][champ]
	matchup nested // [role][lane opponent][champ]
}

// Build counts every pick in games, by role and patch and by lane opponent.
func Build(games []draft.Game) *Counts {
	c := &Counts{freq: nested{}, matchup: nested{}}
	for _, g := range games {
		for i := 0; i < 5; i++ {
			c.freq.add(g.Patch, g.Roles1[i], g.Picks1[i])
			c.freq.add(g.Patch, g.Roles2[i], g.Picks2[i])
			if o := g.LaneOpponent(1, i); o != "" {
				c.matchup.add(g.Roles1[i], o, g.Picks1[i])
			}
			if o := g.LaneOpponent(2, i); o != "" {
				c.matchup.add(g.Roles2[i], o, g.Picks2[i])
			}
		}
	}
	return c
}

// Candidate is one champion that could have been the R5 pick.
type Candidate struct {
	Champ string
	Freq  int // picks in R5's role and patch, excluding this game
	VsOpp int // picks in R5's role against R5's lane opponent, excluding this game
}

// R5Pool returns every champion picked in R5's role and patch elsewhere, minus
// anything already taken in this draft, most picked first.
func (c *Counts) R5Pool(g draft.Game) []Candidate {
	role, actual, opp := g.Roles2[4], g.Picks2[4], g.LaneOpponent(2, 4)
	taken := g.TakenBeforeR5()
	var pool []Candidate
	for champ, n := range c.freq[g.Patch][role] {
		cand := Candidate{Champ: champ, Freq: n, VsOpp: c.matchup[role][opp][champ]}
		if champ == actual {
			cand.Freq--
			cand.VsOpp--
		}
		if cand.Freq > 0 && !taken[champ] {
			pool = append(pool, cand)
		}
	}
	sort.Slice(pool, func(i, j int) bool {
		if pool[i].Freq != pool[j].Freq {
			return pool[i].Freq > pool[j].Freq
		}
		return pool[i].Champ < pool[j].Champ
	})
	return pool
}
