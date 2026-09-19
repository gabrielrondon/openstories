package store

import (
	"testing"
)

const sampleStory = `---
id: OS-TEST-001
industry: fintech
domain: payments
title: Test Payment Story
demand_score: 9.5
status: verified
persona:
  role: Payment Engineer
  context: High volume SaaS
story:
  as_a: Payment Engineer
  i_want: Idempotency keys on charge requests
  so_that: Network retries do not double charge customers
acceptance_criteria:
  - scenario: Replay same key
    given: A charge with key K1 succeeded
    when: Key K1 is posted again
    then: Return original 200 response
edge_cases:
  - Expired keys after 24h
  - Concurrent requests with same key
evidence:
  - source: https://github.com/test/repo/issues/1
    type: github_issue
    quote: We had duplicate transactions during a network blip.
    date: 2025-05-10
tags:
  - payments
  - idempotency
---

# Details
More details here.
`

func TestParseStory(t *testing.T) {
	st, err := ParseStory([]byte(sampleStory), "test.story.md")
	if err != nil {
		t.Fatalf("ParseStory failed: %v", err)
	}

	if st.ID != "OS-TEST-001" {
		t.Errorf("expected ID OS-TEST-001, got %s", st.ID)
	}
	if st.DemandScore != 9.5 {
		t.Errorf("expected DemandScore 9.5, got %f", st.DemandScore)
	}
	if len(st.AcceptanceCriteria) != 1 {
		t.Errorf("expected 1 acceptance criterion, got %d", len(st.AcceptanceCriteria))
	}
	if len(st.EdgeCases) != 2 {
		t.Errorf("expected 2 edge cases, got %d", len(st.EdgeCases))
	}
	if len(st.Evidence) != 1 {
		t.Errorf("expected 1 evidence item, got %d", len(st.Evidence))
	}
	if st.Body != "# Details\nMore details here." {
		t.Errorf("unexpected body: %s", st.Body)
	}
}

func TestStoreSearch(t *testing.T) {
	st, err := ParseStory([]byte(sampleStory), "test.story.md")
	if err != nil {
		t.Fatal(err)
	}

	s := New("")
	s.stories[st.ID] = st
	s.rebuildSlice()

	// Exact word match
	results := s.Search("idempotency", "", "", 0)
	if len(results) == 0 {
		t.Fatal("expected at least 1 result for search 'idempotency'")
	}
	if results[0].Story.ID != "OS-TEST-001" {
		t.Errorf("expected OS-TEST-001, got %s", results[0].Story.ID)
	}

	// Filter by industry
	fintechResults := s.Search("", "fintech", "", 0)
	if len(fintechResults) != 1 {
		t.Fatalf("expected 1 result for fintech filter, got %d", len(fintechResults))
	}

	// Filter by mismatching industry
	healthResults := s.Search("", "healthcare", "", 0)
	if len(healthResults) != 0 {
		t.Fatalf("expected 0 results for healthcare filter, got %d", len(healthResults))
	}
}
