package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gabrielrondon/openstories/internal/evaluator"
	"github.com/gabrielrondon/openstories/internal/model"
	"github.com/gabrielrondon/openstories/internal/store"
)

// Server implements a Model Context Protocol (MCP) server over standard I/O (stdio).
type Server struct {
	store     *store.Store
	evaluator *evaluator.Evaluator
	logger    *log.Logger
}

// NewServer creates a new MCP Server.
func NewServer(s *store.Store) *Server {
	return &Server{
		store:     s,
		evaluator: evaluator.New(s),
		logger:    log.New(os.Stderr, "[openstories-mcp] ", log.LstdFlags),
	}
}

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]PropertyDef `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

type PropertyDef struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Run listens for JSON-RPC messages on stdin and writes responses to stdout.
func (s *Server) Run(stdin io.Reader, stdout io.Writer) error {
	s.logger.Println("OpenStories MCP Server starting on stdio...")
	scanner := bufio.NewScanner(stdin)

	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 4*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.logger.Printf("Invalid JSON: %v", err)
			continue
		}

		resp := s.handleRequest(&req)
		if resp != nil {
			bytes, err := json.Marshal(resp)
			if err != nil {
				s.logger.Printf("Error marshaling response: %v", err)
				continue
			}
			fmt.Fprintf(stdout, "%s\n", string(bytes))
		}
	}

	return scanner.Err()
}

func (s *Server) handleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo": map[string]string{
					"name":    "openstories",
					"version": "1.0.0",
				},
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
			},
		}

	case "notifications/initialized":
		return nil

	case "ping":
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]interface{}{},
		}

	case "tools/list":
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": s.getToolDefinitions(),
			},
		}

	case "tools/call":
		return s.handleToolCall(req)

	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *Server) getToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		{
			Name: "search_stories",
			Description: "Search the OpenStories evidence-backed library for real-world user stories, production post-mortems, and edge cases by problem, technology, or domain.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"query": {
						Type:        "string",
						Description: "Search terms, technical concepts, user pain points, or feature ideas (e.g. 'idempotency retry', 'magic link pre-fetch', 'out-of-order webhooks').",
					},
					"industry": {
						Type:        "string",
						Description: "Optional industry filter (e.g. 'devtools', 'fintech', 'ai-infra', 'b2b-saas').",
					},
					"domain": {
						Type:        "string",
						Description: "Optional domain filter within an industry.",
					},
					"locale": {
						Type:        "string",
						Description: "Optional language/locale filter (defaults to 'en').",
					},
				},
			},
		},
		{
			Name: "get_story",
			Description: "Retrieve a complete evidence-backed user story by ID, including persona, Gherkin acceptance criteria, production edge cases, and real-world citations.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"id": {
						Type:        "string",
						Description: "Story identifier (e.g. 'OS-DEV-001', 'OS-AI-001', 'OS-FIN-001').",
					},
				},
				Required: []string{"id"},
			},
		},
		{
			Name: "evaluate_spec",
			Description: "Empirical reality stress-tester: evaluates a feature plan, PR description, or technical specification against documented real-world production failures and edge cases.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"spec": {
						Type:        "string",
						Description: "The architectural spec, feature draft, or PR description to evaluate.",
					},
					"industry": {
						Type:        "string",
						Description: "Target industry (e.g. 'devtools', 'fintech', 'ai-infra', 'b2b-saas').",
					},
					"domain": {
						Type:        "string",
						Description: "Optional specific domain.",
					},
				},
				Required: []string{"spec"},
			},
		},
		{
			Name: "list_taxonomies",
			Description: "List all indexed industries and domains currently registered in OpenStories.",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]PropertyDef{},
			},
		},
		{
			Name: "save_custom_story",
			Description: "Persist a private custom user story into the local workspace (.openstories/stories/) so AI coding agents can ground and remember internal company rules.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertyDef{
					"id":          {Type: "string", Description: "Unique story identifier (e.g. 'CUSTOM-001')"},
					"industry":    {Type: "string", Description: "Industry name (e.g. 'internal', 'fintech')"},
					"domain":      {Type: "string", Description: "Domain or subsystem name"},
					"title":       {Type: "string", Description: "Clear title of the user story"},
					"as_a":        {Type: "string", Description: "Target persona (As a...)"},
					"i_want":      {Type: "string", Description: "Action or capability (I want...)"},
					"so_that":     {Type: "string", Description: "Core value or benefit (So that...)"},
					"edge_cases":  {Type: "string", Description: "Semicolon-separated production edge cases"},
					"evidence":    {Type: "string", Description: "Real feedback, support ticket quotes, or context"},
				},
				Required: []string{"id", "industry", "domain", "title", "as_a", "i_want", "so_that"},
			},
		},
	}
}

func (s *Server) handleToolCall(req *JSONRPCRequest) *JSONRPCResponse {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid parameters: %v", err),
			},
		}
	}

	switch params.Name {
	case "search_stories":
		query, _ := params.Arguments["query"].(string)
		industry, _ := params.Arguments["industry"].(string)
		domain, _ := params.Arguments["domain"].(string)
		locale, _ := params.Arguments["locale"].(string)

		results := s.store.SearchWithLocale(query, industry, domain, 0, locale)
		type SummaryItem struct {
			ID            string   `json:"id"`
			Locale        string   `json:"locale"`
			Title         string   `json:"title"`
			Industry      string   `json:"industry"`
			Domain        string   `json:"domain"`
			DemandScore   float64  `json:"demand_score"`
			EdgeCases     []string `json:"edge_cases"`
			EvidenceCount int      `json:"evidence_count"`
		}

		out := make([]SummaryItem, 0, len(results))
		for _, r := range results {
			out = append(out, SummaryItem{
				ID:            r.Story.ID,
				Locale:        r.Story.Locale,
				Title:         r.Story.Title,
				Industry:      r.Story.Industry,
				Domain:        r.Story.Domain,
				DemandScore:   r.Story.DemandScore,
				EdgeCases:     r.Story.EdgeCases,
				EvidenceCount: len(r.Story.Evidence),
			})
		}

		data, _ := json.MarshalIndent(out, "", "  ")
		return createToolResponse(req.ID, string(data))

	case "get_story":
		id, _ := params.Arguments["id"].(string)
		st, ok := s.store.Get(id)
		if !ok {
			return createToolResponse(req.ID, fmt.Sprintf("Story with ID '%s' not found.", id))
		}

		data, _ := json.MarshalIndent(st, "", "  ")
		return createToolResponse(req.ID, string(data))

	case "evaluate_spec":
		specText, _ := params.Arguments["spec"].(string)
		industry, _ := params.Arguments["industry"].(string)
		domain, _ := params.Arguments["domain"].(string)

		evaluation := s.evaluator.Evaluate(specText, industry, domain)
		data, _ := json.MarshalIndent(evaluation, "", "  ")
		return createToolResponse(req.ID, string(data))

	case "list_taxonomies":
		taxonomies := s.store.ListTaxonomies()
		data, _ := json.MarshalIndent(taxonomies, "", "  ")
		return createToolResponse(req.ID, string(data))

	case "save_custom_story":
		id, _ := params.Arguments["id"].(string)
		industry, _ := params.Arguments["industry"].(string)
		domain, _ := params.Arguments["domain"].(string)
		title, _ := params.Arguments["title"].(string)
		asA, _ := params.Arguments["as_a"].(string)
		iWant, _ := params.Arguments["i_want"].(string)
		soThat, _ := params.Arguments["so_that"].(string)
		edgeCasesRaw, _ := params.Arguments["edge_cases"].(string)
		evidenceRaw, _ := params.Arguments["evidence"].(string)

		edgeCases := make([]string, 0)
		for _, part := range strings.Split(edgeCasesRaw, ";") {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				edgeCases = append(edgeCases, trimmed)
			}
		}

		evList := make([]model.Evidence, 0)
		if strings.TrimSpace(evidenceRaw) != "" {
			evList = append(evList, model.Evidence{
				Source: "internal_feedback",
				Type:   "custom_evidence",
				Quote:  evidenceRaw,
				Date:   "custom",
			})
		}

		customStory := &model.Story{
			ID:          id,
			Locale:      "en",
			Industry:    industry,
			Domain:      domain,
			Title:       title,
			DemandScore: 8.0,
			Status:      "custom",
			Persona: model.Persona{
				Role:    asA,
				Context: "Custom workspace story",
			},
			Statement: model.UserStoryStatement{
				AsA:    asA,
				IWant:  iWant,
				SoThat: soThat,
			},
			EdgeCases: edgeCases,
			Evidence:  evList,
		}

		path, err := s.store.SaveCustom(customStory)
		if err != nil {
			return createToolResponse(req.ID, fmt.Sprintf("Error saving custom story: %v", err))
		}

		return createToolResponse(req.ID, fmt.Sprintf("Story '%s' saved successfully to: %s", id, path))

	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Unknown tool: %s", params.Name),
			},
		}
	}
}

func createToolResponse(id interface{}, text string) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]string{
				{
					"type": "text",
					"text": text,
				},
			},
		},
	}
}
