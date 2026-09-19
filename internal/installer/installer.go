package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the MCP client configuration file.
type Config struct {
	MCPServers map[string]MCPServerEntry `json:"mcpServers"`
}

type MCPServerEntry struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// InstallClaude configures OpenStories inside Claude Desktop's configuration file.
func InstallClaude() error {
	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "openstories"
	}

	configPath, err := getClaudeConfigPath()
	if err != nil {
		return err
	}

	return injectMCPServer(configPath, "openstories", binaryPath, []string{"mcp"})
}

// InstallCursor configures OpenStories inside Cursor's configuration.
func InstallCursor() error {
	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "openstories"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Cursor MCP config can be placed in ~/.cursor/mcp.json
	configPath := filepath.Join(homeDir, ".cursor", "mcp.json")
	return injectMCPServer(configPath, "openstories", binaryPath, []string{"mcp"})
}

func getClaudeConfigPath() (string, error) {
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
		// Linux
		return filepath.Join(homeDir, ".config", "Claude", "claude_desktop_config.json"), nil
	}
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
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", configPath, err)
	}

	fmt.Printf("✅ Sucesso! Servidor MCP configurado em: %s\n", configPath)
	fmt.Println("🚀 Reinicie o Claude Desktop ou Cursor para que as ferramentas do OpenStories fiquem disponíveis imediatamente!")
	return nil
}
