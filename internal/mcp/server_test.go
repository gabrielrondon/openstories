package mcp

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gabrielrondon/openstories/internal/store"
	"github.com/gabrielrondon/openstories/stories"
)

func TestMCPServerStdio(t *testing.T) {
	st := store.New("")
	if err := st.LoadFromFS(stories.EmbeddedFS, ".", false); err != nil {
		t.Fatalf("failed to load embedded stories: %v", err)
	}

	srv := NewServer(st)

	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"search_stories","arguments":{"query":"idempotency"}}}` + "\n"

	inBuf := bytes.NewBufferString(input)
	outBuf := &bytes.Buffer{}

	if err := srv.Run(inBuf, outBuf); err != nil {
		t.Fatalf("srv.Run failed: %v", err)
	}

	output := outBuf.String()
	lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))

	if len(lines) < 3 {
		t.Fatalf("expected at least 3 JSON-RPC responses, got %d. Output:\n%s", len(lines), output)
	}

	// Verify initialize response
	var initResp JSONRPCResponse
	if err := json.Unmarshal(lines[0], &initResp); err != nil {
		t.Fatalf("invalid init response: %v", err)
	}
	if initResp.Error != nil {
		t.Errorf("init error: %v", initResp.Error)
	}

	// Verify tools/list response
	var listResp JSONRPCResponse
	if err := json.Unmarshal(lines[1], &listResp); err != nil {
		t.Fatalf("invalid tools/list response: %v", err)
	}
	if listResp.Error != nil {
		t.Errorf("tools/list error: %v", listResp.Error)
	}

	// Verify tools/call response
	var callResp JSONRPCResponse
	if err := json.Unmarshal(lines[2], &callResp); err != nil {
		t.Fatalf("invalid tools/call response: %v", err)
	}
	if callResp.Error != nil {
		t.Errorf("tools/call error: %v", callResp.Error)
	}
}
