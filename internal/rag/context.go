package rag

import (
	"fmt"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// ContextAssembler formats retrieved search chunks into legal grounding passages
type ContextAssembler struct{}

// NewContextAssembler creates a new context builder
func NewContextAssembler() *ContextAssembler {
	return &ContextAssembler{}
}

// BuildContextBlocks creates structured passages with explicit statutory citations
func (ca *ContextAssembler) BuildContextBlocks(results []models.SearchResult) []string {
	var blocks []string
	for i, r := range results {
		block := fmt.Sprintf(
			"CITABLE SOURCE [%d]:\n- Statute: %s\n- Provision: %s\n- Authority: %s\n- Official Excerpt: \"%s\"\n- Official Link: %s",
			i+1, r.DocumentTitle, r.SectionReference, r.AuthorityName, r.Excerpt, r.SourceURL,
		)
		blocks = append(blocks, block)
	}
	return blocks
}
