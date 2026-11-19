// spike measures how predictable the last red pick (R5) is from simple pro-play
// signals. It is the yardstick any scoring engine has to beat: if a ranking cannot
// put the real R5 near the top more often than these baselines, it adds nothing.
//
// Every count is leave-one-out: the game being ranked is removed from its own
// statistics before ranking.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/mistahuman/draftle/pipeline/internal/draft"
)

type counter map[string]map[string]map[string]int

func (c counter) add(a, b, champ string) {
	if c[a] == nil {
		c[a] = map[string]map[string]int{}
	}
	if c[a][b] == nil {
		c[a][b] = map[string]int{}
	}
	c[a][b][champ]++
}

type signals struct {
	freq    int // picks of the candidate in this role and patch
	freqTot int // all picks in this role and patch
	vsOpp   int // picks of the candidate in this role against this lane opponent, all year
}

type ranking struct {
	name  string
	score func(s signals) float64
}

var rankings = []ranking{
	{"meta (patch frequency)", func(s signals) float64 { return float64(s.freq) }},
	{"lane opponent (year)", func(s signals) float64 { return float64(s.vsOpp) + 1e-6*float64(s.freq) }},
	{"lane opponent + meta k=30", func(s signals) float64 {
		return float64(s.vsOpp) + 30*float64(s.freq)/float64(s.freqTot)
	}},
}

func main() {
	in := flag.String("in", "data/drafts.json", "drafts written by cmd/fetch")
	majorOnly := flag.Bool("major", false, "rank only games from major leagues and international events")
	majorStats := flag.Bool("major-stats", false, "count picks only in major games (implies -major)")
	flag.Parse()

	body, err := os.ReadFile(*in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spike:", err, "(run make data first)")
		os.Exit(1)
	}
	var games []draft.Game
	if err := json.Unmarshal(body, &games); err != nil {
		fmt.Fprintln(os.Stderr, "spike:", err)
		os.Exit(1)
	}

	var major []draft.Game
	for _, g := range games {
		if draft.Major(g.Overview) {
			major = append(major, g)
		}
	}
	statsGames, ranked := games, games
	if *majorOnly || *majorStats {
		ranked = major
	}
	if *majorStats {
		statsGames = major
	}

	freq := counter{}    // [patch][role][champ]
	matchup := counter{} // [role][opponent][champ]
	for _, g := range statsGames {
		for i := 0; i < 5; i++ {
			freq.add(g.Patch, g.Roles1[i], g.Picks1[i])
			freq.add(g.Patch, g.Roles2[i], g.Picks2[i])
			if o := g.LaneOpponent(1, i); o != "" {
				matchup.add(g.Roles1[i], o, g.Picks1[i])
			}
			if o := g.LaneOpponent(2, i); o != "" {
				matchup.add(g.Roles2[i], o, g.Picks2[i])
			}
		}
	}

	ks := []int{1, 3, 5, 10}
	fmt.Printf("%d drafts ranked, stats from %d\n\n%-28s", len(ranked), len(statsGames), "ranking R5 by")
	for _, k := range ks {
		fmt.Printf("  top%-3d", k)
	}
	fmt.Println()
	for _, r := range rankings {
		hits := map[int]int{}
		for _, g := range ranked {
			role, actual, opp := g.Roles2[4], g.Picks2[4], g.LaneOpponent(2, 4)
			taken := g.TakenBeforeR5()
			tot := 0
			for _, n := range freq[g.Patch][role] {
				tot += n
			}
			tot-- // leave-one-out
			type cand struct {
				champ string
				score float64
			}
			var pool []cand
			for c, n := range freq[g.Patch][role] {
				s := signals{freq: n, freqTot: tot, vsOpp: matchup[role][opp][c]}
				if c == actual {
					s.freq--
					s.vsOpp--
				}
				if s.freq > 0 && !taken[c] {
					pool = append(pool, cand{c, r.score(s)})
				}
			}
			sort.Slice(pool, func(i, j int) bool {
				if pool[i].score != pool[j].score {
					return pool[i].score > pool[j].score
				}
				return pool[i].champ < pool[j].champ
			})
			for i, p := range pool {
				if p.champ == actual {
					for _, k := range ks {
						if i < k {
							hits[k]++
						}
					}
					break
				}
			}
		}
		fmt.Printf("%-28s", r.name)
		for _, k := range ks {
			fmt.Printf("  %5.1f%%", 100*float64(hits[k])/float64(len(ranked)))
		}
		fmt.Println()
	}
}
