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
	st := store.New(".openstories/stories")

	if err := st.LoadFromFS(stories.EmbeddedFS, ".", false); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load embedded stories: %v\n", err)
	}

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
		industry := searchCmd.String("industry", "", "Filter by industry")
		domain := searchCmd.String("domain", "", "Filter by domain")
		locale := searchCmd.String("locale", "", "Filter by locale (e.g. en, pt-br)")
		jsonOut := searchCmd.Bool("json", false, "Output results as formatted JSON")
		searchCmd.Parse(os.Args[2:])

		query := strings.Join(searchCmd.Args(), " ")
		results := st.SearchWithLocale(query, *industry, *domain, 0, *locale)

		if *jsonOut {
			type JSONItem struct {
				ID          string   `json:"id"`
				Locale      string   `json:"locale"`
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
					Locale:      r.Story.Locale,
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

		fmt.Printf("\n🔍 Found %d stories matching: \"%s\"\n\n", len(results), query)
		for _, r := range results {
			s := r.Story
			fmt.Printf("• \033[1;36m[%s]\033[0m \033[1m%s\033[0m (Score: \033[1;32m%.1f\033[0m)\n", s.ID, s.Title, s.DemandScore)
			fmt.Printf("  Industry: %s | Domain: %s | Evidence: %d sources\n", s.Industry, s.Domain, len(s.Evidence))
			if s.Statement.IWant != "" {
				fmt.Printf("  \033[90mWant: %s\033[0m\n", s.Statement.IWant)
			}
			fmt.Println()
		}

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: openstories get <ID>")
			os.Exit(1)
		}
		id := os.Args[2]
		s, ok := st.Get(id)
		if !ok {
			fmt.Printf("Story with ID '%s' not found.\n", id)
			os.Exit(1)
		}

		fmt.Print(banner)
		fmt.Printf("ID: \033[1;36m%s\033[0m | Industry: %s | Domain: %s | Demand Score: \033[1;32m%.1f\033[0m\n", s.ID, s.Industry, s.Domain, s.DemandScore)
		fmt.Printf("Title: \033[1m%s\033[0m\n\n", s.Title)

		fmt.Println("📋 \033[1mUSER STORY STATEMENT:\033[0m")
		fmt.Printf("  As a: %s\n", s.Statement.AsA)
		fmt.Printf("  I want: %s\n", s.Statement.IWant)
		fmt.Printf("  So that: %s\n\n", s.Statement.SoThat)

		if len(s.AcceptanceCriteria) > 0 {
			fmt.Println("✅ \033[1mGHERKIN ACCEPTANCE CRITERIA:\033[0m")
			for _, ac := range s.AcceptanceCriteria {
				fmt.Printf("  • Scenario: %s\n", ac.Scenario)
				fmt.Printf("    GIVEN: %s\n", ac.Given)
				fmt.Printf("    WHEN:  %s\n", ac.When)
				fmt.Printf("    THEN:  %s\n", ac.Then)
			}
			fmt.Println()
		}

		if len(s.EdgeCases) > 0 {
			fmt.Println("⚠️ \033[1mOBSERVED PRODUCTION EDGE CASES:\033[0m")
			for _, ec := range s.EdgeCases {
				fmt.Printf("  • %s\n", ec)
			}
			fmt.Println()
		}

		if len(s.Evidence) > 0 {
			fmt.Println("🔗 \033[1mGROUNDED REAL-WORLD EVIDENCE:\033[0m")
			for _, ev := range s.Evidence {
				fmt.Printf("  • [\"%s\"] - %s (%s)\n", ev.Quote, ev.Source, ev.Date)
			}
			fmt.Println()
		}

	case "eval":
		evalCmd := flag.NewFlagSet("eval", flag.ExitOnError)
		industry := evalCmd.String("industry", "", "Context industry (e.g. devtools, fintech, ai-infra)")
		filePath := evalCmd.String("file", "", "Path to markdown spec file")
		evalCmd.Parse(os.Args[2:])

		var specContent string
		if *filePath != "" {
			data, err := os.ReadFile(*filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", *filePath, err)
				os.Exit(1)
			}
			specContent = string(data)
		} else {
			specContent = strings.Join(evalCmd.Args(), " ")
		}

		if strings.TrimSpace(specContent) == "" {
			fmt.Println("Please provide the specification text or a file path (--file <path>).")
			os.Exit(1)
		}

		ev := evaluator.New(st)
		res := ev.Evaluate(specContent, *industry, "")

		fmt.Println("\n🛡️  \033[1mREALITY STRESS-TEST REPORT\033[0m")
		fmt.Printf("Reality Coverage Score: \033[1;32m%d/100\033[0m\n", res.Score)
		fmt.Printf("Summary: %s\n\n", res.Summary)

		if len(res.EdgeCaseAlerts) > 0 {
			fmt.Println("🚨 \033[1;31mPRODUCTION BLIND SPOTS UNHANDLED IN SPEC:\033[0m")
			for _, alert := range res.EdgeCaseAlerts {
				fmt.Printf("  ❌ %s\n", alert)
			}
			fmt.Println()
		}

		if len(res.Recommendations) > 0 {
			fmt.Println("💡 \033[1;33mACTIONABLE RECOMMENDATIONS:\033[0m")
			for _, rec := range res.Recommendations {
				fmt.Printf("  • %s\n", rec)
			}
			fmt.Println()
		}

	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: openstories install [claude|claude-code|cursor]")
			os.Exit(1)
		}
		target := strings.ToLower(os.Args[2])
		switch target {
		case "claude", "claude-code":
			if err := installer.InstallClaude(); err != nil {
				fmt.Fprintf(os.Stderr, "Error configuring Claude: %v\n", err)
				os.Exit(1)
			}
		case "cursor":
			if err := installer.InstallCursor(); err != nil {
				fmt.Fprintf(os.Stderr, "Error configuring Cursor: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unknown target: %s. Use 'claude', 'claude-code', or 'cursor'.\n", target)
			os.Exit(1)
		}

	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.Int("port", 8080, "Web server port")
		serveCmd.Parse(os.Args[2:])

		srv := server.New(st, *port)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
			os.Exit(1)
		}

	case "harvest":
		harvestCmd := flag.NewFlagSet("harvest", flag.ExitOnError)
		repo := harvestCmd.String("repo", "", "Target GitHub repository (e.g. stripe/stripe-go)")
		label := harvestCmd.String("label", "bug", "Issue label filter")
		industry := harvestCmd.String("industry", "devtools", "Target story industry")
		domain := harvestCmd.String("domain", "general", "Target story domain")
		harvestCmd.Parse(os.Args[2:])

		if *repo == "" {
			fmt.Println("Please specify the target repository with --repo <owner/repo>")
			os.Exit(1)
		}

		h := harvester.New(st)
		story, err := h.HarvestGitHub(*repo, *label, *industry, *domain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Harvest failure: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Synthesized Story: [%s] %s\n", story.ID, story.Title)

	case "list":
		all := st.ListAll()
		fmt.Printf("\n📚 Total Registered Stories: %d\n\n", len(all))
		for _, s := range all {
			fmt.Printf("• \033[1;36m%-12s\033[0m \033[1m%-60s\033[0m [\033[33m%s\033[0m] (Score: %.1f)\n", s.ID, s.Title, s.Industry, s.DemandScore)
		}
		fmt.Println()

	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Print(banner)
	fmt.Println("COMMANDS:")
	fmt.Println("  openstories mcp                  Run native MCP stdio server for Claude, Cursor, or Antigravity")
	fmt.Println("  openstories search <terms>       Search stories by problem, technology, or user pain point")
	fmt.Println("  openstories get <id>             Inspect full Gherkin criteria and verified evidence links")
	fmt.Println("  openstories eval --file <spec>   Stress-test a feature specification against real production failures")
	fmt.Println("  openstories install [claude]     1-Click setup for Claude Code CLI and Claude Desktop")
	fmt.Println("  openstories install [cursor]     1-Click setup for Cursor IDE")
	fmt.Println("  openstories serve [--port 8080]  Launch interactive Web Dashboard & REST API (openstories.tuturama.com)")
	fmt.Println("  openstories harvest --repo <r>   Mine public GitHub issues into structured user stories")
	fmt.Println("  openstories list                 List all indexed stories across industries")
	fmt.Println()
}
