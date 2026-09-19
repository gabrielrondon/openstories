---
id: OS-DEV-002
locale: en
industry: devtools
domain: auth-and-iam
title: "Preventing Magic Link Token Invalidation by Enterprise Email Security Scanners (Pre-fetching)"
demand_score: 9.4
status: verified
persona:
  role: "Security Engineer / Full-Stack Lead"
  context: "B2B SaaS applications serving enterprise customers equipped with automated mail defense suites (Microsoft Defender SafeLinks, Proofpoint, Mimecast)"
story:
  as_a: "Corporate user logging in via an emailed Magic Link"
  i_want: "The link in my email to avoid instant invalidation before my deliberate human click"
  so_that: "I can sign in without seeing the dreaded 'Token expired or already used' error caused by automated antivirus bots pre-fetching every URL"
acceptance_criteria:
  - scenario: "Automated GET request scanner pre-fetch by corporate antivirus"
    given: "A single-use authentication token link sent via email"
    when: "An automated enterprise mail security bot issues a background GET request to inspect the target page for phishing"
    then: "The server must render an intermediate landing page with a 'Confirm Sign-In' button WITHOUT burning or consuming the token on GET"
  - scenario: "Deliberate human confirmation click (POST request)"
    given: "The user has arrived at the confirmation landing page"
    when: "The user clicks the 'Confirm Sign-In' button submitting a POST request"
    then: "The token is consumed atomically, an authenticated session is established, and the token is invalidated"
edge_cases:
  - "Advanced email security crawlers with headless browsers simulating mouse movements or executing JavaScript."
  - "Strict 15-minute token TTL expiration regardless of click state."
  - "Client fingerprint matching (IP subnet or User-Agent heuristics) between token issuance and confirmation."
evidence:
  - source: "https://github.com/supabase/gotrue/issues/1204"
    type: "github_issue"
    quote: "Enterprise customers using Outlook SafeLinks / Microsoft Defender have all magic links pre-fetched. When users click it 5 seconds later, it says link is expired."
    date: "2024-03-22"
  - source: "https://reddit.com/r/webdev/comments/182k9p1"
    type: "reddit_thread"
    quote: "Never, ever consume a magic link token on a GET request. Mail security gateways fetch every link automatically and burn the token instantly."
    date: "2024-11-05"
evaluation_rubric:
  - "Does the initial GET handler explicitly avoid consuming or expiring the one-time authentication token?"
  - "Is there an explicit human interaction (POST form submission or button click) required to finalize authentication?"
  - "Are brute-force rate limiters active on the token verification endpoint?"
tags:
  - auth
  - magic-link
  - security
  - enterprise
---

# Production Failure Context

Hundreds of startups implement passwordless email authentication only to discover that 100% of their enterprise customers fail to log in. Enterprise email security solutions (Office 365 Advanced Threat Protection, Proofpoint URL Defense, Mimecast) crawl incoming links programmatically. If the endpoint `/auth/callback?token=XYZ` burns the token on the initial HTTP GET, the human recipient arrives at an already-expired link.
