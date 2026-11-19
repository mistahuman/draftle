# draftle — Claude Code context

Daily puzzle on real pro drafts. The MVP shows a draft at the last red pick (R5) and
asks for that pick. Started 2026-11-19 from an Italian brief that discussed three modes
(counterpick puzzle, guess the pro pick, full draft against the computer). The puzzle
mode and the guess-the-pro-pick mode turned out to be the same game when stopped at R5.

## Layout

| Path | What |
|---|---|
| `pipeline/` | Go module, stdlib only. Offline steps: fetch → export |
| `pipeline/internal/draft` | `Game` model: picks, roles, bans in draft order. Team 1 = blue, team 2 = red |
| `pipeline/internal/leaguepedia` | Cargo API client with an on-disk page cache |
| `pipeline/cmd/fetch` | `make data` → `pipeline/data/drafts.json` (complete drafts only) |
| `pipeline/internal/stats` | Leave-one-out pick counts by role/patch and by lane opponent; `R5Pool` |
| `pipeline/internal/ddragon` | Data Dragon champion list (id, name, class tags), cached per version |
| `pipeline/cmd/export` | `make export` → `ui/static/data/puzzles.json`: 365 puzzles from the last 6 patches, shuffled with a fixed seed, plus hint stats and the champion list |
| `pipeline/cmd/spike` | `make spike`: baseline R5 predictability. The yardstick for any scorer |
| `ui/` | Copy of `sveltekit-skeleton-starter`: adapter-static, GitHub Pages under `/draftle`. Deploy workflow is **manual only** until launch |

**Look: simple on purpose, restyled later by Marco.** The theme picker stays, as a testing
tool. `Logo.svelte` is the starter's placeholder until draftle has its own. No nav: one
page. Game colours that are not branding live as CSS variables at the bottom of
`src/routes/layout.css` (`--side-blue`, `--side-red`). Components use theme presets and
those tokens, never hard-coded colours.

Folder names follow `fullstack-webapp-starter` (`ui/`, later `backend/`), so adding a
server is additive: `backend/` from `gin-mongodb-starter` + compose + nginx, and `ui/`
builds as a Node server with `DOCKER_BUILD=true`. Do it only when a feature needs shared
state (community stats, accounts, leaderboards).

There is **no backend**. The pipeline precomputes everything and writes
`ui/static/data/puzzles.json`, which **is committed**: CI only builds the web app and
does not run Go. Rerun `make export` and commit after `make data`.

## The game (ui)

- `src/routes/+page.svelte`: the whole game. Daily puzzle = `dayNumber()` (day 1 is
  2026-11-19 UTC) modulo the puzzle count. The daily state is saved in `localStorage` under
  `draftle:day:<n>`. "Practice puzzle" (shown after the end) picks a random puzzle that is not saved.
- `src/lib/game/logic.ts`: rules. 6 tries. Each guess shows class match (Data Dragon
  tags) and whether the answer is more or less picked in that role and patch. One hint
  unlocks per miss (`unlockedHints`): class → picks against the lane opponent → meta rank →
  first letter → one of three. The class hint is skipped once a guess matched the class
  exactly. That is checked miss by miss, so a hint already shown never disappears.
- **Every number the player sees is per patch** (popularity, picks vs lane opponent, meta
  rank). `stats.Counts` counts per patch for that reason. Per patch, the vs-opponent count
  for an answer has a median of 4 and is 0 in about 20% of puzzles, and a 0 is a hint too.
- Champion names from Leaguepedia match Data Dragon `name` exactly (checked on
  2026-11-19, 173 champions). Only `None` / `Missing Data` placeholders differ. Icons load
  from the Data Dragon CDN at the version in the JSON.

## Legal

Riot "Legal Jibber Jabber" disclaimer is in the footer. No "LoL" or "League" in the
name. No monetisation.

## Commands

See `make help`.
