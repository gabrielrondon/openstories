---
id: OS-DEV-001
industry: devtools
domain: apis-and-sdks
title: "Idempotência obrigatória em endpoints de mutação com retries automáticos"
demand_score: 9.6
status: verified
persona:
  role: "Engenheiro de Integração de APIs / Backend Lead"
  context: "Microsserviços e sistemas de missão crítica lidando com transações e webhooks"
story:
  as_a: "Engenheiro integrador de APIs de terceiros"
  i_want: "Que todos os endpoints de criação e mutação (POST/PATCH) suportem o cabeçalho Idempotency-Key com deduplicação atômica e cache da resposta"
  so_that: "Falhas transitórias de rede, timeouts ou reconexões de cliente nunca dupliquem registros no banco nem cobrem o cliente duas vezes"
acceptance_criteria:
  - scenario: "Replay de requisição idêntica com mesma chave dentro de 24h"
    given: 'Uma requisição anterior com Idempotency-Key "idemp_9921_x" concluiu com status 201'
    when: 'Uma nova requisição for enviada com a mesma chave "idemp_9921_x" e mesmo payload'
    then: 'A API deve retornar HTTP 201 com o payload original, adicionando o cabeçalho "Idempotent-Replayed: true", sem reexecutar a lógica'
  - scenario: "Conflito de payload com a mesma chave"
    given: 'A chave "idemp_9921_x" foi usada com o payload A'
    when: 'O cliente enviar a mesma chave "idemp_9921_x" mas com payload B'
    then: 'A API deve rejeitar com HTTP 422 Unprocessable Entity e código "idempotency_key_payload_mismatch"'
  - scenario: "Requisição concorrente em trânsito com a mesma chave"
    given: 'Uma requisição com a chave "idemp_9921_x" ainda está sendo processada no backend'
    when: "Uma segunda requisição chegar simultaneamente com a mesma chave"
    then: "A API deve retornar HTTP 409 Conflict ou segurar o lock distribuído até a primeira concluir, nunca executando duas threads em paralelo"
edge_cases:
  - "Expiração da chave de idempotência após a janela de retenção (24 horas)."
  - "Concorrência real de milissegundos tratada via distributed lock (ex: Redis Redlock ou Postgres advisory lock)."
  - "Salvamento do cabeçalho de resposta original e status code original, não apenas do body JSON."
evidence:
  - source: "https://github.com/stripe/stripe-go/issues/632"
    type: "github_issue"
    quote: "During network drops on cloud providers, client retries without idempotency keys produced duplicate customer charges totaling over $45k in less than 20 minutes."
    date: "2024-09-14"
  - source: "https://news.ycombinator.com/item?id=38192801"
    type: "hackernews"
    quote: "Building distributed systems without first-class idempotency in API design is tech debt that will eventually wake someone up at 3 AM with data corruption."
    date: "2025-01-18"
evaluation_rubric:
  - "O código valida se a mesma chave foi reutilizada com payload diferente e retorna 422?"
  - "Existe mecanismo de lock distribuído para requisições em trânsito com a mesma chave?"
  - "A resposta retransmitida inclui indicação explícita de replay (ex: Idempotent-Replayed: true)?"
tags:
  - api-design
  - idempotency
  - distributed-systems
  - reliability
---

# Contexto Arquitetural

A ausência de idempotência em APIs modernas é uma das principais fontes de corrupção silenciosa de dados em microsserviços. Quando um timeout de 30 segundos ocorre no client HTTP (ex: Axios, Faraday, Fetch), o servidor pode já ter completado o commit no banco de dados. Sem suporte nativo a `Idempotency-Key`, o retry automático inevitavelmente cria uma entidade duplicada.
