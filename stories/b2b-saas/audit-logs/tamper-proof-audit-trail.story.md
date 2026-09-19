---
id: OS-SaaS-002
industry: b2b-saas
domain: audit-logs
title: Trilha de auditoria imutável (Audit Trail) com encadeamento criptográfico para conformidade SOC2
demand_score: 9.2
status: verified
persona:
  role: Enterprise Compliance Officer / Staff Security Engineer
  context: Aplicações B2B SaaS vendendo para clientes Enterprise que exigem certificação SOC2 Type II e auditoria forense
story:
  as_a: Auditor de segurança e conformidade de um cliente corporativo
  i_want: Que todas as operações de mutação de privilégios, exportação de dados e acessos confidenciais gerem registros de auditoria append-only assinados com hash anterior
  so_that: Nenhum administrador de banco ou invasor com acesso interno consiga adulterar ou apagar logs retroativamente para ocultar um vazamento
acceptance_criteria:
  - scenario: Criação de novo evento de auditoria
    given: Um administrador alterando as permissões de um usuário de "Viewer" para "Admin"
    when: O evento for gravado na tabela de audit log
    then: O registro deve conter obrigatoriamente actor_id, target_id, ip_address, user_agent, diff de permissões em formato estruturado e prev_hash criptográfico (SHA-256)
  - scenario: Tentativa de update ou delete direto no banco de logs
    given: A tabela de trilha de auditoria configurada
    when: Uma instrução SQL UPDATE ou DELETE for executada diretamente na tabela
    then: Triggers de banco ou políticas de retenção WORM (Write Once, Read Many) devem abortar a operação com erro de integridade estrita
edge_cases:
  - Anonimização controlada de PII (dados pessoais sensíveis) para cumprir LGPD/GDPR sem quebrar a integridade criptográfica da cadeia de hashes.
  - Exportação em massa em formato JSON Lines (JSONL) ou CEF (Common Event Format) para integração com SIEMs (Splunk, Datadog).
  - Tolerância a alta vazão de escrita através de buffers de agregação ou particionamento mensal.
evidence:
  - source: "https://news.ycombinator.com/item?id=35890214"
    type: "hackernews"
    quote: "Failed our SOC2 audit because our audit logs were just normal mutable Postgres rows where any dev with DB access could UPDATE them without a trace."
    date: "2023-05-11"
  - source: "https://github.com/boxyhq/jackson/issues/489"
    type: "github_issue"
    quote: "Enterprise buyers require non-repudiation in audit events. If you cannot prove the log wasn't edited in PostgreSQL, procurement won't approve the deal."
    date: "2024-02-19"
evaluation_rubric:
  - "Os eventos de auditoria contêm contexto forense completo (ator, IP, tenant, timestamp UTC e payload de diff)?"
  - "Existe garantia técnica de imutabilidade (tabela append-only, permissões de banco restritas ou hash chain)?"
  - "O sistema fornece endpoint para streaming ou exportação de logs para ferramentas SIEM do cliente?"
tags:
  - security
  - soc2
  - audit-logs
  - compliance
  - enterprise
  - b2b-saas
---

# Contexto de Conformidade e Negócios

A trilha de auditoria é frequentemente o requisito eliminatório ("deal breaker") para fechar vendas com clientes corporativos (Fortune 500 ou bancos). Empresas que implementam logs como simples tabelas relacionais mutáveis falham em auditorias SOC 2 Type II e ISO 27001 por não conseguirem garantir a não-repudiação dos dados.
