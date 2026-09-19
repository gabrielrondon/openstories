package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents standard MCP client configuration.
type Config struct {
	MCPServers map[string]MCPServerEntry `json:"mcpServers"`
}

type MCPServerEntry struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
	Type    string            `json:"type,omitempty"`
}

// InstallClaude configures OpenStories for both Claude Desktop and Claude Code CLI.
func InstallClaude() error {
	var errors []error

	// 1. Configure Claude Code CLI (~/.claude.json)
	if err := InstallClaudeCode(); err != nil {
		errors = append(errors, fmt.Errorf("Claude Code CLI: %w", err))
	}

	// 2. Configure Claude Desktop (GUI app)
	if err := InstallClaudeDesktop(); err != nil {
		// Non-fatal if user only uses Claude Code CLI
		fmt.Printf("Note: Claude Desktop app config skipped or not found: %v\n", err)
	}

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// InstallClaudeCode configures OpenStories inside ~/.claude.json used by Claude Code CLI.
func InstallClaudeCode() error {
	binaryPath, err := getBinaryPath()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".claude.json")
	return injectIntoClaudeCodeJSON(configPath, "openstories", binaryPath, []string{"mcp"})
}

// InstallClaudeDesktop configures OpenStories inside Claude Desktop's app config.
func InstallClaudeDesktop() error {
	binaryPath, err := getBinaryPath()
	if err != nil {
		return err
	}

	configPath, err := getClaudeDesktopConfigPath()
	if err != nil {
		return err
	}

	return injectMCPServer(configPath, "openstories", binaryPath, []string{"mcp"})
}

// InstallCursor configures OpenStories inside Cursor's configuration.
func InstallCursor() error {
	binaryPath, err := getBinaryPath()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".cursor", "mcp.json")
	return injectMCPServer(configPath, "openstories", binaryPath, []string{"mcp"})
}

func getBinaryPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "openstories", nil
	}
	// Return absolute path so Claude Code or Desktop can invoke it from any directory
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return execPath, nil
	}
	return absPath, nil
}

func getClaudeDesktopConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "Claude", "claude_desktop_config.json"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Claude", "claude_desktop_config.json"), nil
	default:
		return filepath.Join(homeDir, ".config", "Claude", "claude_desktop_config.json"), nil
	}
}

func injectIntoClaudeCodeJSON(configPath, serverName, command string, args []string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("failed to parse %s: %w", configPath, err)
	}

	var mcpServers map[string]interface{}
	if rawServers, ok := root["mcpServers"].(map[string]interface{}); ok {
		mcpServers = rawServers
	} else {
		mcpServers = make(map[string]interface{})
		root["mcpServers"] = mcpServers
	}

	mcpServers[serverName] = map[string]interface{}{
		"command": command,
		"args":    args,
		"env":     map[string]string{},
		"type":    "stdio",
	}

	updated, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if err := os.WriteFile(configPath, updated, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", configPath, err)
	}

	fmt.Printf("✅ Configured OpenStories MCP in Claude Code CLI: %s\n", configPath)
	return nil
}

func injectMCPServer(configPath, serverName, command string, args []string) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	cfg := Config{
		MCPServers: make(map[string]MCPServerEntry),
	}

	if data, err := os.ReadFile(configPath); err == nil {
		_ = json.Unmarshal(data, &cfg)
		if cfg.MCPServers == nil {
			cfg.MCPServers = make(map[string]MCPServerEntry)
		}
	}

	cfg.MCPServers[serverName] = MCPServerEntry{
		Command: command,
		Args:    args,
		Type:    "stdio",
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", configPath, err)
	}

	fmt.Printf("✅ Configured OpenStories MCP in: %s\n", configPath)
	return nil
}
