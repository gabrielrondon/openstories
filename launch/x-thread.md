# X thread (6 tweets)

Attach: 1 → `site/assets/og.png`; 2 → `term-eval.png`; 3 → `term-get.png`; 4 → `ui-story-modal.png`; 5 → `term-install.png`. Each tweet under 280 with the link counted as 23.

**1/**
Coding agents write the happy path. Production writes the rest.

OpenStories is an open-source library of user stories with the edge cases that broke real systems, plus an MCP server so Claude Code and Cursor read them before writing a line.

Go, MIT: github.com/gabrielrondon/openstories

**2/**
Ask an agent for "a charge endpoint that retries Stripe on timeout" and it will hand you a duplicate-charge bug.

openstories eval scores the spec against the library: 54/100, three idempotency edge cases missing, each one citing the incident that made someone write it down.

**3/**
Every story has the same shape: As a / I want / So that, Gherkin scenarios you can turn into tests, the observed production edge cases, and the evidence: GitHub issues, post-mortems, HN and Reddit threads, quoted and linked.

openstories get OS-DEV-001

**4/**
The one I keep coming back to: magic-link logins broken by enterprise mail scanners. Microsoft Defender pre-fetches the link and burns the one-time token before the human clicks. Documented in public issues for years. Agents still ship it.

**5/**
Setup is two commands:

go install github.com/gabrielrondon/openstories/cmd/openstories@latest
openstories install claude   (or: cursor)

The agent gets search_stories, get_story and evaluate_spec as tools. There is also a local dashboard: openstories serve.

**6/**
20 industries, 104 domains, 1,628 stories loaded. Stories are Markdown files with YAML front matter; a story needs a real citation to be merged.

Site: openstories.tuturama.com
Repo: github.com/gabrielrondon/openstories

Built at Tuturama.
