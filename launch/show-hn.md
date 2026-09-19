# Show HN

**Title** (79 chars, HN cap is 80):

Show HN: OpenStories – Evidence-backed user stories and a reality linter for AI agents

**URL:** https://github.com/gabrielrondon/openstories

**Text:**

I run coding agents (Claude Code, Cursor) on real projects and kept hitting the same thing: the agent writes a clean happy path and has no idea what production already taught everyone else. Ask for a magic-link login and you get one that Microsoft Defender's link pre-fetch will break, because the one-time token gets burned before a human clicks. That failure has been documented in public GitHub issues for years. The agent has never read them.

OpenStories is an attempt to put that knowledge in the agent's context before it writes code. It is a Go binary (zero dependencies, library embedded) with three parts:

1. A library of user stories in Markdown with YAML front matter. Each has the As a / I want / So that statement, Gherkin acceptance scenarios, the observed production edge cases, and evidence: quotes from GitHub issues, post-mortems, HN and Reddit threads, with links. `openstories get OS-DEV-001` shows the idempotency-key one, with the $45k duplicate-charge incident behind it.

2. A native MCP server. `openstories install claude` or `openstories install cursor` registers it; the agent gets `search_stories`, `get_story`, `evaluate_spec`, `list_taxonomies` and `save_custom_story` as tools.

3. `openstories eval`, a reality linter for specs. It matches your spec text against the library and counts every known edge case the spec does not mention. "Charge endpoint that retries Stripe on timeout" scores 54/100 with three idempotency gaps listed; add the idempotency key, TTL, distributed lock and payload-mismatch handling and it scores 100. Same engine runs at `/api/eval` under `openstories serve`.

Honest notes. Eight stories are hand-written and carefully cited; the rest of the 1,628 are generated from 30 templates across 10 failure modes and 104 domains, so coverage is wide but the evidence density is uneven, and the labels in the repo say which is which. `openstories harvest --repo` mines a public repo's issues into new stories; that is the path to making the number real rather than wide. Contributions need a citation to be merged.

Site with a playground for the linter: https://openstories.tuturama.com

I would like to hear which production failure you wish your agent knew about. That is the next story.
