---
id: OS-FIN-002
locale: en
industry: fintech
domain: reconciliation
title: "Proportional Interchange Fee (MDR) and Fixed Gateway Fee Accounting on Partial Refunds"
demand_score: 9.1
status: verified
persona:
  role: "Financial Ledger Engineer / ERP Integration Specialist"
  context: "High-volume e-commerce platforms issuing daily partial returns with card acquiring networks"
story:
  as_a: "Financial controller or ledger engineer"
  i_want: "The reconciliation engine to proportionately calculate merchant discount rates (MDR) on partial refunds while properly attributing non-refundable fixed fees"
  so_that: "Internal accounting ledgers match acquiring settlement bank statements to the penny without accumulating phantom discrepancies"
acceptance_criteria:
  - scenario: "Partial refund with percentage fee (3% MDR)"
    given: "An original sale of $100.00 with a 3% MDR acquiring deduction ($3.00)"
    when: "A customer initiates a partial return for 1 item valued at $40.00"
    then: "The ledger must credit back $1.20 in proportional MDR (3% of $40) and debit $38.80 from merchant payout balances"
  - scenario: "Non-refundable fixed gateway authorization fee"
    given: "The payment processor charges a fixed $0.30 transaction fee that is non-refundable under all circumstances"
    when: "A full or partial refund is processed"
    then: "The $0.30 fixed fee must be recorded as an irrecoverable operational expense and never credited back against the refund principal"
edge_cases:
  - "Half-even rounding (Banker's Rounding) on sub-cent fee calculations to avoid systematic financial drift over millions of transactions."
  - "Multi-currency foreign exchange (FX) spread fluctuations between purchase authorization and refund execution."
  - "Refunds processed across closed fiscal calendar months requiring adjusting ledger entries."
evidence:
  - source: "https://news.ycombinator.com/item?id=37728192"
    type: "hackernews"
    quote: "Our billing engine accumulated an $18,000 reconciliation discrepancy over 6 months simply because it assumed acquiring fees were zeroed out on partial refunds."
    date: "2023-10-02"
  - source: "https://reddit.com/r/accounting/comments/16u810d"
    type: "reddit_thread"
    quote: "Every junior developer implements refund as 'total - refund_amount' without accounting for whether the payment processor returns the interchange fee."
    date: "2024-01-14"
evaluation_rubric:
  - "Does the financial calculation utilize Banker's Rounding (Half-Even) to prevent statistical drift?"
  - "Does the model explicitly distinguish variable percentage fees from fixed non-refundable transaction fees?"
  - "Are refund operations persisted as immutable double-entry ledger journal entries?"
tags:
  - fintech
  - reconciliation
  - accounting
  - refunds
  - payments
---

# Financial Accounting Hazards

Credit card acquiring reconciliation is full of subtle traps. Most payment processors (Stripe, Adyen, Stone, Cielo) rebate the variable percentage component (interchange/MDR) on partial refunds, but retain the fixed per-transaction fee (e.g. $0.30). Naive systems that simply subtract the gross refund amount accumulate serious audit discrepancies that trigger external investigation.
