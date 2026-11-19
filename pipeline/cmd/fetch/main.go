// fetch downloads a year of pro drafts from Leaguepedia and writes the
// complete ones to a single JSON file for the later pipeline steps.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mistahuman/draftle/pipeline/internal/draft"
	"github.com/mistahuman/draftle/pipeline/internal/leaguepedia"
)

func main() {
	year := flag.Int("year", 2026, "season year to download")
	cache := flag.String("cache", "data/raw", "directory for cached API pages")
	out := flag.String("out", "data/drafts.json", "output file")
	flag.Parse()

	games, err := leaguepedia.Fetch(*year, *cache)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch:", err)
		os.Exit(1)
	}
	complete := make([]draft.Game, 0, len(games))
	for _, g := range games {
		if g.Complete() {
			complete = append(complete, g)
		}
	}
	body, err := json.Marshal(complete)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*out, body, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "fetch:", err)
		os.Exit(1)
	}
	fmt.Printf("%d drafts fetched, %d complete → %s\n", len(games), len(complete), *out)
}
