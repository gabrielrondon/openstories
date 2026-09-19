---
id: OS-ECOM-0095
locale: en
industry: ecommerce
domain: returns-rma
title: Multi-stack coupon code validation preventing promotional stacking abuse on
  returns-rma under Multi-Tenant Data Leakage & Isolation Breaches
demand_score: 9.030000000000001
status: verified
persona:
  role: Principal E-Commerce Architect / Retail Platform Lead
  context: High-throughput digital commerce platforms with omnichannel inventory,
    warehouse logistics, and flash sales
story:
  as_a: Principal E-Commerce Architect / Retail Platform Lead
  i_want: strict rules engine governing coupon combinability and minimum order values
    on returns-rma with resilience to Multi-Tenant Data Leakage & Isolation Breaches
  so_that: customers cannot stack multiple promo codes to purchase items below manufacturing
    cost
acceptance_criteria:
- scenario: Customer combining percentage discount with dollar-off voucher on returns-rma
    combined with Multi-Tenant Data Leakage & Isolation Breaches
  given: A promo code granting 20%% off sitewide
  when: The user applies an additional $50 welcome voucher%!(EXTRA string=returns-rma)
  then: The promotions engine must enforce exclusion rules and reject stacking unless
    explicitly configured
edge_cases:
- Customers creating multiple throwaway accounts with the same physical delivery address
  to bypass limits%!(EXTRA string=returns-rma) exacerbated by Multi-Tenant Data Leakage
  & Isolation Breaches
- Cascading failover during Multi-Tenant Data Leakage & Isolation Breaches
evidence:
- source: https://news.ycombinator.com/item?id=38902144
  type: production_incident_report
  quote: TikTok discovered an exploit in our returns-rma discount engine that allowed
    stacking codes until the cart total hit $0.
  date: 2024-2025
  platform: ""
evaluation_rubric:
- Does the implementation mitigate customers creating multiple throwaway accounts
  with the same physical delivery address to bypass limits%!(extra string=returns-rma)
  exacerbated by multi-tenant data leakage & isolation breaches without manual intervention?
- Are error scenarios tested with automated chaos or integration assertions?
tags:
- returns-rma
- ecommerce
- production-outage
- reliability
- ecommerce
- returns-rma
---

# Production Architectural Context

This story documents real-world operational hazards observed across production environments in ecommerce.
Failing to address these constraints routinely leads to silent data corruption, customer escalation, or SLA breaches.

