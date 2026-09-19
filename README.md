# OpenStories 📖⚡

<div align="center">

```
   ____                   ____  _             _           
  / __ \                 / ____|| |           (_)          
 | |  | | _ __   ___  _ | (___  | |_  ___  _ __ _  ___  ___ 
 | |  | || '_ \ / _ \| '_ \___ \ | __|/ _ \| '__| |/ _ \/ __|
 | |__| || |_) |  __/| | | |___) || |_| (_) | |  | |  __/\__ \
  \____/ | .__/ \___||_| |_|_____/ \__|\___/|_|  |_|\___||___/
         | |   Open-Source Evidence-Backed User Stories
         |_|   for AI Coding Agents & Systems • v1.0.0
```

[![Go Report Card](https://goreportcard.com/badge/github.com/gabrielrondon/openstories)](https://goreportcard.com/report/github.com/gabrielrondon/openstories)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![MCP Compliant](https://img.shields.io/badge/MCP-2024--11--05-8b5cf6.svg)](https://modelcontextprotocol.io)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://golang.org)
[![Tuturama](https://img.shields.io/badge/Maintained%20by-Tuturama-indigo.svg)](https://tuturama.com)

**Aterre as features e PRs da sua IA em problemas reais da humanidade.**  
*O primeiro motor open-source de User Stories baseadas em evidências reais (GitHub Issues, post-mortems e discussões de produção) com servidor MCP nativo e Spec Stress-Tester.*

[Web Dashboard](https://openstories.tuturama.com) • [Quickstart](#-quickstart-em-30-segundos) • [Integração MCP](#-integração-com-claude-desktop--cursor) • [Spec Stress-Tester](#-o-spec-stress-tester-linter-de-realidade) • [Taxonomia](#-árvore-de-classificação)

</div>

---

## 💥 O Problema: O Viés do "Happy Path" nas IAs

Quando você pede para **Claude Code**, **Cursor**, **Antigravity** ou **Codex** arquitetar uma funcionalidade:
- Elas geram especificações e testes de brinquedo (*"happy path"*).
- Elas não têm ideia das dezenas de incidentes bizarros que acontecem em produção (ex: *scanners corporativos como Microsoft Defender fazendo pre-fetching de Magic Links e invalidando tokens antes do clique humano*, ou *webhooks de pagamento chegando fora de ordem em filas distribuídas*).
- O desenvolvedor sofre em produção com bugs que usuários reais da internet já reportaram centenas de vezes.

## 💡 A Solução: OpenStories

O **OpenStories** transforma a experiência coletiva de engenheiros e usuários em uma **biblioteca viva de User Stories respaldadas por evidências reais**. 

Cada história contém:
1. **User Story Formal**: `As a... I want... So that...`
2. **Critérios de Aceitação Gherkin**: `Scenario: ... Given ... When ... Then`
3. **Casos de Borda de Produção**: Armadilhas documentadas em post-mortems reais.
4. **Pacote de Evidências com Citações**: Links para GitHub issues, Reddit threads e post-mortems de empresas reais.
5. **Score Empírico de Demanda**: Calculado a partir da frequência e gravidade dos relatos.

---

## 🚀 Quickstart em 30 Segundos

### 1. Instalação do Binário
Você pode compilar localmente ou instalar via Go:
```bash
go install github.com/gabrielrondon/openstories/cmd/openstories@latest
```
Ou clonando o repositório:
```bash
git clone https://github.com/gabrielrondon/openstories.git
cd openstories
make build
./openstories
```

### 2. Buscar Histórias pelo Terminal
```bash
# Busca inteligente por problema ou conceito
openstories search "idempotência em APIs"

# Filtrando por indústria
openstories search "webhooks fora de ordem" --industry fintech

# Visualizar todos os detalhes de uma história
openstories get OS-DEV-001
```

---

## 🛠️ Integração com Claude Desktop & Cursor (1-Click MCP)

O OpenStories implementa nativamente o **Model Context Protocol (MCP)** sobre `stdio`. Ele roda **localmente na sua máquina, com zero custo e zero infraestrutura de nuvem**.

### Instalação Automática (1 Comando):

```bash
# Configura o Claude Desktop automaticamente
openstories install claude

# Configura o Cursor IDE automaticamente
openstories install cursor
```

Pronto! Ao abrir o Claude Desktop ou Cursor, você terá as ferramentas:
- `search_stories`: Busca semântica de requisitos reais por indústria/dor.
- `get_story`: Recupera critérios Gherkin e links de incidentes.
- `evaluate_spec`: Avalia o plano da feature e aponta pontos cegos de produção.
- `save_custom_story`: Salva histórias internas na pasta `.openstories/stories/` do seu projeto.

### Configuração Manual (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "openstories": {
      "command": "openstories",
      "args": ["mcp"]
    }
  }
}
```

---

## 🛡️ O Spec Stress-Tester (Linter de Realidade)

Antes de escrever qualquer linha de código, você ou o agente podem passar a especificação pelo Stress-Tester:

```bash
openstories eval --industry devtools "Feature: Criar endpoint POST /charge com retry se houver timeout"
```

### Exemplo de Retorno do OpenStories:
```
🛡️  RELATÓRIO DE STRESS-TEST DE ESPECIFICAÇÃO
Score de Cobertura de Realidade: 15/100
Resumo: Atenção: Alta vulnerabilidade a incidentes de produção. A especificação negligenciou múltiplos casos de borda severos com alto volume de reclamações reais na internet.

🚨 CASOS DE BORDA CRÍTICOS NÃO TRATADOS NA SPEC:
  ❌ [OS-DEV-001] Idempotência obrigatória: Expiração da chave de idempotência após a janela de retenção (24 horas).
  ❌ [OS-DEV-001] Idempotência obrigatória: Concorrência real de milissegundos tratada via distributed lock.
  ❌ [OS-DEV-001] Idempotência obrigatória: Salvamento do cabeçalho de resposta original e status code original.

💡 RECOMENDAÇÕES:
  • Adicionar testes automatizados específicos para os casos de borda listados nos alertas.
  • Incorporar na documentação do PR ou RFC as respostas aos itens do checklist não atendidos.
```

---

## 🌐 Interface Web Embutida (`openstories serve`)

O OpenStories vem com uma interface web moderna, interativa e responsiva embutida diretamente no binário (via `embed.FS`):

```bash
openstories serve --port 8080
```
Acesse `http://localhost:8080` para:
- Navegar interativamente pelo catálogo de histórias.
- Filtrar por indústrias e domínios.
- Inspecionar critérios Gherkin e citações originais.
- Usar o playground visual do **Spec Stress-Tester**.

> 💡 **Deploy para Produção**: O código da pasta `web/` está preparado para deploy estático ou via subdomínio **`openstories.tuturama.com`**.

---

## 🌾 Harvester: Mineração Automatizada da Web

O OpenStories possui um coletor que lê discussões e issues públicas e rascunha histórias prontas com evidências:

```bash
openstories harvest --repo stripe/stripe-go --label bug --industry devtools
```

---

## 🌳 Árvore de Classificação

```
stories/
├── devtools/                        # APIs, Auth, CI/CD, Observability
│   ├── apis-and-sdks/               # OS-DEV-001: Idempotency Keys
│   └── auth-and-iam/                # OS-DEV-002: Magic Link Pre-fetching
├── ai-infra/                        # LLM Gateways, Agents, RAG
│   ├── agentic-workflows/           # OS-AI-001: Tool Calling Loop Breaker
│   └── llm-gateways/                # OS-AI-002: Token Rate Limit Fallbacks
├── fintech/                         # Pagamentos, Conciliação, Fraude
│   ├── payments-and-checkout/       # OS-FIN-001: Out-of-Order Webhooks
│   └── reconciliation/              # OS-FIN-002: Partial Refund MDR
└── b2b-saas/                        # Multi-tenancy, Auditoria, Workspaces
    ├── multi-tenancy/               # OS-SaaS-001: Slug Hijacking Prevention
    └── audit-logs/                  # OS-SaaS-002: SOC2 Immutable Audit Trail
```

---

## 🤝 Como Contribuir

1. Faça um Fork do projeto.
2. Crie uma nova história em `stories/<industry>/<domain>/<nome>.story.md`.
3. Garanta que a história contenha **links e citações de evidências reais**.
4. Rode os testes: `make test`.
5. Abra seu Pull Request!

---

## 📄 Licença

Distribuído sob a licença **MIT**. Veja [`LICENSE`](LICENSE) para mais detalhes.

Criado com dedicação por **Gabriel Rondon** & **Tuturama**.
