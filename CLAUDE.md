# draftle — Claude Code context

Daily puzzle on real pro drafts. The MVP shows a draft at the last red pick (R5) and
asks for that pick. Started 2026-11-19 from an Italian brief that discussed three modes
(counterpick puzzle, guess the pro pick, full draft against the computer). The puzzle
mode and the guess-the-pro-pick mode turned out to be the same game when stopped at R5.

## Layout

| Path | What |
|---|---|
| `pipeline/` | Go module, stdlib only. Offline steps: fetch → (aggregate → export, not written yet) |
| `pipeline/internal/draft` | `Game` model: picks, roles, bans in draft order. Team 1 = blue, team 2 = red |
| `pipeline/internal/leaguepedia` | Cargo API client with an on-disk page cache |
| `pipeline/cmd/fetch` | `make data` → `pipeline/data/drafts.json` (complete drafts only) |
| `pipeline/cmd/spike` | `make spike`: baseline R5 predictability. The yardstick for any scorer |
| `web/` | Copy of `sveltekit-skeleton-starter`: adapter-static, GitHub Pages under `/draftle` |

There is **no backend**. The plan is that the pipeline precomputes everything (each
puzzle has ~45 candidates) and writes static JSON into `web/static/data/`.

## Data

- **Source: Leaguepedia Cargo API** (`lol.fandom.com/api.php`), tables `PicksAndBansS7`
  joined to `ScoreboardGames` on `GameId` for `Patch`. Licence CC BY-SA: attribution is
  in the footer. Oracle's Elixir was dropped because its old S3 URLs return 404.
- **Rate limiting is aggressive.** It cut in after about four quick queries. The client
  waits 10 s between pages and backs off on errors, and every page is cached under
  `pipeline/data/raw/<year>/`. Never delete the cache casually: a full 2026 fetch is
  23 pages.
- `pipeline/data/` is gitignored.
- Riot API (soloq matchups) is **not used**. A production key would be needed for a
  public app.

## Baseline (2026-11-19, 11,072 drafts, 30 patches)

How often the real R5 is in the top-k of a ranking. Candidates = champions picked in
that role and patch, minus anything already picked or banned (≈44 on average).
Counts are leave-one-out.

| ranking | top1 | top3 | top5 | top10 |
|---|---|---|---|---|
| meta (patch frequency) | 9.5% | 26.7% | 41.3% | 68.5% |
| lane opponent (year) | 11.9% | 31.2% | 46.9% | 72.2% |
| lane opponent + meta | 11.9% | 31.3% | 47.0% | 72.3% |

Reading: knowing the meta alone does not solve the puzzle, so there is a real game.
Lane opponent adds a little. Only 2.5% of R5s were never picked elsewhere in that
role and patch. A first-try guess will rarely be right, so the daily needs multiple
guesses with hints, not a single shot.

## Legal

Riot "Legal Jibber Jabber" disclaimer is in the footer. No "LoL" or "League" in the
name. No monetisation.

## Commands

See `make help`.
