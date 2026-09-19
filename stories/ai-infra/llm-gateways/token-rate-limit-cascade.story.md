---
id: OS-AI-002
locale: en
industry: ai-infra
domain: llm-gateways
title: "Graceful Multi-Provider Fallback Cascade Against HTTP 429 Token Rate Limits"
demand_score: 9.5
status: verified
persona:
  role: "AI Platform Infrastructure Lead / Site Reliability Engineer"
  context: "High-throughput production LLM gateways serving enterprise traffic across OpenAI, Anthropic, and Gemini with spiky demand"
story:
  as_a: "Platform engineer maintaining AI infrastructure"
  i_want: "The AI gateway to proactively manage rate limits using a local Token Bucket algorithm and trigger seamless failover to secondary providers"
  so_that: "Sudden concurrency surges or exhausted Tokens Per Minute (TPM) quotas never surface HTTP 500 error screens to end users"
acceptance_criteria:
  - scenario: "Primary model provider returns HTTP 429 Too Many Requests"
    given: "The default routing policy directs requests to Claude 3.5 Sonnet"
    when: 'The upstream provider returns HTTP 429 with header "Retry-After: 15"'
    then: "The gateway must dynamically reroute the request to an equivalent fallback model (e.g. Gemini 1.5 Pro) in under 100ms without exposing upstream failure to the client"
  - scenario: "Pre-flight token estimation exceeds local safety threshold"
    given: "The current rolling TPM usage tracking window"
    when: "A large prompt is submitted and estimated token consumption exceeds 90% of provider quota"
    then: "The gateway must preemptively divert the prompt to an alternative provider pool before hitting the upstream rate limit"
edge_cases:
  - "Translating provider-specific dialect nuances (e.g. tool call formatting, function calling schemas, JSON mode across OpenAI vs Anthropic)."
  - "Preserving Server-Sent Events (SSE) streaming connections across failovers."
  - "Emitting real-time Prometheus / Datadog telemetry alerts when fallback routes are activated."
evidence:
  - source: "https://news.ycombinator.com/item?id=38419201"
    type: "hackernews"
    quote: "During OpenAI dev day outages, our entire platform died because we didn't have multi-provider fallbacks. A single 429 cascaded into 20,000 failed user sessions."
    date: "2023-11-28"
  - source: "https://github.com/BerriAI/litellm/issues/1402"
    type: "github_issue"
    quote: "Without local pre-emptive token tracking, sudden concurrent spikes trigger burst 429 limits that lock out the API key for minutes."
    date: "2024-04-09"
evaluation_rubric:
  - "Does the architecture implement automatic, transparent model switching on 429/503 upstream responses?"
  - "Is there local client-side token bucket buffering to prevent sudden burst overruns?"
  - "Are system prompts and tool schemas dynamically normalized across target providers?"
tags:
  - llm
  - gateway
  - rate-limiting
  - fallback
  - resiliency
  - ai-infra
---

# Operational Reliability Context

Coupling production systems to a single model provider without dynamic fallback routing is an unacceptable availability risk. Token limits (TPM) and request limits (RPM) are enforced abruptly during global traffic surges. A resilient AI gateway estimates prompt tokens locally and maintains hot standby routes to ensure zero user-visible disruption.
