package citation

import (
	"regexp"
	"strconv"
	"time"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// Verifier inspects LLM answers for footnote citations and anchors them to verified database chunks
type Verifier struct{}

// NewVerifier creates a citation verifier
func NewVerifier() *Verifier {
	return &Verifier{}
}

// ExtractAndVerifyCitations maps markers like [^1], [^2] in the answer to real SearchResult chunks
func (v *Verifier) ExtractAndVerifyCitations(answer string, sources []models.SearchResult) []models.Citation {
	var citations []models.Citation

	// Match footnote markers [^1], [^2], etc.
	markerRegex := regexp.MustCompile(`\[\^(\d+)\]`)
	matches := markerRegex.FindAllStringSubmatch(answer, -1)

	seen := make(map[int]bool)

	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		idx, err := strconv.Atoi(m[1])
		if err != nil || idx <= 0 || seen[idx] {
			continue
		}
		seen[idx] = true

		// Check if source index exists in retrieved results
		sourceIdx := idx - 1
		if sourceIdx < len(sources) {
			s := sources[sourceIdx]
			citations = append(citations, models.Citation{
				CitationID:    idx,
				CitationIndex: idx,
				ChunkID:       s.ChunkID,
				Statute:       s.DocumentTitle,
				Section:       s.SectionReference,
				Authority:     s.AuthorityName,
				URL:           s.SourceURL,
				Page:          s.PageNumber,
				Excerpt:       s.Excerpt,
				Relevance:     s.Score,
				CreatedAt:     time.Now().UTC(),
			})
		}
	}

	// If no footnote markers were used by the LLM but sources were retrieved, attach top 2 sources as grounded citations
	if len(citations) == 0 && len(sources) > 0 {
		for i := 0; i < len(sources) && i < 2; i++ {
			s := sources[i]
			citations = append(citations, models.Citation{
				CitationID:    i + 1,
				CitationIndex: i + 1,
				ChunkID:       s.ChunkID,
				Statute:       s.DocumentTitle,
				Section:       s.SectionReference,
				Authority:     s.AuthorityName,
				URL:           s.SourceURL,
				Page:          s.PageNumber,
				Excerpt:       s.Excerpt,
				Relevance:     s.Score,
				CreatedAt:     time.Now().UTC(),
			})
		}
	}

	return citations
}
