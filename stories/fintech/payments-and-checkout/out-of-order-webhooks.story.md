---
id: OS-FIN-001
industry: fintech
domain: payments-and-checkout
title: "Resiliência a entrega assíncrona e fora de ordem de webhooks de pagamento"
demand_score: 9.7
status: verified
persona:
  role: "Tech Lead / Engenheiro de Pagamentos"
  context: "E-commerce e plataformas de assinatura integrando múltiplos gateways de pagamento (Stripe, Mercado Pago, Adyen)"
story:
  as_a: "Engenheiro de pagamentos"
  i_want: "Que o processador de webhooks do sistema processe eventos com controle de máquina de estados finita e versionamento temporal (event timestamp/version)"
  so_that: "Um webhook que chegue fora de ordem (ex: charge.refunded antes de charge.succeeded, ou retry de charge.failed após aprovação) não reverta indevidamente o status do pedido nem libere acesso indevido"
acceptance_criteria:
  - scenario: "Evento de sucesso chega antes da criação local do pedido"
    given: 'Um webhook de "payment.approved" recebido antes do endpoint síncrono de checkout salvar a transação'
    when: "O webhook for processado"
    then: "O sistema deve enfileirar o evento com backoff exponencial ou realizar um upsert seguro sem lançar erro 500 no gateway"
  - scenario: "Webhook com timestamp anterior a uma transição mais recente"
    given: 'O pedido já está com status "DELIVERED" ou "APPROVED"'
    when: 'Um webhook antigo de retry de "payment.processing" ou "payment.pending" for entregue'
    then: "O processador deve ignorar a transição de estado, registrar no log de auditoria como evento ignorado por obsolescência e responder HTTP 200 ao gateway"
edge_cases:
  - "Gateway que não fornece sequência monotônica ou timestamp com resolução de milissegundos."
  - "Concorrência de múltiplos webhooks do mesmo checkout disparados em paralelo pela adquirente."
  - "Retentativas do gateway disparadas por timeout temporário na aplicação."
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
  - "O código trata estados terminais como imutáveis contra eventos obsoletos?"
  - "O webhook handler responde HTTP 200 mesmo em eventos ignorados para evitar tempestade de retries do gateway?"
  - "Existe lock transacional por ID de transação ou pedido no momento do processamento do webhook?"
tags:
  - fintech
  - webhooks
  - payments
  - state-machine
  - concurrency
---

# Contexto Técnico

A premissa de que webhooks chegam em ordem cronológica de emissão é a causa número 1 de bugs silenciosos de faturamento em SaaS e e-commerce. Como gateways operam em clusters distribuídos com filas concorrentes (Kafka/SQS) e retries exponenciais independentes por falha de rede, um evento disparado 5 segundos antes pode ser entregue 30 minutos depois de um evento posterior.
