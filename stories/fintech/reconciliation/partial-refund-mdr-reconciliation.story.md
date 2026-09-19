---
id: OS-FIN-002
industry: fintech
domain: reconciliation
title: Conciliação proporcional de taxas de MDR e tarifas fixas em estornos parciais
demand_score: 9.1
status: verified
persona:
  role: Analista Financeiro / Desenvolvedor de ERP Contábil
  context: Empresas transacionando alto volume de pagamentos com estornos frequentes (devoluções parciais de itens)
story:
  as_a: Analista financeiro de um e-commerce
  i_want: Que o motor de conciliação do sistema deduza proporcionalmente o percentual da taxa de transação (MDR) nos estornos parciais e contabilize corretamente a tarifa fixa não estornada pela adquirente
  so_that: O saldo contábil no ERP coincida exatamente com o extrato bancário de liquidação da adquirente sem divergências de centavos
acceptance_criteria:
  - scenario: Estorno parcial com taxa percentual (MDR de 3%)
    given: Uma venda de R$ 100,00 com MDR de 3% (R$ 3,00 retidos)
    when: O cliente solicitar a devolução de 1 item no valor de R$ 40,00
    then: O estorno contábil deve creditar de volta R$ 1,20 de MDR proporcional (3% de R$ 40) e debitar R$ 38,80 do lojista
  - scenario: Taxa fixa de gateway não estornável
    given: A operadora cobra R$ 0,50 fixos por transação que não são devolvidos em nenhum caso de cancelamento
    when: Um estorno (total ou parcial) for liquidado
    then: A tarifa de R$ 0,50 deve ser classificada como despesa operacional irrecuperável e nunca abatida como estorno
edge_cases:
  - Arredondamento monetário de meio centavo (utilizar algoritmo Half-Even / Banker's Rounding para evitar acúmulo de viés contábil).
  - Estornos que ocorrem em moedas diferentes da moeda de liquidação (flutuação de spread cambial FX).
  - Estornos solicitados após o fechamento da competência fiscal do mês.
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
  - "O cálculo contábil utiliza Banker's Rounding (Half-Even) para lidar com frações de centavos?"
  - "A lógica diferencia componentes percentuais de taxas (MDR) de taxas fixas não-estornáveis por transação?"
  - "O modelo de dados armazena o evento de estorno como registro contábil imutável em partidas dobradas?"
tags:
  - fintech
  - reconciliation
  - accounting
  - refunds
  - payments
---

# Contexto Contábil e Financeiro

A conciliação bancária de pagamentos com cartão de crédito é repleta de armadilhas. A maioria das adquirentes (Stone, Cielo, Stripe, Adyen) devolve a taxa percentual (MDR) proporcional ao valor do estorno, mas retém a tarifa fixa de autorização (ex: R$ 0,39 ou $0.30). Sistemas que tratam o estorno apenas subtraindo o valor bruto geram um furo contábil cumulativo que estoura em auditorias externas.
