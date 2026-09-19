package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAllCuratedStories(t *testing.T) {
	storiesDir := "../../stories"
	err := filepath.Walk(storiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || (!strings.HasSuffix(path, ".story.md") && !strings.HasSuffix(path, ".md")) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", path, err)
		}
		_, err = ParseStory(data, path)
		if err != nil {
			t.Errorf("Failed to parse %s: %v", path, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
