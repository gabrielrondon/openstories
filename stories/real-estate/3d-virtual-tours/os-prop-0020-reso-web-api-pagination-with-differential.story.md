---
id: OS-PROP-0020
locale: en
industry: real-estate
domain: 3d-virtual-tours
title: RESO Web API pagination with differential delta synchronization for 3d-virtual-tours
  under Disaster Recovery & Cascading Failover
demand_score: 9.13
status: verified
persona:
  role: PropTech Architect / MLS Integration Engineer
  context: Real estate MLS feeds (RETS/RESO Web API), virtual walkthroughs, and escrow
    workflows
story:
  as_a: PropTech Architect / MLS Integration Engineer
  i_want: replication-state tracking and delta synchronization for 3d-virtual-tours
    property listings with resilience to Disaster Recovery & Cascading Failover
  so_that: property search portals update price drops and pending sale statuses within
    60 seconds of MLS changes
acceptance_criteria:
- scenario: MLS server drops pagination token during large 50k listing pull for 3d-virtual-tours
    combined with Disaster Recovery & Cascading Failover
  given: A background synchronization job consuming RESO API
  when: The upstream server returns HTTP 500 midway through a paginated sync
  then: The job must resume from the last committed ModificationTimestamp without
    re-pulling identical records
edge_cases:
- Listings deleted or marked private by agents leaving phantom listings active on
  public search exacerbated by Disaster Recovery & Cascading Failover
- Cascading failover during Disaster Recovery & Cascading Failover
evidence:
- source: https://github.com/reso-standards/reso-web-api/issues/102
  type: production_incident_report
  quote: Our 3d-virtual-tours sync dropped pending status updates, showing houses
    as available that had already closed escrow.
  date: 2024-2025
  platform: ""
evaluation_rubric:
- Does the implementation mitigate listings deleted or marked private by agents leaving
  phantom listings active on public search exacerbated by disaster recovery & cascading
  failover without manual intervention?
- Are error scenarios tested with automated chaos or integration assertions?
tags:
- 3d-virtual-tours
- real-estate
- production-outage
- reliability
- real-estate
- 3d-virtual-tours
---

# Production Architectural Context

This story documents real-world operational hazards observed across production environments in real-estate.
Failing to address these constraints routinely leads to silent data corruption, customer escalation, or SLA breaches.

