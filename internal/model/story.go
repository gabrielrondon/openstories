package model

import "time"

// Story represents an evidence-backed user story.
type Story struct {
	ID                 string                `yaml:"id" json:"id"`
	Industry           string                `yaml:"industry" json:"industry"`
	Domain             string                `yaml:"domain" json:"domain"`
	Title              string                `yaml:"title" json:"title"`
	DemandScore        float64               `yaml:"demand_score" json:"demand_score"`
	Status             string                `yaml:"status" json:"status"` // verified, emerging, community
	Persona            Persona               `yaml:"persona" json:"persona"`
	Statement          UserStoryStatement   `yaml:"story" json:"story"`
	AcceptanceCriteria []AcceptanceCriterion `yaml:"acceptance_criteria" json:"acceptance_criteria"`
	EdgeCases          []string              `yaml:"edge_cases" json:"edge_cases"`
	Evidence           []Evidence            `yaml:"evidence" json:"evidence"`
	EvaluationRubric   []string              `yaml:"evaluation_rubric" json:"evaluation_rubric"`
	Tags               []string              `yaml:"tags" json:"tags"`
	Body               string                `yaml:"-" json:"body,omitempty"`
	FilePath           string                `yaml:"-" json:"file_path,omitempty"`
	IsCustom           bool                  `yaml:"-" json:"is_custom"`
}

// Persona describes who is experiencing this need.
type Persona struct {
	Role    string `yaml:"role" json:"role"`
	Context string `yaml:"context" json:"context"`
}

// UserStoryStatement follows the classic As a / I want / So that pattern.
type UserStoryStatement struct {
	AsA   string `yaml:"as_a" json:"as_a"`
	IWant string `yaml:"i_want" json:"i_want"`
	SoThat string `yaml:"so_that" json:"so_that"`
}

// AcceptanceCriterion is a Gherkin-style testable requirement.
type AcceptanceCriterion struct {
	Scenario string `yaml:"scenario" json:"scenario"`
	Given    string `yaml:"given" json:"given"`
	When     string `yaml:"when" json:"when"`
	Then     string `yaml:"then" json:"then"`
}

// Evidence links the story to verifiable real-world reports, complaints, or post-mortems.
type Evidence struct {
	Source    string    `yaml:"source" json:"source"`
	Type      string    `yaml:"type" json:"type"` // github_issue, reddit_thread, hackernews, post_mortem, app_review
	Quote     string    `yaml:"quote" json:"quote"`
	Date      string    `yaml:"date" json:"date"`
	Platform  string    `yaml:"platform" json:"platform,omitempty"`
	Upvotes   int       `yaml:"upvotes,omitempty" json:"upvotes,omitempty"`
	Timestamp time.Time `yaml:"-" json:"-"`
}

// SpecEvaluationResult contains the output of evaluating a technical spec against real stories.
type SpecEvaluationResult struct {
	TargetIndustry string              `json:"target_industry"`
	Score          int                 `json:"score"` // 0 - 100
	Summary        string              `json:"summary"`
	MatchedStories []MatchedStoryCheck `json:"matched_stories"`
	CriticalGaps   []string            `json:"critical_gaps"`
	EdgeCaseAlerts []string            `json:"edge_case_alerts"`
	Recommendations []string           `json:"recommendations"`
}

// MatchedStoryCheck evaluates coverage for a single story.
type MatchedStoryCheck struct {
	StoryID     string   `json:"story_id"`
	Title       string   `json:"title"`
	DemandScore float64  `json:"demand_score"`
	Covered     bool     `json:"covered"`
	MissedCases []string `json:"missed_cases"`
}

// IndustryTaxonomy maps industries to their domains.
type IndustryTaxonomy struct {
	Industry string   `json:"industry"`
	Label    string   `json:"label"`
	Domains  []string `json:"domains"`
}
