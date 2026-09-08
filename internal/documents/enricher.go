package documents

import (
	"encoding/json"
	"strings"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// Enricher attaches structured legal metadata to document chunks
type Enricher struct{}

// NewEnricher creates a new metadata enricher
func NewEnricher() *Enricher {
	return &Enricher{}
}

// EnrichChunk analyzes chunk text and section reference to attach actionable legal tags
func (e *Enricher) EnrichChunk(doc *ParsedDocument, chunk *models.DocumentChunk) {
	meta := map[string]any{
		"category":       doc.CategoryID,
		"jurisdiction":   doc.JurisdictionID,
		"official_code":  doc.OfficialCode,
		"authority":      doc.Authority,
		"section_ref":    chunk.SectionReference,
	}

	lowerText := strings.ToLower(chunk.ChunkText)
	lowerSec := strings.ToLower(chunk.SectionReference)

	// Flag statutory patent bars
	if strings.Contains(lowerSec, "section 3(p)") || strings.Contains(lowerText, "traditional knowledge") {
		meta["statutory_bar"] = "Section 3(p) Traditional Knowledge Exclusion"
		meta["is_patent_bar"] = true
	}
	if strings.Contains(lowerSec, "section 3(e)") || strings.Contains(lowerText, "mere admixture") {
		meta["statutory_bar"] = "Section 3(e) Mere Admixture Exclusion"
		meta["is_patent_bar"] = true
	}
	if strings.Contains(lowerSec, "section 3(d)") || strings.Contains(lowerText, "known efficacy") {
		meta["statutory_bar"] = "Section 3(d) Enhanced Efficacy Requirement"
		meta["is_patent_bar"] = true
	}

	// Flag NBA Access and Benefit Sharing mandates
	if strings.Contains(lowerSec, "section 6") || strings.Contains(lowerText, "national biodiversity authority") {
		meta["requires_nba_approval"] = true
	}
	if strings.Contains(lowerSec, "section 7") && strings.Contains(lowerText, "vaids and hakims") {
		meta["abs_exemption"] = "Local Vaidyas & Hakims Exemption"
	}

	rawMeta, err := json.Marshal(meta)
	if err == nil {
		chunk.Metadata = rawMeta
	}
}
