// export picks a deterministic set of R5 puzzles from recent patches and writes
// them, with their precomputed hints and the champion list, as one static JSON
// file for the web app.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mistahuman/draftle/pipeline/internal/ddragon"
	"github.com/mistahuman/draftle/pipeline/internal/draft"
	"github.com/mistahuman/draftle/pipeline/internal/stats"
)

type pick struct {
	Champ string `json:"champ"`
	Role  string `json:"role"`
}

type side struct {
	Picks []pick   `json:"picks"`
	Bans  []string `json:"bans"`
}

type puzzle struct {
	ID       string            `json:"id"`
	Patch    string            `json:"patch"`
	Event    string            `json:"event"`
	Winner   int               `json:"winner"`
	Blue     side              `json:"blue"`
	Red      side              `json:"red"` // first four picks only
	Role     string            `json:"role"`
	Answer   string            `json:"answer"`
	Opponent string            `json:"opponent"`
	MetaRank int               `json:"metaRank"` // answer's rank in the pool, 0 if absent
	PoolSize int               `json:"poolSize"`
	Stats    map[string][2]int `json:"stats"` // champ → [freq, vsOpp], pool only
}

type output struct {
	Version   string             `json:"version"`
	Champions []ddragon.Champion `json:"champions"`
	Puzzles   []puzzle           `json:"puzzles"`
}

func patchKey(p string) [2]int {
	parts := strings.SplitN(p, ".", 2)
	var k [2]int
	for i := range parts {
		k[i], _ = strconv.Atoi(parts[i])
	}
	return k
}

func cleanBans(bans [5]string) []string {
	out := []string{}
	for _, b := range bans {
		if b != "" && b != "None" && b != "Missing Data" {
			out = append(out, b)
		}
	}
	return out
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "export:", err)
	os.Exit(1)
}

func main() {
	in := flag.String("in", "data/drafts.json", "drafts written by cmd/fetch")
	out := flag.String("out", "../ui/static/data/puzzles.json", "output file")
	recent := flag.Int("patches", 6, "draw puzzles from the most recent N patches")
	n := flag.Int("n", 365, "number of puzzles")
	seed := flag.Int64("seed", 1, "shuffle seed; changing it reorders every daily puzzle")
	flag.Parse()

	body, err := os.ReadFile(*in)
	if err != nil {
		die(err)
	}
	var games []draft.Game
	if err := json.Unmarshal(body, &games); err != nil {
		die(err)
	}
	version, champs, err := ddragon.Load("data/ddragon")
	if err != nil {
		die(err)
	}
	known := map[string]bool{}
	for _, c := range champs {
		known[c.Name] = true
	}

	counts := stats.Build(games)

	patches := map[string]bool{}
	for _, g := range games {
		patches[g.Patch] = true
	}
	ordered := make([]string, 0, len(patches))
	for p := range patches {
		ordered = append(ordered, p)
	}
	sort.Slice(ordered, func(i, j int) bool {
		a, b := patchKey(ordered[i]), patchKey(ordered[j])
		return a[0] > b[0] || (a[0] == b[0] && a[1] > b[1])
	})
	if len(ordered) > *recent {
		ordered = ordered[:*recent]
	}
	keep := map[string]bool{}
	for _, p := range ordered {
		keep[p] = true
	}

	var eligible []draft.Game
	for _, g := range games {
		if !keep[g.Patch] {
			continue
		}
		ok := true
		for i := 0; i < 5; i++ {
			ok = ok && known[g.Picks1[i]] && known[g.Picks2[i]]
		}
		if ok {
			eligible = append(eligible, g)
		}
	}
	sort.Slice(eligible, func(i, j int) bool { return eligible[i].GameID < eligible[j].GameID })
	rand.New(rand.NewSource(*seed)).Shuffle(len(eligible), func(i, j int) {
		eligible[i], eligible[j] = eligible[j], eligible[i]
	})
	if len(eligible) > *n {
		eligible = eligible[:*n]
	}

	o := output{Version: version, Champions: champs}
	for _, g := range eligible {
		pool := counts.R5Pool(g)
		p := puzzle{
			ID: g.GameID, Patch: g.Patch, Event: g.Overview, Winner: g.Winner,
			Blue:     side{Bans: cleanBans(g.Bans1)},
			Red:      side{Bans: cleanBans(g.Bans2)},
			Role:     g.Roles2[4],
			Answer:   g.Picks2[4],
			Opponent: g.LaneOpponent(2, 4),
			PoolSize: len(pool),
			Stats:    map[string][2]int{},
		}
		for i := 0; i < 5; i++ {
			p.Blue.Picks = append(p.Blue.Picks, pick{g.Picks1[i], g.Roles1[i]})
			if i < 4 {
				p.Red.Picks = append(p.Red.Picks, pick{g.Picks2[i], g.Roles2[i]})
			}
		}
		for i, c := range pool {
			p.Stats[c.Champ] = [2]int{c.Freq, c.VsOpp}
			if c.Champ == p.Answer {
				p.MetaRank = i + 1
			}
		}
		o.Puzzles = append(o.Puzzles, p)
	}

	body, err = json.Marshal(o)
	if err != nil {
		die(err)
	}
	if err := os.WriteFile(*out, body, 0o644); err != nil {
		die(err)
	}
	fmt.Printf("%d puzzles from patches %v, %d champions (Data Dragon %s) → %s (%d KB)\n",
		len(o.Puzzles), ordered, len(champs), version, *out, len(body)/1024)
}
