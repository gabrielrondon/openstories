---
id: OS-ED-0020
locale: en
industry: edtech
domain: live-proctoring
title: Resilient offline exam state synchronization and autosave for live-proctoring
  under Disaster Recovery & Cascading Failover
demand_score: 9.13
status: verified
persona:
  role: EdTech Lead Architect / Educational Platform Engineer
  context: Scalable learning management systems (LMS), online examinations, and interactive
    educational streaming
story:
  as_a: EdTech Lead Architect / Educational Platform Engineer
  i_want: local IndexedDB draft caching with differential background sync for live-proctoring
    with resilience to Disaster Recovery & Cascading Failover
  so_that: students never lose essay answers or exam progress during intermittent
    Wi-Fi disconnects
acceptance_criteria:
- scenario: Internet connection drops while student is submitting timed exam on live-proctoring
    combined with Disaster Recovery & Cascading Failover
  given: A student actively answering a 60-minute certification test
  when: The browser loses network connectivity 3 minutes before the timer expires
  then: The client must continue storing encrypted keystrokes locally and automatically
    synchronize upon reconnect
edge_cases:
- System clock tampering on student laptops to artificially extend examination time
  limits exacerbated by Disaster Recovery & Cascading Failover
- Cascading failover during Disaster Recovery & Cascading Failover
evidence:
- source: https://reddit.com/r/professors/comments/18k192a
  type: production_incident_report
  quote: Hundreds of university students lost their final exam essays when campus
    Wi-Fi dropped on live-proctoring.
  date: 2024-2025
  platform: ""
evaluation_rubric:
- Does the implementation mitigate system clock tampering on student laptops to artificially
  extend examination time limits exacerbated by disaster recovery & cascading failover
  without manual intervention?
- Are error scenarios tested with automated chaos or integration assertions?
tags:
- live-proctoring
- edtech
- production-outage
- reliability
- edtech
- live-proctoring
---

# Production Architectural Context

This story documents real-world operational hazards observed across production environments in edtech.
Failing to address these constraints routinely leads to silent data corruption, customer escalation, or SLA breaches.

