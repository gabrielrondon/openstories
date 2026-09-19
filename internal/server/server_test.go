package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielrondon/openstories/internal/store"
	"github.com/gabrielrondon/openstories/stories"
)

func setupTestServer(t *testing.T) *HTTPServer {
	st := store.New("")
	if err := st.LoadFromFS(stories.EmbeddedFS, ".", false); err != nil {
		t.Fatalf("failed to load embedded stories: %v", err)
	}
	return New(st, 8080)
}

func TestStaticIndex(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.handleStatic(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("OpenStories")) {
		t.Errorf("expected HTML to contain 'OpenStories'")
	}
}

func TestAPIStories(t *testing.T) {
	srv := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/stories?industry=devtools", nil)
	w := httptest.NewRecorder()
	srv.handleStories(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var storiesList []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &storiesList); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if len(storiesList) == 0 {
		t.Fatalf("expected at least 1 devtools story, got 0")
	}
}

func TestAPIEval(t *testing.T) {
	srv := setupTestServer(t)

	payload := map[string]string{
		"spec":     "POST /charge endpoint with stripe retries",
		"industry": "devtools",
	}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/eval", bytes.NewReader(data))
	w := httptest.NewRecorder()
	srv.handleEval(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if _, ok := result["score"]; !ok {
		t.Errorf("expected 'score' field in result")
	}
}
