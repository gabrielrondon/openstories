# OpenStories hotsite

Static site for **https://openstories.tuturama.com**. No framework, no build step: `index.html`, `styles.css`, `app.js`, plus JSON data and PNG assets. Deployed on Vercel as project `openstories-site` (team `gabrielrondons-projects`).

## Layout

```
site/
  index.html          page (hero, problem, playground, incidents, MCP setup, CLI, dashboard, get)
  styles.css          violet palette on #0b0b10, Inter Tight + JetBrains Mono, reduced-motion safe
  app.js              counters, copy buttons, playground (local evaluator or live API), showcase cards
  vercel.json         clean URLs, cache headers for /assets and /data, security headers
  data/
    stories.min.json  compact corpus (fields the Go store scores on), built by scripts/build-data.py
    taxonomies.json   industries and domains, from /api/taxonomies
    playground.json   7 sample specs with their precomputed /api/eval result
    showcase.json     the 5 curated stories shown in "Incident showcase" (hand-picked)
  assets/
    og.png            1200x630 social banner, rendered from assets-src/og.html
    term-*.png        terminal renders from assets-src/terminal.html (2x, 1202 css px wide)
    ui-*.png          dashboard captures of `openstories serve`
    favicon.svg
  assets-src/         HTML sources for og.png and the terminal renders
  scripts/build-data.py
```

## The playground is the real evaluator

`app.js` mirrors `internal/store/store.go` (`tokenize`, `calculateScore`, `Search`) and `internal/evaluator/evaluator.go` (`Evaluate`, `extractKeyPhrases`) line for line, running over `data/stories.min.json` in the browser. Parity was checked against `POST /api/eval` on 13 specs: identical score and alert count on all 13. Matched story IDs can differ only among tied siblings, because Go's `sort.Slice` is not stable and the generated sibling stories score identically.

The 7 sample chips use `playground.json`, which is the API's own output. Free text runs the in-browser mirror. If `window.OPENSTORIES_API` is set to a URL in `index.html`, free text goes to that API instead.

`window.OpenStoriesEval` exposes `tokenize`, `search`, `evaluateLocal` and `corpus()` for debugging in the console.

## Rebuilding the data

Whenever the corpus changes (new stories, regenerated stories):

```bash
# from the repo root, with the current binary
./openstories serve --port 8090 &
python3 site/scripts/build-data.py http://127.0.0.1:8090
```

This rewrites `stories.min.json`, `taxonomies.json` and re-runs the 7 playground samples through `/api/eval`. Then update the three counters in `index.html` (hero `data-count` values and the `<title>`) and in `assets-src/og.html` if the totals moved, and re-render `og.png`.

## Re-rendering images

Serve the folder and render with headless Chrome or Playwright:

```bash
cd site && python3 -m http.server 8091 &
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless=new --hide-scrollbars \
  --force-device-scale-factor=1 --window-size=1200,630 \
  --screenshot="$PWD/assets/og.png" http://127.0.0.1:8091/assets-src/og.html
```

Terminal renders are element screenshots of `#list .term`, `#get .term`, `#eval .term`, `#install .term` in `assets-src/terminal.html` at device scale 2. The text in that file is copied from real command output; when the CLI output changes, update the HTML first.

## Deploying

```bash
cd site
vercel deploy --prod --yes --scope gabrielrondons-projects
```

The folder is linked to the project (`.vercel/`, git-ignored). DNS: `openstories CNAME cname.vercel-dns.com` at Hostinger, domain attached to the project with `vercel domains add openstories.tuturama.com openstories-site`.

## Local preview

```bash
cd site && python3 -m http.server 8091
```

Open http://127.0.0.1:8091. To test free-text evaluation against the Go server instead of the mirror, run `openstories serve --port 8090` and set `window.OPENSTORIES_API = "http://127.0.0.1:8090"` in `index.html` (the Go server must allow the origin).
