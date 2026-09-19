---
id: OS-ECOM-0099
locale: en
industry: ecommerce
domain: returns-rma
title: Multi-stack coupon code validation preventing promotional stacking abuse on
  returns-rma under Zero-Trust Authentication & Token Invalidation
demand_score: 8.809999999999999
status: verified
persona:
  role: Principal E-Commerce Architect / Retail Platform Lead
  context: High-throughput digital commerce platforms with omnichannel inventory,
    warehouse logistics, and flash sales
story:
  as_a: Principal E-Commerce Architect / Retail Platform Lead
  i_want: strict rules engine governing coupon combinability and minimum order values
    on returns-rma with resilience to Zero-Trust Authentication & Token Invalidation
  so_that: customers cannot stack multiple promo codes to purchase items below manufacturing
    cost
acceptance_criteria:
- scenario: Customer combining percentage discount with dollar-off voucher on returns-rma
    combined with Zero-Trust Authentication & Token Invalidation
  given: A promo code granting 20%% off sitewide
  when: The user applies an additional $50 welcome voucher%!(EXTRA string=returns-rma)
  then: The promotions engine must enforce exclusion rules and reject stacking unless
    explicitly configured
edge_cases:
- Customers creating multiple throwaway accounts with the same physical delivery address
  to bypass limits%!(EXTRA string=returns-rma) exacerbated by Zero-Trust Authentication
  & Token Invalidation
- Cascading failover during Zero-Trust Authentication & Token Invalidation
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
  exacerbated by zero-trust authentication & token invalidation without manual intervention?
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

