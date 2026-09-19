---
id: OS-DEV-002
industry: devtools
domain: auth-and-iam
title: Prevenção de invalidação de Magic Link por scanners corporativos de e-mail (Pre-fetching)
demand_score: 9.4
status: verified
persona:
  role: Engenheiro de Segurança / Full Stack Lead
  context: Aplicações B2B SaaS onde os clientes utilizam suites empresariais (Microsoft Defender, Proofpoint, Mimecast)
story:
  as_a: Usuário corporativo acessando a plataforma via Magic Link por e-mail
  i_want: Que o link do e-mail não seja invalidado antes do meu clique real
  so_that: Eu consiga fazer login sem receber a mensagem de erro "Token expirado ou já utilizado" causada pelo bot antivírus da empresa
acceptance_criteria:
  - scenario: Varredura automatizada (GET request) por bot de antivírus corporativo
    given: Um token de login temporário de uso único gerado no link enviado por e-mail
    when: O bot de segurança da empresa realiza uma requisição GET automática no link para inspecionar malware
    then: O sistema deve renderizar uma página intermediária com botão "Confirmar Login" (ou desafio JS leve) SEM consumir nem queimar o token no GET
  - scenario: Clique humano deliberado no botão de confirmação (POST request)
    given: O usuário acessou a página de confirmação
    when: O usuário clica no botão "Entrar na minha conta" enviando um POST
    then: O token é validado atomicamente, a sessão do usuário é iniciada e o token é marcado como expirado
edge_cases:
  - Scanners com headless browsers automatizados que simulam cliques ou executam JavaScript.
  - Expiração estrita do token em no máximo 15 minutos.
  - Associação de IP/User-Agent ou fingerprint leve entre a solicitação inicial e a confirmação.
evidence:
  - source: "https://github.com/supabase/gotrue/issues/1204"
    type: "github_issue"
    quote: "Enterprise customers using Outlook SafeLinks / Microsoft Defender have all magic links pre-fetched. When users click it 5 seconds later, it says link is expired."
    date: "2024-03-22"
  - source: "https://reddit.com/r/webdev/comments/182k9p1"
    type: "reddit_thread"
    quote: "Never, ever consume a magic link token on a GET request. Mail security gateways fetch every link automatically and burn the token instantly."
    date: "2024-11-05"
evaluation_rubric:
  - "O endpoint que recebe o clique do link por GET evita queimar o token imediatamente?"
  - "Existe uma ação explícita (POST ou clique) na landing page para consolidar a sessão?"
  - "Há rate limiting e proteção contra brute-force nos tokens de magic link?"
tags:
  - auth
  - magic-link
  - security
  - enterprise
---

# Contexto de Falha em Produção

Dezenas de startups adotam Magic Links como mecanismo padrão de login sem senha ("passwordless") e descobrem no primeiro cliente Enterprise que 100% dos usuários não conseguem logar. Serviços como Microsoft Defender for Office 365, Proofpoint URL Defense e Mimecast abrem programaticamente todo link que chega na caixa de entrada para checar vírus e phishing. Se o endpoint `/auth/callback?token=XYZ` consome o token na requisição GET, o usuário humano recebe erro ao clicar.
