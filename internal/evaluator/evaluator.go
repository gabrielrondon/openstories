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
	matches := e.store.Search(specText, industry, domain, 0)
	if len(matches) == 0 {
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
					fmt.Sprintf("[%s] Unaddressed Quality Rubric: %s", st.ID, rubric))
			}
		}
	}

	if totalEdgeCases == 0 {
		result.Score = 75
		result.Summary = "No domain-specific stories with failure criteria found for detailed comparison. Specification appears reasonable but lacks empirical validation."
	} else {
		scorePercent := int((float64(coveredEdgeCases) / float64(totalEdgeCases)) * 100)
		result.Score = scorePercent

		if scorePercent >= 85 {
			result.Summary = "High resilience! The technical specification preemptively mitigates known production failure modes documented in field reports."
		} else if scorePercent >= 60 {
			result.Summary = "Moderate coverage. The happy path is addressed, but critical production blind spots frequently reported by users in the wild are missing."
		} else {
			result.Summary = "Severe Production Vulnerability Alert: Multiple high-impact failure modes observed in real-world outages were ignored in this specification."
		}
	}

	if len(result.EdgeCaseAlerts) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Add explicit integration and failure-injection tests covering the flagged production edge cases.")
	}
	if len(result.CriticalGaps) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Document explicit architectural answers to the quality checklist rubrics before writing code.")
	}
	result.Recommendations = append(result.Recommendations,
		"Utilize the Gherkin acceptance criteria from matched user stories as an automated E2E test suite.")

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
