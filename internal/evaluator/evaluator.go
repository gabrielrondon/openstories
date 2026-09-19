package evaluator

import (
	"fmt"
	"strings"

	"github.com/gabrielrondon/openstories/internal/model"
	"github.com/gabrielrondon/openstories/internal/store"
)

// Evaluator checks technical specifications or PR descriptions against evidence-backed stories.
type Evaluator struct {
	store *store.Store
}

// New creates an Evaluator backed by a store.
func New(s *store.Store) *Evaluator {
	return &Evaluator{store: s}
}

// Evaluate analyzes the provided spec text against relevant stories in the domain/industry.
func (e *Evaluator) Evaluate(specText string, industry string, domain string) *model.SpecEvaluationResult {
	// Search for relevant stories matching the spec
	matches := e.store.Search(specText, industry, domain, 0)
	if len(matches) == 0 {
		// Fallback to searching by industry if specified
		if industry != "" {
			matches = e.store.Search("", industry, domain, 0)
		}
	}

	result := &model.SpecEvaluationResult{
		TargetIndustry:  industry,
		MatchedStories:  make([]model.MatchedStoryCheck, 0),
		CriticalGaps:    make([]string, 0),
		EdgeCaseAlerts:  make([]string, 0),
		Recommendations: make([]string, 0),
	}

	specLower := strings.ToLower(specText)
	totalEdgeCases := 0
	coveredEdgeCases := 0

	// Consider top 5 most relevant stories
	limit := 5
	if len(matches) < limit {
		limit = len(matches)
	}

	for i := 0; i < limit; i++ {
		st := matches[i].Story
		check := model.MatchedStoryCheck{
			StoryID:     st.ID,
			Title:       st.Title,
			DemandScore: st.DemandScore,
			MissedCases: make([]string, 0),
		}

		storyCovered := true
		for _, ec := range st.EdgeCases {
			totalEdgeCases++
			tokens := extractKeyPhrases(ec)
			matched := false
			for _, tok := range tokens {
				if strings.Contains(specLower, tok) {
					matched = true
					break
				}
			}

			if matched {
				coveredEdgeCases++
			} else {
				storyCovered = false
				check.MissedCases = append(check.MissedCases, ec)
				result.EdgeCaseAlerts = append(result.EdgeCaseAlerts,
					fmt.Sprintf("[%s] %s: %s", st.ID, st.Title, ec))
			}
		}

		check.Covered = storyCovered
		result.MatchedStories = append(result.MatchedStories, check)

		// Check evaluation rubrics
		for _, rubric := range st.EvaluationRubric {
			tokens := extractKeyPhrases(rubric)
			found := false
			for _, tok := range tokens {
				if strings.Contains(specLower, tok) {
					found = true
					break
				}
			}
			if !found {
				result.CriticalGaps = append(result.CriticalGaps,
					fmt.Sprintf("[%s] Checklist não atendido: %s", st.ID, rubric))
			}
		}
	}

	// Calculate realistic score
	if totalEdgeCases == 0 {
		result.Score = 75
		result.Summary = "Nenhuma história específica com casos de borda encontrada para comparação detalhada. A especificação parece aceitável mas sem validação empírica."
	} else {
		scorePercent := int((float64(coveredEdgeCases) / float64(totalEdgeCases)) * 100)
		result.Score = scorePercent

		if scorePercent >= 85 {
			result.Summary = "Excelente cobertura! A especificação antecipou os principais pontos de falha reportados por usuários em produção."
		} else if scorePercent >= 60 {
			result.Summary = "Cobertura moderada. A especificação trata o fluxo principal, mas possui pontos cegos críticos observados em incidentes reais de produção."
		} else {
			result.Summary = "Atenção: Alta vulnerabilidade a incidentes de produção. A especificação negligenciou múltiplos casos de borda severos com alto volume de reclamações reais na internet."
		}
	}

	// Generate actionable recommendations
	if len(result.EdgeCaseAlerts) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Adicionar testes automatizados específicos para os casos de borda listados nos alertas.")
	}
	if len(result.CriticalGaps) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Incorporar na documentação do PR ou RFC as respostas aos itens do checklist não atendidos.")
	}
	result.Recommendations = append(result.Recommendations,
		"Utilizar os critérios de aceitação no formato Gherkin das histórias mapeadas como suíte de testes E2E.")

	return result
}

func extractKeyPhrases(s string) []string {
	clean := strings.ToLower(s)
	words := strings.FieldsFunc(clean, func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'))
	})
	var phrases []string
	for _, w := range words {
		if len(w) > 4 {
			phrases = append(phrases, w)
		}
	}
	return phrases
}
