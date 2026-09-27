package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/database"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/llm"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// Retriever performs hybrid semantic vector, keyword search, query expansion, and cross-encoder re-ranking
type Retriever struct {
	db        *database.DB
	llmClient llm.Client
	glossary  *BilingualGlossary
	expander  *QueryExpander
	reranker  Reranker
}

// NewRetriever creates a new Retriever instance with full expansion and re-ranking capabilities
func NewRetriever(db *database.DB, llmClient llm.Client) *Retriever {
	glossary := NewBilingualGlossary()
	expander := NewQueryExpander(glossary)
	reranker := NewStatutoryCrossEncoder()

	return &Retriever{
		db:        db,
		llmClient: llmClient,
		glossary:  glossary,
		expander:  expander,
		reranker:  reranker,
	}
}

// SetReranker allows plugging in a custom re-ranker (e.g., external LLM-assisted re-ranker)
func (r *Retriever) SetReranker(reranker Reranker) {
	if reranker != nil {
		r.reranker = reranker
	}
}

// GetGlossary returns the underlying bilingual legal glossary
func (r *Retriever) GetGlossary() *BilingualGlossary {
	return r.glossary
}

// GetExpander returns the underlying query expander
func (r *Retriever) GetExpander() *QueryExpander {
	return r.expander
}

// Retrieve expands the query, performs multi-query hybrid search in PostgreSQL,
// and applies cross-encoder re-ranking to return the most authoritative chunks.
func (r *Retriever) Retrieve(ctx context.Context, query, jurisdiction, category string, limit int) ([]models.SearchResult, error) {
	if limit <= 0 {
		limit = 4
	}

	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil, fmt.Errorf("empty query provided")
	}

	// 1. Expand and decompose query
	expanded := r.expander.Expand(ctx, cleanQuery, "")

	var candidates []models.SearchResult
	seenChunks := make(map[string]bool)

	// 2. Multi-query candidate collection
	if r.db != nil && r.db.Ping(ctx) == nil {
		// Target pool size: fetch 3-5 candidates per subquery up to max 15
		subLimit := limit + 2
		for _, sq := range expanded.SubQueries {
			// Generate query embedding vector (1536-dim)
			queryVec, err := r.llmClient.Embed(ctx, sq)
			if err != nil {
				queryVec = make([]float32, 1536)
			}

			results, err := r.db.SearchHybridWithFilters(ctx, sq, queryVec, jurisdiction, category, subLimit)
			if err == nil {
				for _, res := range results {
					if !seenChunks[res.ChunkID] {
						seenChunks[res.ChunkID] = true
						candidates = append(candidates, res)
					}
				}
			}
		}
	}

	// 3. Graceful offline fallback if DB offline or zero candidates found
	if len(candidates) == 0 {
		mockResults := r.fallbackMockRetrieval(cleanQuery, category, expanded)
		for _, mr := range mockResults {
			if !seenChunks[mr.ChunkID] {
				seenChunks[mr.ChunkID] = true
				candidates = append(candidates, mr)
			}
		}
	}

	// 4. Precision Cross-Encoder Re-Ranking
	reranked, err := r.reranker.Rerank(ctx, cleanQuery, candidates, limit)
	if err != nil || len(reranked) == 0 {
		// If re-ranking encounters an error, gracefully fallback to candidates truncated to limit
		if len(candidates) > limit {
			return candidates[:limit], nil
		}
		return candidates, nil
	}

	return reranked, nil
}

func (r *Retriever) fallbackMockRetrieval(query, category string, expanded ExpandedQuery) []models.SearchResult {
	lower := strings.ToLower(query)
	for _, sq := range expanded.SubQueries {
		lower += " " + strings.ToLower(sq)
	}

	var list []models.SearchResult

	// Section 3(p) Patent Exclusion
	if strings.Contains(lower, "patent") || strings.Contains(lower, "herb") ||
		strings.Contains(lower, "ashwagandha") || strings.Contains(lower, "churna") ||
		strings.Contains(lower, "traditional knowledge") || category == "PATENT" {
		list = append(list, models.SearchResult{
			ChunkID:          "33333333-3333-3333-3333-333333333301",
			DocumentID:       "22222222-2222-2222-2222-222222222201",
			DocumentTitle:    "The Patents Act, 1970",
			AuthorityName:    "Indian Patent Office (CGPDTM)",
			CategoryID:       "PATENT",
			SectionReference: "Section 3(p)",
			PageNumber:       12,
			Excerpt:          "What are not inventions: An invention which in effect is traditional knowledge or which is an aggregation or duplication of known properties of traditionally known component or components.",
			SourceURL:        "https://ipindia.gov.in/patents-act.htm",
			Score:            0.98,
		})

		list = append(list, models.SearchResult{
			ChunkID:          "33333333-3333-3333-3333-333333333302",
			DocumentID:       "22222222-2222-2222-2222-222222222201",
			DocumentTitle:    "The Patents Act, 1970",
			AuthorityName:    "Indian Patent Office (CGPDTM)",
			CategoryID:       "PATENT",
			SectionReference: "Section 3(e)",
			PageNumber:       11,
			Excerpt:          "What are not inventions: A substance obtained by a mere admixture resulting only in the aggregation of the properties of the components thereof or a process for producing such substance.",
			SourceURL:        "https://ipindia.gov.in/patents-act.htm",
			Score:            0.94,
		})
	}

	// Biological Diversity Act 2023
	if strings.Contains(lower, "abs") || strings.Contains(lower, "biodiversity") ||
		strings.Contains(lower, "nba") || strings.Contains(lower, "biological resource") || category == "ABS" {
		list = append(list, models.SearchResult{
			ChunkID:          "33333333-3333-3333-3333-333333333303",
			DocumentID:       "22222222-2222-2222-2222-222222222202",
			DocumentTitle:    "The Biological Diversity Act, 2002 & 2023",
			AuthorityName:    "National Biodiversity Authority (NBA)",
			CategoryID:       "ABS",
			SectionReference: "Section 6",
			PageNumber:       6,
			Excerpt:          "No person shall apply for any intellectual property right for any invention based on any biological resource obtained from India without previous approval of the National Biodiversity Authority.",
			SourceURL:        "http://nbaindia.org/act",
			Score:            0.96,
		})

		list = append(list, models.SearchResult{
			ChunkID:          "33333333-3333-3333-3333-333333333304",
			DocumentID:       "22222222-2222-2222-2222-222222222202",
			DocumentTitle:    "The Biological Diversity Act, 2002 & 2023",
			AuthorityName:    "National Biodiversity Authority (NBA)",
			CategoryID:       "ABS",
			SectionReference: "Section 7 Proviso",
			PageNumber:       7,
			Excerpt:          "Provided that the provisions of this section shall not apply to the local people and communities of the area, including growers and cultivators of biological resources, and vaids and hakims, who have been practising indigenous medicine.",
			SourceURL:        "http://nbaindia.org/act",
			Score:            0.95,
		})
	}

	// Drugs and Magic Remedies Act 1954
	if strings.Contains(lower, "cure") || strings.Contains(lower, "diabetes") ||
		strings.Contains(lower, "cancer") || strings.Contains(lower, "advertis") ||
		strings.Contains(lower, "dmra") {
		list = append(list, models.SearchResult{
			ChunkID:          "33333333-3333-3333-3333-333333333305",
			DocumentID:       "22222222-2222-2222-2222-222222222204",
			DocumentTitle:    "Drugs and Magic Remedies (Objectionable Advertisements) Act, 1954",
			AuthorityName:    "Ministry of AYUSH & CDSCO",
			CategoryID:       "REGULATORY",
			SectionReference: "Section 3",
			PageNumber:       2,
			Excerpt:          "Prohibition of advertisement of certain drugs for treatment of certain diseases and disorders: No person shall take any part in the publication of any advertisement referring to any drug in terms which suggest that it cures diabetes, cancer, or any scheduled condition.",
			SourceURL:        "https://legislative.gov.in",
			Score:            0.97,
		})
	}

	return list
}
