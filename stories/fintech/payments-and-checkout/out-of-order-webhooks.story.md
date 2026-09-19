---
id: OS-FIN-001
locale: en
industry: fintech
domain: payments-and-checkout
title: "Resilience to Asynchronous and Out-of-Order Payment Webhook Delivery"
demand_score: 9.7
status: verified
persona:
  role: "Fintech Lead / Staff Payment Engineer"
  context: "High-volume checkout and subscription engines consuming asynchronous webhooks from Stripe, Adyen, and PayPal across distributed queues"
story:
  as_a: "Payment platform engineer"
  i_want: "The webhook processor to enforce finite state machine (FSM) transitions with temporal event versioning"
  so_that: "An out-of-order webhook delivery (e.g. charge.refunded arriving prior to charge.succeeded, or a delayed charge.failed retry) never reverts completed orders or unlocks unauthorized access"
acceptance_criteria:
  - scenario: "Payment success webhook arrives before local order creation commits"
    given: 'A payment.approved webhook received before the synchronous checkout endpoint completes database writes'
    when: "The webhook is processed by the ingestion worker"
    then: "The system must enqueue the event with exponential backoff retry or perform an atomic upsert without returning HTTP 500 to the payment gateway"
  - scenario: "Delayed webhook with timestamp older than current order state"
    given: 'The order is already marked with a terminal state "PAID" or "DELIVERED"'
    when: 'A delayed retry of an older "payment.pending" or "payment.failed" webhook arrives hours later'
    then: "The state machine must discard the invalid transition, log an obsolescence audit trail event, and return HTTP 200 to acknowledge gateway delivery"
edge_cases:
  - "Gateways that omit monotonic sequence numbers or millisecond-precision event timestamps."
  - "Concurrent webhook deliveries for the same payment intent fired in parallel by redundant gateway regions."
  - "Locking per order ID (distributed redis lock or postgres row lock) during webhook evaluation."
evidence:
  - source: "https://stripe.com/docs/webhooks#delivery-order"
    type: "official_doc_postmortem"
    quote: "Stripe does not guarantee delivery of webhooks in the order in which they are generated. Your webhook handling endpoint must not assume events arrive sequentially."
    date: "2024-05-15"
  - source: "https://reddit.com/r/rails/comments/16k1wz8"
    type: "reddit_thread"
    quote: "Had customers receiving free subscription access because the initial failed attempt webhook was delayed 4 hours and overwrote the subsequent successful payment state."
    date: "2024-09-17"
evaluation_rubric:
  - "Does the webhook processor treat terminal order states as immutable against older event timestamps?"
  - "Does the endpoint always return HTTP 200 on handled obsolete events to halt gateway retry storms?"
  - "Is there transactional row locking per order ID during event processing?"
tags:
  - fintech
  - webhooks
  - payments
  - state-machine
  - concurrency
---

# Technical Incident Background

Assuming that webhooks arrive in exact chronological order is the primary source of billing discrepancies in subscription businesses. Because modern payment gateways dispatch events via distributed partitioning clusters (Kafka/SQS) with independent retry policies, a webhook generated 5 seconds earlier can easily land 30 minutes after subsequent events.
