package llm

import (
	"context"
	"strings"
	"testing"
)

func TestMockClientGenerationAndCitations(t *testing.T) {
	client := NewMockClient()

	ctx := context.Background()
	prompt := BuildUserPrompt("Can I patent an Ayurvedic formula with Ashwagandha and Turmeric?", "IN", "PATENT", "en", []string{"Section 3(p) Traditional knowledge bar"})

	ans, err := client.Generate(ctx, SystemPromptAyurvedaIP, prompt)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	if !strings.Contains(ans, "Section 3(p)") {
		t.Errorf("expected Section 3(p) in mock response, got: %s", ans)
	}

	if !strings.Contains(ans, "[^1]") {
		t.Errorf("expected footnote marker [^1] in response, got: %s", ans)
	}
}

func TestMockClientEmbeddings(t *testing.T) {
	client := NewMockClient()

	ctx := context.Background()
	vec, err := client.Embed(ctx, "The Patents Act, 1970")
	if err != nil {
		t.Fatalf("embed failed: %v", err)
	}

	if len(vec) != 1536 {
		t.Fatalf("expected 1536 dimensions, got %d", len(vec))
	}

	// Verify deterministic property: same input produces exact same vector
	vec2, _ := client.Embed(ctx, "The Patents Act, 1970")
	for i := range vec {
		if vec[i] != vec2[i] {
			t.Fatalf("embeddings are not deterministic at index %d", i)
		}
	}
}
