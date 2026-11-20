// Package ddragon reads champion metadata from Riot's Data Dragon CDN.
package ddragon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

const cdn = "https://ddragon.leagueoflegends.com"

// Champion is the slice of Data Dragon metadata the game uses.
type Champion struct {
	ID   string   `json:"id"`   // asset id, e.g. "Kaisa"
	Name string   `json:"name"` // display name, matches Leaguepedia, e.g. "Kai'Sa"
	Tags []string `json:"tags"` // classes, e.g. ["Marksman"]
	// Range is "melee" or "ranged" (base attack range under 300 is melee).
	Range string `json:"range"`
}

type rawChampion struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Stats struct {
		AttackRange float64 `json:"attackrange"`
	} `json:"stats"`
}

func (r rawChampion) champion() Champion {
	// No damage type: Data Dragon's attack/magic ratings call Hwei, Vex or Fizz "mixed".
	c := Champion{ID: r.ID, Name: r.Name, Tags: r.Tags, Range: "melee"}
	if r.Stats.AttackRange >= 300 {
		c.Range = "ranged"
	}
	return c
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// Load returns the latest Data Dragon version and its champions sorted by name.
// champion.json is cached under cacheDir/<version>/.
func Load(cacheDir string) (string, []Champion, error) {
	body, err := get(cdn + "/api/versions.json")
	if err != nil {
		return "", nil, err
	}
	var versions []string
	if err := json.Unmarshal(body, &versions); err != nil || len(versions) == 0 {
		return "", nil, fmt.Errorf("versions.json: %v", err)
	}
	version := versions[0]

	path := filepath.Join(cacheDir, version, "champion.json")
	body, err = os.ReadFile(path)
	if err != nil {
		if body, err = get(fmt.Sprintf("%s/cdn/%s/data/en_US/champion.json", cdn, version)); err != nil {
			return "", nil, err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", nil, err
		}
		if err := os.WriteFile(path, body, 0o644); err != nil {
			return "", nil, err
		}
	}
	var r struct {
		Data map[string]rawChampion `json:"data"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return "", nil, err
	}
	champs := make([]Champion, 0, len(r.Data))
	for _, c := range r.Data {
		champs = append(champs, c.champion())
	}
	sort.Slice(champs, func(i, j int) bool { return champs[i].Name < champs[j].Name })
	return version, champs, nil
}
