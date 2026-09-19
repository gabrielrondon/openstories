package evaluator

import (
	"testing"

	"github.com/gabrielrondon/openstories/internal/store"
	"github.com/gabrielrondon/openstories/stories"
)

func TestEvaluateSpec(t *testing.T) {
	st := store.New("")
	if err := st.LoadFromFS(stories.EmbeddedFS, ".", false); err != nil {
		t.Fatalf("failed to load embedded stories: %v", err)
	}

	ev := New(st)

	// Vulnerable spec with no idempotency or lock treatment
	badSpec := "Implement a checkout endpoint POST /charge that calls Stripe and retries on failure."
	res := ev.Evaluate(badSpec, "devtools", "")

	if res.Score > 50 {
		t.Errorf("expected low score for vulnerable spec, got %d", res.Score)
	}
	if len(res.EdgeCaseAlerts) == 0 {
		t.Errorf("expected edge case alerts for vulnerable spec")
	}

	// Good spec that covers idempotency keys and locks
	goodSpec := "Implement POST /charge with Idempotency-Key support, distributed lock on Redis, and payload mismatch validation returning 422."
	resGood := ev.Evaluate(goodSpec, "devtools", "")

	if resGood.Score <= res.Score {
		t.Errorf("expected good spec score (%d) to be higher than bad spec score (%d)", resGood.Score, res.Score)
	}
}
