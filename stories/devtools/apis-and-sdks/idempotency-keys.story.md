---
id: OS-DEV-001
locale: en
industry: devtools
domain: apis-and-sdks
title: "Mandatory Idempotency Keys on Mutation Endpoints with Automatic Retries"
demand_score: 9.6
status: verified
persona:
  role: "Backend Integration Engineer / API Platform Lead"
  context: "Mission-critical microservices and payment gateways handling distributed transactions and webhooks"
story:
  as_a: "Backend engineer integrating third-party APIs"
  i_want: "Mutation endpoints (POST/PATCH) to support an Idempotency-Key header with atomic deduplication and response caching"
  so_that: "Transient network drops, client timeouts, or connection resets never duplicate state mutations or charge customers twice"
acceptance_criteria:
  - scenario: "Replaying an identical request with the same idempotency key within 24h"
    given: 'A previous request with Idempotency-Key "idemp_9921_x" completed with HTTP 201 and payload {"id": "tx_123"}'
    when: 'A new request arrives with the exact same key "idemp_9921_x" and identical payload'
    then: 'The API must return HTTP 201 with the cached payload and header "Idempotent-Replayed: true", without re-executing business logic'
  - scenario: "Payload mismatch conflict on identical key"
    given: 'The key "idemp_9921_x" was already processed with payload A'
    when: 'A client submits the same key "idemp_9921_x" with a different payload B'
    then: 'The API must reject with HTTP 422 Unprocessable Entity and error code "idempotency_key_payload_mismatch"'
  - scenario: "Concurrent in-flight requests sharing the same key"
    given: 'A request with key "idemp_9921_x" is currently being processed by a worker'
    when: "A second concurrent request arrives with the exact same key before completion"
    then: "The API must return HTTP 409 Conflict or hold a distributed lock until the first completes, never executing duplicate parallel routines"
edge_cases:
  - "Idempotency key TTL expiration after the standard retention window (24 hours)."
  - "Sub-millisecond race conditions mitigated via distributed locks (e.g. Redis Redlock or Postgres advisory locks)."
  - "Full header and status code preservation (replaying original response headers, not just raw JSON body)."
evidence:
  - source: "https://github.com/stripe/stripe-go/issues/632"
    type: "github_issue"
    quote: "During network drops on cloud providers, client retries without idempotency keys produced duplicate customer charges totaling over $45k in less than 20 minutes."
    date: "2024-09-14"
  - source: "https://news.ycombinator.com/item?id=38192801"
    type: "hackernews"
    quote: "Building distributed systems without first-class idempotency in API design is tech debt that will eventually wake someone up at 3 AM with data corruption."
    date: "2025-01-18"
evaluation_rubric:
  - "Does the spec or implementation reject reused keys with altered payloads with HTTP 422?"
  - "Is there a distributed lock mechanism for in-flight concurrent requests sharing the same key?"
  - "Does the replayed response explicitly notify the client via an Idempotent-Replayed header?"
tags:
  - api-design
  - idempotency
  - distributed-systems
  - reliability
---

# Architectural Context

The absence of native idempotency in APIs is the leading cause of silent data corruption in distributed microservices. When a 30-second client-side HTTP timeout fires, the server may have already committed the database transaction. Without first-class `Idempotency-Key` deduplication, automated client retries inevitably spawn duplicate resources.
