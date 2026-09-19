---
id: OS-AI-001
industry: ai-infra
domain: agentic-workflows
title: Detecção e interrupção automática de loops infinitos de tool calling com degradação graciosa
demand_score: 9.8
status: verified
persona:
  role: AI Engineer / LLM Platform Architect
  context: Construindo agentes autônomos que orquestram múltiplos tool calls em produção
story:
  as_a: Engenheiro de plataforma de agentes de IA
  i_want: Que o runtime de execução do agente detecte padrões repetitivos de tool calls (mesmos parâmetros ou oscilação entre duas ferramentas) e aborte o loop antes de esgotar o orçamento
  so_that: A aplicação não gaste centenas de dólares em tokens desnecessários e o usuário não fique esperando infinitamente por uma resposta que nunca chega
acceptance_criteria:
  - scenario: Agente invocando a mesma ferramenta com argumentos idênticos 3 vezes consecutivas
    given: Um agente autônomo em execução com orçamento configurado de 10 passos
    when: O agente invocar a ferramenta "search_api" com os mesmos parâmetros exatos pela 3ª vez consecutiva
    then: O runtime deve interceptar a chamada, injetar uma mensagem do sistema alertando sobre repetição inútil e forçar o agente a sintetizar uma resposta ou admitir que não encontrou o dado
  - scenario: Oscilação pendular entre duas ferramentas (A -> B -> A -> B)
    given: O histórico das últimas 4 ações do agente
    when: As chamadas alternarem ciclicamente entre duas ferramentas sem produzir novos dados
    then: O disjuntor (circuit breaker) deve disparar, emitindo evento de alerta de telemetria e finalizando o fluxo com status "agent_loop_detected"
edge_cases:
  - Ferramentas paginadas que legítimamente chamam a mesma tool incrementando `page` ou `offset`.
  - Mudança imperceptível de argumentos gerada pelo LLM para tentar burlar a checagem exata (requer checagem por similaridade de argumentos).
  - Preservação do histórico e motivo da interrupção para debugging do desenvolvedor.
evidence:
  - source: "https://github.com/langchain-ai/langgraph/issues/412"
    type: "github_issue"
    quote: "An agent got stuck calling the weather tool with minor variations in loop overnight, resulting in a $1,400 OpenAI bill and no final answer."
    date: "2024-08-03"
  - source: "https://news.ycombinator.com/item?id=39912044"
    type: "hackernews"
    quote: "If your LLM agent framework doesn't have a hard circuit breaker for repetitive tool loops, you are running an unmonitored money furnace in production."
    date: "2024-12-19"
evaluation_rubric:
  - "O runtime possui teto rígido configurável de passos máximos (`max_steps` ou `budget`)?"
  - "Existe detecção de repetição de ferramentas e de oscilação cíclica (A-B-A-B)?"
  - "Ao atingir o limite, o sistema entrega ao usuário uma resposta parcial explicativa em vez de crashar silenciosamente?"
tags:
  - ai-agents
  - tool-calling
  - cost-control
  - reliability
  - circuit-breaker
---

# Contexto de Riscos em Produção

O problema de "Looping Hallucination" em agentes autônomos é um dos maiores desafios de segurança financeira em AI Infra. Diferente de sistemas determinísticos, um modelo de linguagem pode entrar em um ciclo vicioso quando a resposta de uma ferramenta não satisfaz o prompt interno, tentando repetidamente formular pequenas variações da mesma pergunta até estourar a cota de cartão de crédito.
