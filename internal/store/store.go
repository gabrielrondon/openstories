package store

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gabrielrondon/openstories/internal/model"
)

// Store manages the in-memory index of user stories.
type Store struct {
	mu           sync.RWMutex
	stories      map[string]*model.Story
	allStories   []*model.Story
	customDir    string
}

// New creates an empty Store.
func New(customDir string) *Store {
	return &Store{
		stories:    make(map[string]*model.Story),
		allStories: make([]*model.Story, 0),
		customDir:  customDir,
	}
}

// LoadFromFS recursively walks an fs.FS and indexes all .story.md or .md files.
func (s *Store) LoadFromFS(fileSys fs.FS, root string, isCustom bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return fs.WalkDir(fileSys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".story.md") && !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := fs.ReadFile(fileSys, path)
		if err != nil {
			return err
		}

		story, err := ParseStory(data, path)
		if err != nil {
			// Skip files that aren't valid stories (e.g. regular markdown docs)
			return nil
		}

		story.IsCustom = isCustom
		s.stories[story.ID] = story
		s.rebuildSlice()
		return nil
	})
}

// LoadLocalDir loads custom stories from a local filesystem directory.
func (s *Store) LoadLocalDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil
	}
	return s.LoadFromFS(os.DirFS(dirPath), ".", true)
}

func (s *Store) rebuildSlice() {
	s.allStories = make([]*model.Story, 0, len(s.stories))
	for _, story := range s.stories {
		s.allStories = append(s.allStories, story)
	}
	// Sort by demand score descending by default
	sort.Slice(s.allStories, func(i, j int) bool {
		return s.allStories[i].DemandScore > s.allStories[j].DemandScore
	})
}

// Get retrieves a story by ID.
func (s *Store) Get(id string) (*model.Story, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	story, ok := s.stories[strings.ToUpper(strings.TrimSpace(id))]
	if !ok {
		// Try case-insensitive lookup
		for _, st := range s.stories {
			if strings.EqualFold(st.ID, id) {
				return st, true
			}
		}
		return nil, false
	}
	return story, true
}

// SearchResult pairs a story with its calculated relevance score.
type SearchResult struct {
	Story *model.Story `json:"story"`
	Score float64      `json:"score"`
}

// Search queries stories by query terms, industry, domain and minimum demand score.
func (s *Store) Search(query string, industry string, domain string, minScore float64) []SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tokens := tokenize(query)
	results := make([]SearchResult, 0)

	for _, story := range s.allStories {
		if industry != "" && !strings.EqualFold(story.Industry, industry) {
			continue
		}
		if domain != "" && !strings.EqualFold(story.Domain, domain) {
			continue
		}
		if minScore > 0 && story.DemandScore < minScore {
			continue
		}

		if len(tokens) == 0 {
			// No query string provided: return all matching the filters, scored by demand
			results = append(results, SearchResult{
				Story: story,
				Score: story.DemandScore,
			})
			continue
		}

		matchScore := calculateScore(story, tokens)
		if matchScore > 0 {
			results = append(results, SearchResult{
				Story: story,
				Score: matchScore,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// ListAll returns all indexed stories.
func (s *Store) ListAll() []*model.Story {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*model.Story, len(s.allStories))
	copy(out, s.allStories)
	return out
}

// ListTaxonomies returns all unique industries and domains currently registered.
func (s *Store) ListTaxonomies() []model.IndustryTaxonomy {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tree := make(map[string]map[string]bool)
	for _, story := range s.allStories {
		ind := strings.ToLower(story.Industry)
		dom := strings.ToLower(story.Domain)
		if ind == "" {
			ind = "general"
		}
		if _, ok := tree[ind]; !ok {
			tree[ind] = make(map[string]bool)
		}
		if dom != "" {
			tree[ind][dom] = true
		}
	}

	taxonomies := make([]model.IndustryTaxonomy, 0, len(tree))
	for ind, doms := range tree {
		domList := make([]string, 0, len(doms))
		for d := range doms {
			domList = append(domList, d)
		}
		sort.Strings(domList)

		taxonomies = append(taxonomies, model.IndustryTaxonomy{
			Industry: ind,
			Label:    formatLabel(ind),
			Domains:  domList,
		})
	}

	sort.Slice(taxonomies, func(i, j int) bool {
		return taxonomies[i].Industry < taxonomies[j].Industry
	})

	return taxonomies
}

// SaveCustom writes a new story to the local custom directory (.openstories/stories/...).
func (s *Store) SaveCustom(story *model.Story) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.customDir == "" {
		s.customDir = ".openstories/stories"
	}

	targetDir := filepath.Join(s.customDir, story.Industry, story.Domain)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	slug := sanitizeSlug(story.Title)
	if slug == "" {
		slug = strings.ToLower(story.ID)
	}
	fileName := fmt.Sprintf("%s.story.md", slug)
	fullPath := filepath.Join(targetDir, fileName)

	data, err := FormatStory(story)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write story file: %w", err)
	}

	story.FilePath = fullPath
	story.IsCustom = true
	s.stories[story.ID] = story
	s.rebuildSlice()

	return fullPath, nil
}

func tokenize(s string) []string {
	clean := strings.ToLower(s)
	f := func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c >= 'à' && c <= 'ú'))
	}
	words := strings.FieldsFunc(clean, f)
	var filtered []string
	stopWords := map[string]bool{
		"a": true, "o": true, "de": true, "do": true, "da": true, "em": true, "um": true, "uma": true,
		"para": true, "com": true, "e": true, "ou": true, "the": true, "and": true, "of": true, "to": true, "in": true,
		"for": true, "is": true, "on": true, "that": true, "this": true, "as": true, "i": true, "want": true, "so": true,
	}
	for _, w := range words {
		if len(w) > 1 && !stopWords[w] {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

func calculateScore(story *model.Story, tokens []string) float64 {
	score := 0.0
	titleLower := strings.ToLower(story.Title)
	asALower := strings.ToLower(story.Statement.AsA)
	iWantLower := strings.ToLower(story.Statement.IWant)
	soThatLower := strings.ToLower(story.Statement.SoThat)

	for _, t := range tokens {
		if strings.Contains(titleLower, t) {
			score += 25.0
		}
		if strings.Contains(iWantLower, t) {
			score += 15.0
		}
		if strings.Contains(soThatLower, t) || strings.Contains(asALower, t) {
			score += 10.0
		}
		for _, tag := range story.Tags {
			if strings.EqualFold(tag, t) {
				score += 15.0
			}
		}
		for _, ec := range story.EdgeCases {
			if strings.Contains(strings.ToLower(ec), t) {
				score += 8.0
			}
		}
		for _, ev := range story.Evidence {
			if strings.Contains(strings.ToLower(ev.Quote), t) {
				score += 6.0
			}
		}
		for _, ac := range story.AcceptanceCriteria {
			text := strings.ToLower(ac.Scenario + " " + ac.Given + " " + ac.When + " " + ac.Then)
			if strings.Contains(text, t) {
				score += 6.0
			}
		}
	}

	if score > 0 {
		// Weight by demand score
		multiplier := 1.0 + (story.DemandScore / 10.0)
		score = score * multiplier
	}

	return score
}

func formatLabel(key string) string {
	parts := strings.Split(key, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

func sanitizeSlug(title string) string {
	clean := strings.ToLower(title)
	clean = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		if r == ' ' || r == '_' {
			return '-'
		}
		return -1
	}, clean)
	for strings.Contains(clean, "--") {
		clean = strings.ReplaceAll(clean, "--", "-")
	}
	return strings.Trim(clean, "-")
}
