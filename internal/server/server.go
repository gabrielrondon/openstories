package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gabrielrondon/openstories/internal/evaluator"
	"github.com/gabrielrondon/openstories/internal/store"
	"github.com/gabrielrondon/openstories/web"
)

// HTTPServer provides the web interface and REST API for OpenStories.
type HTTPServer struct {
	store     *store.Store
	evaluator *evaluator.Evaluator
	port      int
}

// New creates a new HTTPServer.
func New(s *store.Store, port int) *HTTPServer {
	return &HTTPServer{
		store:     s,
		evaluator: evaluator.New(s),
		port:      port,
	}
}

// Start launches the web server and blocks until termination.
func (srv *HTTPServer) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/stories", srv.handleStories)
	mux.HandleFunc("/api/stories/", srv.handleSingleStory)
	mux.HandleFunc("/api/taxonomies", srv.handleTaxonomies)
	mux.HandleFunc("/api/eval", srv.handleEval)

	// Web UI static files from embed.FS
	mux.HandleFunc("/", srv.handleStatic)

	addr := fmt.Sprintf(":%d", srv.port)
	log.Printf("🌐 OpenStories Web Interface rodando em http://localhost:%d\n", srv.port)
	log.Printf("🚀 Pronto para servir localmente ou via openstories.tuturama.com!\n")

	return http.ListenAndServe(addr, mux)
}

func (srv *HTTPServer) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	data, err := web.EmbeddedFS.ReadFile(path)
	if err != nil {
		// Fallback to index.html for SPA-style routing
		data, err = web.EmbeddedFS.ReadFile("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	if strings.HasSuffix(path, ".html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	}

	w.Write(data)
}

func (srv *HTTPServer) handleStories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	industry := r.URL.Query().Get("industry")
	domain := r.URL.Query().Get("domain")
	query := r.URL.Query().Get("q")

	results := srv.store.Search(query, industry, domain, 0)
	stories := make([]interface{}, 0, len(results))
	for _, res := range results {
		stories = append(stories, res.Story)
	}

	json.NewEncoder(w).Encode(stories)
}

func (srv *HTTPServer) handleSingleStory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := strings.TrimPrefix(r.URL.Path, "/api/stories/")
	st, ok := srv.store.Get(id)
	if !ok {
		http.Error(w, `{"error": "Story not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(st)
}

func (srv *HTTPServer) handleTaxonomies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	taxonomies := srv.store.ListTaxonomies()
	json.NewEncoder(w).Encode(taxonomies)
}

func (srv *HTTPServer) handleEval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to read body"}`, http.StatusBadRequest)
		return
	}

	var payload struct {
		Spec     string `json:"spec"`
		Industry string `json:"industry"`
		Domain   string `json:"domain"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result := srv.evaluator.Evaluate(payload.Spec, payload.Industry, payload.Domain)
	json.NewEncoder(w).Encode(result)
}
