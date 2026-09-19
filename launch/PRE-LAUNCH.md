# Pre-launch gate

What I found while building the hotsite, in the order an HN reader will find it. None of this is visible from the README; all of it is visible in thirty seconds of `openstories list` and one `cat` of a story file.

## 1. 1,620 of the 1,628 stories are generated from 30 templates

`cmd/generator/catalog.go` holds 30 story templates (one to three per industry) and 10 "modifiers" (Network Partitions, Data Drift, Cold-Start Latency, ...). The generator multiplies template × modifier × domain. That is where "1,600+" comes from. Eight stories are hand-written (`OS-*-00N` with three digits): the magic-link one, the out-of-order webhooks one, the idempotency one, the tool-calling loop, the 429 cascade, partial-refund MDR, the SOC 2 audit trail, and the subdomain homoglyph one. Those eight are excellent and are what the hotsite showcases.

Consequences a reader will notice:

- Titles like "Double-entry immutable ledger journaling for tax-compliance transactions under Data Drift & Silent Schema Corruption", repeated across every fintech domain with only the domain word changed.
- The eval report for a weak spec lists the same two edge cases five times, once per sibling domain (see `site/data/playground.json`, sample 1).
- Every generated story is marked `status: verified`.

## 2. The evidence is 30 URLs reused up to 66 times each

`grep -h 'source:' stories -r | sort | uniq -c` shows 30 distinct URLs across 1,628 stories. Three of them (two HN items, one r/stripe thread) each back 66 stories. The generated quotes are template strings with the domain substituted ("We failed an enterprise procurement audit because our sso-saml system..."), attributed to a real issue that says something else. A reader who clicks one citation and finds it does not contain the quote will say so in the thread, and they will be right.

The eight curated stories cite real sources with real quotes (I checked the Stripe and HN ones for OS-DEV-001; both resolve).

## 3. Fixed already: the `%!(EXTRA string=...)` artefacts

Templates without a `%s` were passed through `fmt.Sprintf` with the domain, so 1,782 files and every eval report carried `%!(EXTRA string=<domain>)`. Fixed in [PR #1](https://github.com/gabrielrondon/openstories/pull/1) (generator helper, corpus regenerated, 162 stale duplicates removed, tests pass). Merge it before anything else.

## What I recommend before posting

Pick one; the copy in this kit works with either.

**A. Tell the truth about the shape (one afternoon).** Keep the corpus, change the labels. `status: verified` only on stories with a checked citation; generated ones become `status: template` or `generated`. README and hotsite say: "8 curated stories with verified evidence, expanded into 1,600+ domain variants by a generator you can run yourself." Reframe the number as coverage, not evidence. `eval` dedupes sibling-domain alerts (same template, same modifier) so the report shows one line per failure mode, not five. Quotes on generated stories drop the fabricated attribution and point at the template's source as "related discussion".

**B. Make the number true (weeks).** Run `openstories harvest` against real repos per domain, keep only stories whose quote appears in the linked issue, and let the count be whatever it is. A library of 150 real stories beats 1,600 synthetic ones on HN, and the harvester is already the most interesting part of the project.

Until one of these is done, the launch copy below stays in this folder. The hotsite can go live now: its numbers are real (1,628 loaded, 20 industries, 104 domains) and its showcase only uses the curated eight, but the "evidence-backed" claim in the hero is only as strong as item 2.
