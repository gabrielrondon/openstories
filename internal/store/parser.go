package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gabrielrondon/openstories/internal/model"
	"gopkg.in/yaml.v2"
)

// ParseStory parses a markdown file with YAML frontmatter into a Story struct.
func ParseStory(data []byte, filePath string) (*model.Story, error) {
	content := string(data)
	trimmed := strings.TrimSpace(content)

	if !strings.HasPrefix(trimmed, "---") {
		return nil, fmt.Errorf("file %s does not start with YAML frontmatter delimiter (---)", filePath)
	}

	// Remove leading ---
	rest := trimmed[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		// Try with \r\n---
		idx = strings.Index(rest, "\r\n---")
		if idx == -1 {
			return nil, fmt.Errorf("file %s is missing closing YAML frontmatter delimiter (---)", filePath)
		}
	}

	frontmatter := rest[:idx]
	body := strings.TrimSpace(rest[idx+4:])
	if strings.HasPrefix(body, "---") {
		body = strings.TrimSpace(body[3:])
	}

	var story model.Story
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(frontmatter)))
	if err := decoder.Decode(&story); err != nil {
		return nil, fmt.Errorf("failed to decode YAML frontmatter in %s: %w", filePath, err)
	}

	story.Body = body
	story.FilePath = filePath

	if story.Status == "" {
		story.Status = "verified"
	}
	if story.Locale == "" {
		story.Locale = "en"
	}

	return &story, nil
}

// FormatStory serializes a Story back into Markdown with YAML frontmatter.
func FormatStory(story *model.Story) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")

	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(story); err != nil {
		return nil, fmt.Errorf("failed to encode story to YAML: %w", err)
	}

	buf.WriteString("---\n\n")
	if story.Body != "" {
		buf.WriteString(story.Body)
		buf.WriteString("\n")
	} else {
		buf.WriteString(fmt.Sprintf("# %s\n\n", story.Title))
		buf.WriteString(fmt.Sprintf("> **%s** (%s)\n\n", story.Persona.Role, story.Persona.Context))
		buf.WriteString(fmt.Sprintf("**As a** %s, **I want** %s, **so that** %s.\n",
			story.Statement.AsA, story.Statement.IWant, story.Statement.SoThat))
	}

	return buf.Bytes(), nil
}
