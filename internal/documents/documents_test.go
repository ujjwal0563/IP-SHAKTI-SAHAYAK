package documents

import (
	"testing"
)

func TestStatutoryChunkerSection3p(t *testing.T) {
	doc := &ParsedDocument{
		Title:          "The Patents Act, 1970",
		OfficialCode:   "Act No. 39 of 1970",
		JurisdictionID: "IN",
		CategoryID:     "PATENT",
		RawContent: `
### Section 3 - What are not inventions
The following are not inventions within the meaning of this Act:

#### Section 3(d)
The mere discovery of a new form of a known substance which does not result in the enhancement of the known efficacy of that substance.

#### Section 3(e)
A substance obtained by a mere admixture resulting only in the aggregation of the properties of the components thereof.

#### Section 3(p)
An invention which in effect is traditional knowledge or which is an aggregation or duplication of known properties of traditionally known component or components.
`,
	}

	chunker := NewChunker()
	chunks := chunker.ChunkDocument(doc)

	if len(chunks) < 3 {
		t.Fatalf("expected at least 3 statutory chunks, got %d", len(chunks))
	}

	found3p := false
	for _, c := range chunks {
		if c.SectionReference == "Section 3(p)" {
			found3p = true
			if !contains(c.Text, "traditional knowledge") {
				t.Errorf("Section 3(p) chunk missing traditional knowledge text: %s", c.Text)
			}
		}
	}

	if !found3p {
		t.Errorf("expected Section 3(p) chunk to be isolated cleanly, chunks were: %+v", chunks)
	}
}

func TestEnricherSection3pMetadata(t *testing.T) {
	doc := &ParsedDocument{
		Title:          "The Patents Act, 1970",
		OfficialCode:   "Act No. 39 of 1970",
		JurisdictionID: "IN",
		CategoryID:     "PATENT",
	}

	chunker := NewChunker()
	enricher := NewEnricher()

	rawChunks := []RawChunk{
		{
			Index:            1,
			SectionReference: "Section 3(p)",
			PageNumber:       1,
			Text:             "An invention which in effect is traditional knowledge.",
		},
	}

	modelChunks := chunker.ConvertToModelChunks("test-doc-1", rawChunks)
	enricher.EnrichChunk(doc, &modelChunks[0])

	metaStr := string(modelChunks[0].Metadata)
	if !contains(metaStr, "traditional knowledge") && !contains(metaStr, "Traditional Knowledge Exclusion") {
		t.Errorf("expected Section 3(p) statutory bar metadata tag, got: %s", metaStr)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
