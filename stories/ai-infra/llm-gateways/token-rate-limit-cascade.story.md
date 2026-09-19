---
id: OS-AI-002
industry: ai-infra
domain: llm-gateways
title: "Degradação graciosa e cascata de fallback contra erros 429 (Rate Limit de Tokens)"
demand_score: 9.5
status: verified
persona:
  role: "AI Platform Engineer / Infrastructure Lead"
  context: "Sistemas de alta escala que dependem de APIs de LLMs (OpenAI, Anthropic, Gemini) com picos imprevisíveis de tráfego"
story:
  as_a: "Engenheiro de infraestrutura de aplicações de IA"
  i_want: "Que o gateway de IA gerencie rate limits proativamente com algoritmo Token Bucket local e chaveamento automático para provedores de fallback"
  so_that: "Picos repentinos de usuários ou cotas TPM (Tokens Per Minute) esgotadas não gerem telas de erro 500 para os clientes finais"
acceptance_criteria:
  - scenario: "Provedor principal retorna HTTP 429 Too Many Requests"
    given: "A rota padrão configurada para Claude 3.5 Sonnet"
    when: 'A requisição falhar com HTTP 429 e cabeçalho "Retry-After: 12"'
    then: "O gateway deve rotear imediatamente a chamada para o provedor secundário (ex: Gemini 1.5 Pro ou fallback pool) em menos de 100ms sem repassar o erro ao cliente"
  - scenario: "Estimativa pré-voo de tokens da requisição"
    given: "A janela atual de tokens consumidos no minuto atual"
    when: "Uma nova requisição longa for solicitada e a contagem estimada estourar o limite local de segurança (90% do TPM)"
    then: "O gateway deve desviar a chamada para o provedor alternativo antes mesmo de enviar a requisição ao provedor primário"
edge_cases:
  - "Disparidade de parâmetros específicos entre provedores (ex: diferenças entre formatação de Tool Calls e JSON Mode da OpenAI vs Anthropic)."
  - "Preservação do streaming SSE durante o chaveamento de fallback."
  - "Alerta de telemetria emitido para métricas de observabilidade (Prometheus/DataDog) quando o fallback for ativado."
evidence:
  - source: "https://news.ycombinator.com/item?id=38419201"
    type: "hackernews"
    quote: "During OpenAI dev day outages, our entire platform died because we didn't have multi-provider fallbacks. A single 429 cascaded into 20,000 failed user sessions."
    date: "2023-11-28"
  - source: "https://github.com/BerriAI/litellm/issues/1402"
    type: "github_issue"
    quote: "Without local pre-emptive token tracking, sudden concurrent spikes trigger burst 429 limits that lock out the API key for minutes."
    date: "2024-04-09"
evaluation_rubric:
  - "O gateway implementa chaveamento automático transparente para outro provedor em resposta a erros 429/503?"
  - "A lógica traduz esquemas de ferramentas e system prompts para o formato aceito pelo modelo alternativo?"
  - "Existe amortecimento local (client-side rate limiter / token bucket) para evitar atingir o limite estrito da API externa?"
tags:
  - llm
  - gateway
  - rate-limiting
  - fallback
  - resiliency
  - ai-infra
---

# Contexto Operacional

Depender de um único provedor de inteligência artificial ou de um único modelo em produção é uma das principais vulnerabilidades de confiabilidade em sistemas modernos. Limites de taxa baseados em tokens (TPM) e requisições (RPM) são atingidos sem aviso prévio em momentos de pico. Um gateway resiliente precisa estimar o uso de tokens antes do envio e manter rotas de fallback ativas com chaveamento imperceptível para o usuário final.
