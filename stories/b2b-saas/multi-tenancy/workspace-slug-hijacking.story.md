---
id: OS-SaaS-001
industry: b2b-saas
domain: multi-tenancy
title: "Prevenção de sequestro de subdomínio e colisão de slugs reservados em multi-tenancy"
demand_score: 9.3
status: verified
persona:
  role: "Security Architect / SaaS Product Manager"
  context: "Aplicações multi-tenant com subdomínios ou URLs com slug de organização"
story:
  as_a: "Administrador da plataforma SaaS"
  i_want: "Que o registro de novos workspaces valide estritamente listas de palavras reservadas, termos de infraestrutura e homóglifos"
  so_that: "Usuários mal-intencionados não registrem slugs como admin, api, login, billing ou support, sequestrando cookies de sessão ou aplicando phishing"
acceptance_criteria:
  - scenario: "Tentativa de registrar slug de infraestrutura reservada"
    given: "Um usuário criando uma nova conta de workspace"
    when: 'O usuário submeter o slug "admin" ou "billing" ou "auth"'
    then: "O sistema deve recusar a criação com erro 400 informando que o identificador é reservado pelo sistema"
  - scenario: "Tentativa de spoofing com caracteres Unicode homóglifos"
    given: 'O workspace legítimo "paypal" já registrado'
    when: 'Um atacante tentar registrar "pаypal" usando caractere cirílico visualmente idêntico'
    then: 'O sistema deve normalizar via punycode/NFKC ou rejeitar caracteres não-alfanuméricos ASCII estritos (regex ^[a-z0-9-]+$)'
edge_cases:
  - "Renomeação de slug existente (o slug antigo deve entrar em período de quarentena de 90 dias antes de ser liberado para terceiros)."
  - "Colisões com subdomínios técnicos internos de DNS (ex: ns1, mail, autodiscover, ftp)."
  - "Escopo de cookies definidos em .saas.com que poderiam vazar para subdomínios de usuários sem isolamento de sufixo público."
evidence:
  - source: "https://hackerone.com/reports/409850"
    type: "security_bounty_report"
    quote: "Attacker claimed the subdomain 'support.target.com' by registering the workspace slug 'support', allowing credential harvesting through cross-subdomain cookies."
    date: "2023-11-20"
  - source: "https://news.ycombinator.com/item?id=34281902"
    type: "hackernews"
    quote: "If your SaaS permits custom subdomain creation without an exhaustive denylist of over 500 reserved keywords, you have a critical security flaw from day one."
    date: "2024-04-10"
evaluation_rubric:
  - "O sistema implementa uma denylist rígida e testada de palavras reservadas do sistema?"
  - "Existe validação estrita de charset para impedir ataques de homóglifos Unicode?"
  - "Slugs excluídos ou alterados possuem quarentena temporal para evitar takeover imediato?"
tags:
  - security
  - multi-tenancy
  - saas
  - domain-takeover
  - isolation
---

# Contexto de Segurança

Em arquiteturas SaaS multi-tenant baseadas em subdomínios (`{tenant}.empresa.com`), uma das vulnerabilidades mais exploradas por bug bounty hunters é o registro de nomes internos que o próprio sistema esqueceu de reservar. Sem uma lista rígida de termos de sistema, atacantes registram subdomínios como `api`, `auth` ou `status` para roubar tokens transmitidos para o domínio pai ou enganar funcionários internos.
