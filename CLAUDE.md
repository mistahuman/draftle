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

**Look: a printed scorecard, vintage rather than new.** One theme, `draftle-theme`
(`src/lib/styles/draftle-theme.css`): cream paper / warm ink-brown surfaces, brass primary
(neither blue nor red, which belong to the sides), small radii. Type: Geist for everything
and Geist Mono for figures and labels (`font-mono`), self-hosted via fontsource (system
stacks fall back to `system-ui` on Linux). `font-display` is set in `@theme` in
`layout.css` and is Geist for now. Tried and rejected: Fraunces serif + Courier Prime,
with italics, read as generic "AI vintage". Also there: a faint SVG paper grain on `body`.
Style vocabulary: no italics, thin borders, double and dotted rules, tries as filled
coloured squares, outlined "stamp" badges, dashed outline for the R5 slot, a light sepia
on champion portraits. No theme picker; light/dark switch only. Logo: a minimap, three
lanes, blue base filled and red base still empty (`Logo.svelte`, mirrored with hard-coded
colours in `static/favicon.svg`). No wordmark: plain `draftle` text (a `draftl▢` wordmark
was tried and rejected). No nav: one page. Game colours that are not branding live as CSS
variables at the bottom of `src/routes/layout.css` (`--side-blue`, `--side-red`).
Components use theme presets and those tokens, never hard-coded colours.

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
- `DraftBoard.svelte` shows one row per lane (blue pick | role | red pick), so the R5 slot
  sits next to its lane opponent. Each pick keeps a small B1…R5 label for draft order.
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

## Skeleton docs

Skeleton has no official MCP server (community ones exist and are deliberately not used).
Before writing or changing Skeleton components or classes, read the official LLM docs:

- Index: `https://www.skeleton.dev/llms.txt` (~8 KB). Use the **Svelte** section.
- Then fetch only the page you need, e.g.
  `https://www.skeleton.dev/docs/svelte/framework-components/popover.md`.
- Avoid `llms-svelte.txt` (~720 KB) and `llms-full.txt` (~1.2 MB) unless a broad search
  is really needed.

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
