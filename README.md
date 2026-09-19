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

**Ground your AI agents in real-world human problems and production battle-scars.**  
*The first open-source engine of evidence-backed user stories grounded in public GitHub issues, production post-mortems, and community outages — featuring a native Model Context Protocol (MCP) server and reality stress-tester.*

[Web Dashboard](https://openstories.tuturama.com) • [Quickstart in 30s](#-quickstart-in-30-seconds) • [Claude & Cursor Setup](#-1-click-setup-for-claude-code--cursor) • [Spec Stress-Tester](#-reality-stress-tester-the-production-linter) • [Taxonomy](#-curated-taxonomy)

</div>

---

## 💥 The Problem: The "Happy Path" Bias in AI Coding Agents

When you prompt **Claude Code**, **Cursor**, **Codex**, or **Antigravity** to architect a feature:
- They generate toy specifications and "happy path" assumptions.
- They have zero context on bizarre, real-world production traps (e.g. *enterprise mail scanners like Microsoft Defender pre-fetching magic links and burning one-time tokens before humans click*, or *payment webhooks arriving out-of-order across distributed message queues*).
- Engineering teams face critical production outages caused by failure modes that other developers have already documented hundreds of times on GitHub and Reddit.

## 💡 The Solution: OpenStories

**OpenStories** turns the collective operational scars of the software industry into an **executable library of evidence-backed user stories**.

Every story provides:
1. **User Story Statement**: `As a... I want... So that...`
2. **Gherkin Acceptance Criteria**: Executable `Scenario: ... Given ... When ... Then`
3. **Observed Production Edge Cases**: Specific failure modes extracted from real post-mortems.
4. **Grounded Evidence with Verifiable Citations**: Direct links and quotes from GitHub issues, Reddit rants, and official incident reports.
5. **Empirical Demand Score**: Quantitative score calculated from community frequency and outage severity.

---

## 🚀 Quickstart in 30 Seconds

### 1. Installation
Install directly via Go or compile from source:
```bash
go install github.com/gabrielrondon/openstories/cmd/openstories@latest
```
Or clone and build the zero-dependency static binary:
```bash
git clone https://github.com/gabrielrondon/openstories.git
cd openstories
make build
./openstories
```

### 2. Search Stories via CLI
```bash
# Semantic search across all curated stories
openstories search "idempotency on mutation endpoints"

# Filter by industry
openstories search "out of order webhooks" --industry fintech

# Inspect complete Gherkin criteria & citations
openstories get OS-DEV-001
```

---

## 🛠️ 1-Click Setup for Claude Code & Cursor

OpenStories natively implements the official **Model Context Protocol (MCP)** specification over `stdio`. It runs **100% locally on your machine with zero cloud costs or infrastructure**.

### Instant Configuration:
```bash
# Configures both Claude Code CLI (~/.claude.json) and Claude Desktop
openstories install claude

# Configures Cursor IDE (~/.cursor/mcp.json)
openstories install cursor
```

Restart Claude Code or Cursor, and OpenStories tools are instantly available:
- `search_stories`: Retrieve verified stories matching your current task.
- `get_story`: Load full Gherkin criteria and incident evidence.
- `evaluate_spec`: Reality stress-tester that audits your feature draft against production failures.
- `save_custom_story`: Persist private company stories in your local `.openstories/stories/` folder.

### System Prompt Snippet (Recommended):
Add this line to your project's `CLAUDE.md`, `.cursorrules`, or system prompt:
> *"Before implementing any feature or architectural plan, query the `openstories` MCP tools (`search_stories` and `evaluate_spec`) to ground your specification in real-world production incident edge cases."*

---

## 🛡️ Reality Stress-Tester (The Production Linter)

Before writing code, run your architectural draft or AI-generated PR plan through the Reality Stress-Tester:

```bash
openstories eval --industry devtools "Feature: POST /charge endpoint that retries Stripe on network timeout"
```

### Output:
```
🛡️  REALITY STRESS-TEST REPORT
Reality Coverage Score: 15/100
Summary: Severe Production Vulnerability Alert: Multiple high-impact failure modes observed in real-world outages were ignored in this specification.

🚨 PRODUCTION BLIND SPOTS UNHANDLED IN SPEC:
  ❌ [OS-DEV-001] Mandatory Idempotency Keys: Idempotency key TTL expiration after the standard retention window (24 hours).
  ❌ [OS-DEV-001] Mandatory Idempotency Keys: Sub-millisecond race conditions mitigated via distributed locks.
  ❌ [OS-DEV-001] Mandatory Idempotency Keys: Full header and status code preservation (replaying original response headers).

💡 ACTIONABLE RECOMMENDATIONS:
  • Add explicit integration and failure-injection tests covering the flagged production edge cases.
  • Document explicit architectural answers to the quality checklist rubrics before writing code.
```

---

## 🌐 Embedded Web Dashboard (`openstories serve`)

OpenStories features an embedded, responsive dark-mode web application and REST API compiled directly into the binary:

```bash
openstories serve --port 8080
```
Open `http://localhost:8080` in your browser to:
- Browse and search the evidence-backed catalog.
- Filter by industry and domain.
- Invert specs interactively in the **Spec Stress-Tester Playground**.
- Ready for deployment on **`openstories.tuturama.com`**.

---

## 🌾 Web Harvester: Mining Real World Issues

Automatically ingest public issues and synthesize draft user stories:

```bash
openstories harvest --repo stripe/stripe-go --label bug --industry devtools
```

---

## 🌳 Curated Taxonomy

```
stories/
├── devtools/                        # APIs, Auth, CI/CD, Observability
│   ├── apis-and-sdks/               # OS-DEV-001: Mandatory Idempotency Keys
│   └── auth-and-iam/                # OS-DEV-002: Magic Link Scanner Pre-Fetching
├── ai-infra/                        # LLM Gateways, Agents, Tool Calling
│   ├── agentic-workflows/           # OS-AI-001: Infinite Tool-Calling Loop Breaker
│   └── llm-gateways/                # OS-AI-002: Token Rate Limit Fallback Cascades
├── fintech/                         # Payments, Ledgers, Reconciliation
│   ├── payments-and-checkout/       # OS-FIN-001: Out-of-Order Webhook Delivery
│   └── reconciliation/              # OS-FIN-002: Partial Refund MDR Reconciliation
└── b2b-saas/                        # Multi-Tenancy, Compliance, Workspaces
    ├── multi-tenancy/               # OS-SaaS-001: Subdomain Hijacking & Homoglyphs
    └── audit-logs/                  # OS-SaaS-002: Cryptographic Tamper-Proof Audit Trail
```

---

## 🌍 Multi-Language & Internationalization

Stories support an explicit `locale` tag (defaulting to `en`). Community contributions in any language (Portuguese, Spanish, Japanese, etc.) can be contributed as localized companion stories (e.g. `OS-DEV-001.pt.story.md` with `locale: pt-BR`).

---

## 🤝 Contributing

We welcome community-contributed stories!
1. Fork the repository.
2. Create a `.story.md` file in `stories/<industry>/<domain>/`.
3. Include **verifiable links and citations** to real GitHub issues, forum threads, or post-mortems.
4. Run tests: `make test`.
5. Submit your Pull Request!

---

## 📄 License

Distributed under the **MIT License**. See [`LICENSE`](LICENSE) for details.

Maintained with ❤️ by **Gabriel Rondon** & **Tuturama**.
