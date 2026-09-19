---
id: OS-SaaS-001
locale: en
industry: b2b-saas
domain: multi-tenancy
title: "Preventing Reserved Keyword Collisions, Subdomain Hijacking, and Homoglyph Spoofing in Multi-Tenant SaaS"
demand_score: 9.3
status: verified
persona:
  role: "Chief Information Security Officer / Principal Architect"
  context: "Multi-tenant B2B SaaS platforms provisioning vanity subdomains (tenant.saas.com) or organization URL paths"
story:
  as_a: "SaaS platform security administrator"
  i_want: "The workspace registration flow to rigorously validate an exhaustive reserved keyword denylist, system infrastructure terms, and Unicode homoglyphs"
  so_that: "Malicious actors cannot claim reserved slugs such as 'admin', 'api', 'billing', 'oauth', or 'support' to steal session cookies or launch targeted spear-phishing campaigns"
acceptance_criteria:
  - scenario: "Attempting to claim an internal infrastructure keyword slug"
    given: "A user registering a new workspace organization"
    when: 'The user submits the slug "admin", "billing", "auth", or "api"'
    then: "The system must reject the registration with HTTP 400 Bad Request and an explicit error indicating the identifier is reserved by system policy"
  - scenario: "Attempting homoglyph spoofing via lookalike Unicode characters (e.g. Cyrillic 'а' replacing Latin 'a')"
    given: 'A legitimate enterprise customer is registered with slug "paypal"'
    when: 'An attacker attempts to register "pаypal" using visually indistinguishable Cyrillic runes'
    then: "The system must normalize via NFKC/Punycode and strictly enforce ASCII alphanumeric constraints (regex ^[a-z0-9-]+$)"
edge_cases:
  - "Workspace slug rename quarantine (decommissioned slugs must enter a 90-day cooldown before being released to the public)."
  - "Internal DNS collision risks (e.g. autodiscover, smtp, imap, ns1, cdn)."
  - "Parent domain cookie scoping (ensuring cookies defined on .saas.com do not leak into user-controlled tenant subdomains without Public Suffix List isolation)."
evidence:
  - source: "https://hackerone.com/reports/409850"
    type: "security_bounty_report"
    quote: "Attacker claimed the subdomain 'support.target.com' by registering the workspace slug 'support', allowing credential harvesting through cross-subdomain cookies."
    date: "2023-11-20"
  - source: "https://news.ycombinator.com/item?id=34281902"
    type: "hackernews"
    quote: "If your SaaS permits custom subdomain creation without an exhaustive denylist of over 500 reserved keywords, you have a critical security flaw from day one."
    date: "2024-04-10"
evaluation_rubric:
  - "Does the registration validator enforce an exhaustive system keyword denylist (>200 keywords)?"
  - "Is there strict alphanumeric charset enforcement preventing internationalized domain (IDN) homoglyph spoofing?"
  - "Are renamed or deleted tenant slugs placed into a security quarantine to prevent immediate takeover?"
tags:
  - security
  - multi-tenancy
  - saas
  - domain-takeover
  - isolation
---

# Multi-Tenant Isolation Hazards

In subdomain-based multi-tenant architectures (`{tenant}.company.com`), one of the most common critical findings reported by bug bounty hunters is the lack of reserved subdomain protections. Failing to reserve names like `api`, `auth`, `status`, or `billing` allows attackers to intercept cookies or manipulate single-sign-on (SSO) redirect flows.
