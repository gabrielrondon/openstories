# Product Hunt

**Name:** OpenStories

**Tagline** (60 chars max):
Production battle-scars for your AI coding agent

**Description** (260 chars max, 248 used):
Open-source library of user stories with the edge cases that broke real systems, cited from GitHub issues and post-mortems. Native MCP server for Claude Code and Cursor, plus a reality linter that scores your spec against known failures. Go, MIT.

**Topics:** Developer Tools, Open Source, Artificial Intelligence, GitHub

**Links:** https://openstories.tuturama.com · https://github.com/gabrielrondon/openstories

**Gallery order:** og.png, term-eval.png, ui-story-modal.png, term-get.png, term-install.png, ui-catalog.png

**First comment (maker):**

Hi PH. I build with coding agents every day and kept shipping the same class of bug: the happy path works, production disagrees. Magic links that enterprise mail scanners burn before anyone clicks. Webhooks that arrive out of order. Retries without idempotency keys that charge a customer twice. None of it is new; all of it is written down somewhere the agent never looks.

OpenStories is a single Go binary with the library embedded. Two commands register it as an MCP server for Claude Code or Cursor, and the agent can search stories, read the Gherkin criteria and the citations, and run `evaluate_spec` on its own plan before writing code. `openstories eval` does the same from the terminal or CI.

What it is not: a finished encyclopedia. Eight stories are hand-curated with verified evidence; the rest are generated variants across 20 industries and 104 domains so the taxonomy is complete, and the repo labels which is which. The harvester (`openstories harvest --repo`) turns real issues into new stories, and a story needs a real citation to be merged.

Tell me the production failure you wish your agent had known about. I will write the story.
