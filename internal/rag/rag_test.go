package rag

import (
	"context"
	"strings"
	"testing"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/llm"
)

func TestRetrieverFallback(t *testing.T) {
	mockLLM := llm.NewMockClient()
	retriever := NewRetriever(nil, mockLLM) // nil DB triggers graceful fallback

	ctx := context.Background()
	results, err := retriever.Retrieve(ctx, "How does Section 3(p) affect patenting classical formulations?", "IN", "PATENT", 2)
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least 1 fallback result, got 0")
	}

	found3p := false
	for _, r := range results {
		if strings.Contains(r.SectionReference, "3(p)") {
			found3p = true
		}
	}

	if !found3p {
		t.Errorf("expected Section 3(p) in retrieved chunks, got: %+v", results)
	}
}

func TestContextAssembler(t *testing.T) {
	mockLLM := llm.NewMockClient()
	retriever := NewRetriever(nil, mockLLM)
	assembler := NewContextAssembler()

	results, _ := retriever.Retrieve(context.Background(), "patents and traditional knowledge", "IN", "PATENT", 2)
	blocks := assembler.BuildContextBlocks(results)

	if len(blocks) != len(results) {
		t.Fatalf("expected %d context blocks, got %d", len(results), len(blocks))
	}

	if !strings.Contains(blocks[0], "CITABLE SOURCE") {
		t.Errorf("expected CITABLE SOURCE prefix in block, got: %s", blocks[0])
	}
}
