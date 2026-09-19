package harvester

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabrielrondon/openstories/internal/model"
	"github.com/gabrielrondon/openstories/internal/store"
)

// GitHubIssue represents the subset of GitHub's API issue response.
type GitHubIssue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Harvester handles data collection from public developer communities.
type Harvester struct {
	store  *store.Store
	client *http.Client
}

// New creates a new Harvester instance.
func New(s *store.Store) *Harvester {
	return &Harvester{
		store: s,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// HarvestGitHub fetches issues from a repository and synthesizes them into structured stories.
func (h *Harvester) HarvestGitHub(repo string, label string, industry string, domain string) (*model.Story, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/issues?state=all&per_page=5", repo)
	if label != "" {
		apiURL += fmt.Sprintf("&labels=%s", label)
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "OpenStories-Harvester/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to contact GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var issues []GitHubIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub response: %w", err)
	}

	if len(issues) == 0 {
		return nil, fmt.Errorf("no issues found in %s", repo)
	}

	primary := issues[0]
	storyID := fmt.Sprintf("OS-HARVEST-%d", primary.Number)

	quote := primary.Body
	if len(quote) > 280 {
		quote = quote[:277] + "..."
	}

	evidenceList := make([]model.Evidence, 0)
	evidenceList = append(evidenceList, model.Evidence{
		Source:   primary.HTMLURL,
		Type:     "github_issue",
		Quote:    quote,
		Date:     primary.CreatedAt.Format("2006-01-02"),
		Platform: "github",
	})

	synthesized := &model.Story{
		ID:          storyID,
		Industry:    industry,
		Domain:      domain,
		Title:       primary.Title,
		DemandScore: 8.5,
		Status:      "community",
		Persona: model.Persona{
			Role:    "Engenheiro / Usuário de " + repo,
			Context: fmt.Sprintf("Reportado originalmente na issue #%d de %s", primary.Number, repo),
		},
		Statement: model.UserStoryStatement{
			AsA:    "Desenvolvedor utilizando a biblioteca " + repo,
			IWant:  "Tratamento robusto para: " + primary.Title,
			SoThat: "A aplicação não falhe em cenários inesperados reportados em produção",
		},
		AcceptanceCriteria: []model.AcceptanceCriterion{
			{
				Scenario: "Cenário reportado na issue",
				Given:    "A biblioteca " + repo + " configurada no ambiente",
				When:     "O fluxo principal da issue for executado",
				Then:     "O comportamento deve seguir a especificação sem lançar exceções não tratadas",
			},
		},
		EdgeCases: []string{
			"Falha transitória reportada na issue original",
		},
		Evidence: evidenceList,
		Tags:     []string{"harvested", "github", repo},
	}

	path, err := h.store.SaveCustom(synthesized)
	if err != nil {
		return nil, err
	}
	fmt.Printf("🎉 História minerada com sucesso e salva em: %s\n", path)

	return synthesized, nil
}
