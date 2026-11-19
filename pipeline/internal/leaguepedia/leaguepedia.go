// Package leaguepedia downloads pro drafts from the Leaguepedia Cargo API
// (lol.fandom.com, CC BY-SA). Every page is cached on disk: the API rate-limits
// hard, so a rerun must never hit the network for pages it already has.
package leaguepedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mistahuman/draftle/pipeline/internal/draft"
)

const (
	api       = "https://lol.fandom.com/api.php"
	pageSize  = 500
	userAgent = "draftle-pipeline/0.1"
	pause     = 10 * time.Second
)

func fields() string {
	f := []string{"PB.GameId=gid", "PB.OverviewPage=ov", "PB.Winner=win", "SG.Patch=patch",
		"SG.Team1=t1", "SG.Team2=t2", "SG.DateTime_UTC=date"}
	for t := 1; t <= 2; t++ {
		for i := 1; i <= 5; i++ {
			f = append(f,
				fmt.Sprintf("PB.Team%dPick%d=p%d%d", t, i, t, i),
				fmt.Sprintf("PB.Team%dRole%d=r%d%d", t, i, t, i),
				fmt.Sprintf("PB.Team%dBan%d=b%d%d", t, i, t, i))
		}
	}
	return strings.Join(f, ",")
}

func query(year, offset int) url.Values {
	q := url.Values{}
	q.Set("action", "cargoquery")
	q.Set("format", "json")
	q.Set("tables", "PicksAndBansS7=PB,ScoreboardGames=SG")
	q.Set("join_on", "PB.GameId=SG.GameId")
	q.Set("fields", fields())
	q.Set("where", fmt.Sprintf("PB.OverviewPage LIKE '%%%d%%' AND PB.IsComplete='1'", year))
	q.Set("order_by", "PB.GameId")
	q.Set("limit", strconv.Itoa(pageSize))
	q.Set("offset", strconv.Itoa(offset))
	return q
}

func fetchPage(year, offset int) ([]byte, error) {
	q := query(year, offset)
	for attempt := 0; attempt < 8; attempt++ {
		req, err := http.NewRequest(http.MethodGet, api+"?"+q.Encode(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if !strings.Contains(string(body), `"error"`) {
			return body, nil
		}
		wait := time.Duration(60*(attempt+1)) * time.Second
		fmt.Fprintf(os.Stderr, "offset %d: API error, retrying in %s: %.120s\n", offset, wait, body)
		time.Sleep(wait)
	}
	return nil, fmt.Errorf("offset %d: giving up after repeated API errors", offset)
}

// Fetch returns every completed draft whose tournament page mentions year,
// reading cached pages from cacheDir and downloading only the missing ones.
func Fetch(year int, cacheDir string) ([]draft.Game, error) {
	dir := filepath.Join(cacheDir, strconv.Itoa(year))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	var games []draft.Game
	for offset := 0; ; offset += pageSize {
		path := filepath.Join(dir, fmt.Sprintf("page_%05d.json", offset))
		body, err := os.ReadFile(path)
		if err != nil {
			if body, err = fetchPage(year, offset); err != nil {
				return nil, err
			}
			if err := os.WriteFile(path, body, 0o644); err != nil {
				return nil, err
			}
			time.Sleep(pause)
		}
		page, err := Parse(body)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		games = append(games, page...)
		if len(page) < pageSize {
			return games, nil
		}
	}
}

// Parse turns one cargoquery response body into games.
func Parse(body []byte) ([]draft.Game, error) {
	var r struct {
		Cargoquery []struct {
			Title map[string]string `json:"title"`
		} `json:"cargoquery"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	games := make([]draft.Game, 0, len(r.Cargoquery))
	for _, row := range r.Cargoquery {
		m := row.Title
		win, _ := strconv.Atoi(m["win"])
		g := draft.Game{GameID: m["gid"], Patch: m["patch"], Overview: m["ov"], Winner: win,
			Team1: m["t1"], Team2: m["t2"], Date: m["date"]}
		for i := 0; i < 5; i++ {
			n := i + 1
			g.Picks1[i], g.Picks2[i] = m[fmt.Sprintf("p1%d", n)], m[fmt.Sprintf("p2%d", n)]
			g.Roles1[i], g.Roles2[i] = m[fmt.Sprintf("r1%d", n)], m[fmt.Sprintf("r2%d", n)]
			g.Bans1[i], g.Bans2[i] = m[fmt.Sprintf("b1%d", n)], m[fmt.Sprintf("b2%d", n)]
		}
		games = append(games, g)
	}
	return games, nil
}
