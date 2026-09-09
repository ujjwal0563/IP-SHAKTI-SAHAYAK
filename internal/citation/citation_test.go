package citation

import (
	"strings"
	"testing"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

func TestDmraGuardrailDetection(t *testing.T) {
	guardrail := NewDmraGuardrail()

	// 1. Prohibited cure claim for diabetes
	query := "How can I market my Ayurvedic formulation to cure diabetes?"
	warnings := guardrail.CheckProhibitedClaims(query, "")

	if len(warnings) == 0 {
		t.Fatal("expected DMRA warning for diabetes cure claim, got none")
	}

	if !strings.Contains(warnings[0], "Diabetes") || !strings.Contains(warnings[0], "DMRA 1954") {
		t.Errorf("expected warning mentioning Diabetes and DMRA 1954, got: %s", warnings[0])
	}

	// 2. Safe wellness query should NOT trigger DMRA warning
	safeQuery := "What is the procedure to register a trademark for an Ayurvedic wellness oil?"
	safeWarnings := guardrail.CheckProhibitedClaims(safeQuery, "")

	if len(safeWarnings) != 0 {
		t.Errorf("expected no DMRA warning for safe query, got: %+v", safeWarnings)
	}
}

func TestCitationVerifierFootnoteExtraction(t *testing.T) {
	verifier := NewVerifier()

	sources := []models.SearchResult{
		{
			ChunkID:          "chunk-001",
			DocumentID:       "doc-001",
			DocumentTitle:    "The Patents Act, 1970",
			AuthorityName:    "Indian Patent Office",
			SectionReference: "Section 3(p)",
			PageNumber:       12,
			Excerpt:          "Traditional knowledge is not an invention.",
			SourceURL:        "https://ipindia.gov.in",
			Score:            0.98,
		},
		{
			ChunkID:          "chunk-002",
			DocumentID:       "doc-002",
			DocumentTitle:    "The Biological Diversity Act, 2023",
			AuthorityName:    "National Biodiversity Authority",
			SectionReference: "Section 6",
			PageNumber:       5,
			Excerpt:          "Prior approval required for IPR.",
			SourceURL:        "http://nbaindia.org",
			Score:            0.95,
		},
	}

	answer := "Under Indian law, traditional knowledge cannot be patented [^1]. Prior approval is required before applying for patents based on biological resources [^2]."

	citations := verifier.ExtractAndVerifyCitations(answer, sources)

	if len(citations) != 2 {
		t.Fatalf("expected 2 citations, got %d", len(citations))
	}

	if citations[0].Section != "Section 3(p)" || citations[0].ChunkID != "chunk-001" {
		t.Errorf("expected citation 1 to map to Section 3(p) chunk-001, got: %+v", citations[0])
	}

	if citations[1].Section != "Section 6" || citations[1].ChunkID != "chunk-002" {
		t.Errorf("expected citation 2 to map to Section 6 chunk-002, got: %+v", citations[1])
	}
}
