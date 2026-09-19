---
id: OS-SAAS-0045
locale: en
industry: b2b-saas
domain: audit-logs
title: Granular Role-Based Access Control (RBAC) with attribute constraints for audit-logs
  under Multi-Tenant Data Leakage & Isolation Breaches
demand_score: 9.030000000000001
status: verified
persona:
  role: Enterprise SaaS Architect / Security Lead
  context: Multi-tenant B2B enterprise architectures serving Fortune 500 customers
    with rigorous compliance requirements
story:
  as_a: Enterprise SaaS Architect / Security Lead
  i_want: fine-grained policy evaluation (ABAC/RBAC) on all audit-logs endpoints with
    resilience to Multi-Tenant Data Leakage & Isolation Breaches
  so_that: tenant members with restricted roles cannot access privileged audit records
    or billing settings
acceptance_criteria:
- scenario: Horizontal privilege escalation attempt via direct ID reference on audit-logs
    combined with Multi-Tenant Data Leakage & Isolation Breaches
  given: A user logged into tenant organization Alpha
  when: The user queries resource ID belonging to tenant Beta
  then: The authorization layer must return HTTP 404 Not Found rather than 403 Forbidden
    to prevent resource ID enumeration
edge_cases:
- Users belonging to multiple organizations switching active workspace context concurrently
  exacerbated by Multi-Tenant Data Leakage & Isolation Breaches
- Cascading failover during Multi-Tenant Data Leakage & Isolation Breaches
evidence:
- source: https://hackerone.com/reports/512091
  type: production_incident_report
  quote: A classic IDOR vulnerability in our audit-logs module allowed team members
    to inspect invoice PDFs from rival companies.
  date: 2024-2025
  platform: ""
evaluation_rubric:
- Does the implementation mitigate users belonging to multiple organizations switching
  active workspace context concurrently exacerbated by multi-tenant data leakage &
  isolation breaches without manual intervention?
- Are error scenarios tested with automated chaos or integration assertions?
tags:
- audit-logs
- b2b-saas
- production-outage
- reliability
- b2b-saas
- audit-logs
---

# Production Architectural Context

This story documents real-world operational hazards observed across production environments in b2b-saas.
Failing to address these constraints routinely leads to silent data corruption, customer escalation, or SLA breaches.

