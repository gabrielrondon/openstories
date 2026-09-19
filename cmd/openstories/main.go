package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gabrielrondon/openstories/internal/evaluator"
	"github.com/gabrielrondon/openstories/internal/harvester"
	"github.com/gabrielrondon/openstories/internal/installer"
	"github.com/gabrielrondon/openstories/internal/mcp"
	"github.com/gabrielrondon/openstories/internal/server"
	"github.com/gabrielrondon/openstories/internal/store"
	"github.com/gabrielrondon/openstories/stories"
)

const banner = `
   ____                   ____  _             _           
  / __ \                 / ____|| |           (_)          
 | |  | | _ __   ___  _ | (___  | |_  ___  _ __ _  ___  ___ 
 | |  | || '_ \ / _ \| '_ \___ \ | __|/ _ \| '__| |/ _ \/ __|
 | |__| || |_) |  __/| | | |___) || |_| (_) | |  | |  __/\__ \
  \____/ | .__/ \___||_| |_|_____/ \__|\___/|_|  |_|\___||___/
         | |   Open-Source Evidence-Backed User Stories
         |_|   for AI Coding Agents & Systems • v1.0.0
`

func main() {
	// Initialize store
	st := store.New(".openstories/stories")

	// Load compile-time embedded stories
	if err := st.LoadFromFS(stories.EmbeddedFS, ".", false); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load embedded stories: %v\n", err)
	}

	// Load any local project-level custom stories (.openstories/stories or .stories)
	_ = st.LoadLocalDir(".openstories/stories")
	_ = st.LoadLocalDir(".stories")

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "mcp":
		srv := mcp.NewServer(st)
		if err := srv.Run(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			os.Exit(1)
		}

	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		industry := searchCmd.String("industry", "", "Filtrar por indústria")
		domain := searchCmd.String("domain", "", "Filtrar por domínio")
		jsonOut := searchCmd.Bool("json", false, "Exibir resultado em JSON")
		searchCmd.Parse(os.Args[2:])

		query := strings.Join(searchCmd.Args(), " ")
		results := st.Search(query, *industry, *domain, 0)

		if *jsonOut {
			type JSONItem struct {
				ID          string   `json:"id"`
				Title       string   `json:"title"`
				Industry    string   `json:"industry"`
				Domain      string   `json:"domain"`
				DemandScore float64  `json:"demand_score"`
				EdgeCases   []string `json:"edge_cases"`
			}
			out := make([]JSONItem, 0, len(results))
			for _, r := range results {
				out = append(out, JSONItem{
					ID:          r.Story.ID,
					Title:       r.Story.Title,
					Industry:    r.Story.Industry,
					Domain:      r.Story.Domain,
					DemandScore: r.Story.DemandScore,
					EdgeCases:   r.Story.EdgeCases,
				})
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n🔍 Encontradas %d histórias para: \"%s\"\n\n", len(results), query)
		for _, r := range results {
			s := r.Story
			fmt.Printf("• \033[1;36m[%s]\033[0m \033[1m%s\033[0m (Score: \033[1;32m%.1f\033[0m)\n", s.ID, s.Title, s.DemandScore)
			fmt.Printf("  Indústria: %s | Domínio: %s | Evidências: %d\n", s.Industry, s.Domain, len(s.Evidence))
			if s.Statement.IWant != "" {
				fmt.Printf("  \033[90mQuero: %s\033[0m\n", s.Statement.IWant)
			}
			fmt.Println()
		}

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Uso: openstories get <ID>")
			os.Exit(1)
		}
		id := os.Args[2]
		s, ok := st.Get(id)
		if !ok {
			fmt.Printf("História com ID '%s' não encontrada.\n", id)
			os.Exit(1)
		}

		fmt.Print(banner)
		fmt.Printf("ID: \033[1;36m%s\033[0m | Indústria: %s | Domínio: %s | Score: \033[1;32m%.1f\033[0m\n", s.ID, s.Industry, s.Domain, s.DemandScore)
		fmt.Printf("Título: \033[1m%s\033[0m\n\n", s.Title)

		fmt.Println("📋 \033[1mDECLARAÇÃO DA USER STORY:\033[0m")
		fmt.Printf("  Como: %s\n", s.Statement.AsA)
		fmt.Printf("  Quero: %s\n", s.Statement.IWant)
		fmt.Printf("  Para que: %s\n\n", s.Statement.SoThat)

		if len(s.AcceptanceCriteria) > 0 {
			fmt.Println("✅ \033[1mCRITÉRIOS DE ACEITAÇÃO (GHERKIN):\033[0m")
			for _, ac := range s.AcceptanceCriteria {
				fmt.Printf("  • Cenário: %s\n", ac.Scenario)
				fmt.Printf("    DADO: %s\n", ac.Given)
				fmt.Printf("    QUANDO: %s\n", ac.When)
				fmt.Printf("    ENTÃO: %s\n", ac.Then)
			}
			fmt.Println()
		}

		if len(s.EdgeCases) > 0 {
			fmt.Println("⚠️ \033[1mCASOS DE BORDA OBSERVADOS EM PRODUÇÃO:\033[0m")
			for _, ec := range s.EdgeCases {
				fmt.Printf("  • %s\n", ec)
			}
			fmt.Println()
		}

		if len(s.Evidence) > 0 {
			fmt.Println("🔗 \033[1mEVIDÊNCIAS E GROUNDING REAL:\033[0m")
			for _, ev := range s.Evidence {
				fmt.Printf("  • [\"%s\"] - %s (%s)\n", ev.Quote, ev.Source, ev.Date)
			}
			fmt.Println()
		}

	case "eval":
		evalCmd := flag.NewFlagSet("eval", flag.ExitOnError)
		industry := evalCmd.String("industry", "", "Indústria de contexto")
		filePath := evalCmd.String("file", "", "Caminho de arquivo com a spec markdown")
		evalCmd.Parse(os.Args[2:])

		var specContent string
		if *filePath != "" {
			data, err := os.ReadFile(*filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao ler arquivo %s: %v\n", *filePath, err)
				os.Exit(1)
			}
			specContent = string(data)
		} else {
			specContent = strings.Join(evalCmd.Args(), " ")
		}

		if strings.TrimSpace(specContent) == "" {
			fmt.Println("Por favor informe o texto da especificação ou o caminho de um arquivo (--file <caminho>).")
			os.Exit(1)
		}

		ev := evaluator.New(st)
		res := ev.Evaluate(specContent, *industry, "")

		fmt.Println("\n🛡️  \033[1mRELATÓRIO DE STRESS-TEST DE ESPECIFICAÇÃO\033[0m")
		fmt.Printf("Score de Cobertura de Realidade: \033[1;32m%d/100\033[0m\n", res.Score)
		fmt.Printf("Resumo: %s\n\n", res.Summary)

		if len(res.EdgeCaseAlerts) > 0 {
			fmt.Println("🚨 \033[1;31mCASOS DE BORDA CRÍTICOS NÃO TRATADOS NA SPEC:\033[0m")
			for _, alert := range res.EdgeCaseAlerts {
				fmt.Printf("  ❌ %s\n", alert)
			}
			fmt.Println()
		}

		if len(res.Recommendations) > 0 {
			fmt.Println("💡 \033[1;33mRECOMENDAÇÕES:\033[0m")
			for _, rec := range res.Recommendations {
				fmt.Printf("  • %s\n", rec)
			}
			fmt.Println()
		}

	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Uso: openstories install [claude|cursor]")
			os.Exit(1)
		}
		target := strings.ToLower(os.Args[2])
		switch target {
		case "claude":
			if err := installer.InstallClaude(); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao configurar Claude Desktop: %v\n", err)
				os.Exit(1)
			}
		case "cursor":
			if err := installer.InstallCursor(); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao configurar Cursor: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Target desconhecido: %s. Use 'claude' ou 'cursor'.\n", target)
			os.Exit(1)
		}

	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.Int("port", 8080, "Porta do servidor web")
		serveCmd.Parse(os.Args[2:])

		srv := server.New(st, *port)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Erro no servidor HTTP: %v\n", err)
			os.Exit(1)
		}

	case "harvest":
		harvestCmd := flag.NewFlagSet("harvest", flag.ExitOnError)
		repo := harvestCmd.String("repo", "", "Repositório GitHub (ex: stripe/stripe-go)")
		label := harvestCmd.String("label", "bug", "Label de issue a filtrar")
		industry := harvestCmd.String("industry", "devtools", "Indústria da história minerada")
		domain := harvestCmd.String("domain", "general", "Domínio da história minerada")
		harvestCmd.Parse(os.Args[2:])

		if *repo == "" {
			fmt.Println("Por favor, especifique o repositório com --repo <dono/repo>")
			os.Exit(1)
		}

		h := harvester.New(st)
		story, err := h.HarvestGitHub(*repo, *label, *industry, *domain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Falha na mineração: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("História gerada: [%s] %s\n", story.ID, story.Title)

	case "list":
		all := st.ListAll()
		fmt.Printf("\n📚 Total de histórias cadastradas: %d\n\n", len(all))
		for _, s := range all {
			fmt.Printf("• \033[1;36m%-12s\033[0m \033[1m%-50s\033[0m [\033[33m%s\033[0m] (Score: %.1f)\n", s.ID, s.Title, s.Industry, s.DemandScore)
		}
		fmt.Println()

	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Print(banner)
	fmt.Println("COMANDOS DISPONÍVEIS:")
	fmt.Println("  openstories mcp                  Inicia o servidor MCP stdio para Claude, Cursor ou Antigravity")
	fmt.Println("  openstories search <termos>      Busca histórias por problema, tecnologia ou dor de usuário")
	fmt.Println("  openstories get <id>             Exibe detalhes, critérios Gherkin e links de evidência real")
	fmt.Println("  openstories eval --file <spec>   Stress-testa uma especificação contra casos de borda reais")
	fmt.Println("  openstories install [claude]     Configura o MCP automaticamente no Claude Desktop")
	fmt.Println("  openstories install [cursor]     Configura o MCP automaticamente no Cursor IDE")
	fmt.Println("  openstories serve [--port 8080]  Inicia a interface web e API REST (openstories.tuturama.com)")
	fmt.Println("  openstories harvest --repo <r>   Minera discussões do GitHub e sintetiza novas histórias")
	fmt.Println("  openstories list                 Lista todas as histórias cadastradas no repositório")
	fmt.Println()
}
