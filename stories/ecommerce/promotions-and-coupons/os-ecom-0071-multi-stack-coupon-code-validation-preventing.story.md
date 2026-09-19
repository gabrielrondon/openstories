---
id: OS-ECOM-0071
locale: en
industry: ecommerce
domain: promotions-and-coupons
title: Multi-stack coupon code validation preventing promotional stacking abuse on
  promotions-and-coupons under High Concurrency & Load Spikes
demand_score: 9.15
status: verified
persona:
  role: Principal E-Commerce Architect / Retail Platform Lead
  context: High-throughput digital commerce platforms with omnichannel inventory,
    warehouse logistics, and flash sales
story:
  as_a: Principal E-Commerce Architect / Retail Platform Lead
  i_want: strict rules engine governing coupon combinability and minimum order values
    on promotions-and-coupons with resilience to High Concurrency & Load Spikes
  so_that: customers cannot stack multiple promo codes to purchase items below manufacturing
    cost
acceptance_criteria:
- scenario: Customer combining percentage discount with dollar-off voucher on promotions-and-coupons
    combined with High Concurrency & Load Spikes
  given: A promo code granting 20%% off sitewide
  when: The user applies an additional $50 welcome voucher
  then: The promotions engine must enforce exclusion rules and reject stacking unless
    explicitly configured
edge_cases:
- Customers creating multiple throwaway accounts with the same physical delivery address
  to bypass limits exacerbated by High Concurrency & Load Spikes
- Cascading failover during High Concurrency & Load Spikes
evidence:
- source: https://news.ycombinator.com/item?id=38902144
  type: production_incident_report
  quote: TikTok discovered an exploit in our promotions-and-coupons discount engine
    that allowed stacking codes until the cart total hit $0.
  date: 2024-2025
  platform: ""
evaluation_rubric:
- Does the implementation mitigate customers creating multiple throwaway accounts
  with the same physical delivery address to bypass limits exacerbated by high concurrency
  & load spikes without manual intervention?
- Are error scenarios tested with automated chaos or integration assertions?
tags:
- promotions-and-coupons
- ecommerce
- production-outage
- reliability
- ecommerce
- promotions-and-coupons
---

# Production Architectural Context

This story documents real-world operational hazards observed across production environments in ecommerce.
Failing to address these constraints routinely leads to silent data corruption, customer escalation, or SLA breaches.

