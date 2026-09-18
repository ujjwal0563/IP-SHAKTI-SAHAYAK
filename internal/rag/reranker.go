package rag

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// Reranker defines the contract for secondary precision re-ranking of retrieved candidates
type Reranker interface {
	Rerank(ctx context.Context, query string, candidates []models.SearchResult, topK int) ([]models.SearchResult, error)
}

// StatutoryCrossEncoder scores and re-ranks legal chunks by evaluating exact statutory section matches,
// keyword density, authority tiers, and domain-specific intent.
type StatutoryCrossEncoder struct {
	sectionBonus   float64
	authorityBonus float64
}

// NewStatutoryCrossEncoder creates a new statutory cross-encoder re-ranker
func NewStatutoryCrossEncoder() *StatutoryCrossEncoder {
	return &StatutoryCrossEncoder{
		sectionBonus:   0.35,
		authorityBonus: 0.15,
	}
}

// Rerank re-scores candidates and returns the top-K highest-scoring chunks
func (ce *StatutoryCrossEncoder) Rerank(ctx context.Context, query string, candidates []models.SearchResult, topK int) ([]models.SearchResult, error) {
	if len(candidates) == 0 {
		return candidates, nil
	}

	if topK <= 0 || topK > len(candidates) {
		topK = len(candidates)
	}

	queryLower := strings.ToLower(query)
	queryTokens := tokenize(queryLower)

	type scoredItem struct {
		result models.SearchResult
		score  float64
	}

	var scored []scoredItem

	for _, cand := range candidates {
		baseScore := cand.Score
		if baseScore <= 0 {
			baseScore = 0.5
		}

		// 1. Cross-matching score: Token overlap between query and chunk excerpt
		candText := strings.ToLower(cand.Excerpt + " " + cand.DocumentTitle + " " + cand.SectionReference)
		candTokens := tokenize(candText)

		overlapScore := computeTokenOverlap(queryTokens, candTokens)

		// 2. Exact Statutory Section Reference Boost
		sectionBoost := 0.0
		secRefLower := strings.ToLower(cand.SectionReference)
		if secRefLower != "" {
			// Check if specific section is mentioned in query
			if strings.Contains(queryLower, secRefLower) {
				sectionBoost += ce.sectionBonus
			}
			// Special domain matches:
			if (strings.Contains(queryLower, "traditional knowledge") || strings.Contains(queryLower, "churna") || strings.Contains(queryLower, "classical")) &&
				strings.Contains(secRefLower, "3(p)") {
				sectionBoost += ce.sectionBonus * 1.2
			}
			if (strings.Contains(queryLower, "admixture") || strings.Contains(queryLower, "combination") || strings.Contains(queryLower, "synerg")) &&
				strings.Contains(secRefLower, "3(e)") {
				sectionBoost += ce.sectionBonus * 1.2
			}
			if (strings.Contains(queryLower, "nba") || strings.Contains(queryLower, "patent") && strings.Contains(queryLower, "bio")) &&
				strings.Contains(secRefLower, "section 6") {
				sectionBoost += ce.sectionBonus
			}
			if strings.Contains(queryLower, "cure") && strings.Contains(secRefLower, "section 3") && strings.Contains(strings.ToLower(cand.DocumentTitle), "magic remedies") {
				sectionBoost += ce.sectionBonus * 1.3
			}
		}

		// 3. Authority Tier Boost (Primary statutes get preference over general references)
		authBoost := 0.0
		authLower := strings.ToLower(cand.AuthorityName)
		if strings.Contains(authLower, "patent office") || strings.Contains(authLower, "biodiversity") || strings.Contains(authLower, "ayush") {
			authBoost = ce.authorityBonus
		}

		// Final composite score (Weighted sum bounded [0.0, 1.0])
		finalScore := (0.45 * baseScore) + (0.35 * overlapScore) + sectionBoost + authBoost
		if finalScore > 1.0 {
			finalScore = 1.0
		}

		updatedResult := cand
		updatedResult.Score = math.Round(finalScore*1000) / 1000

		scored = append(scored, scoredItem{
			result: updatedResult,
			score:  finalScore,
		})
	}

	// Sort descending by score
	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Select topK
	results := make([]models.SearchResult, 0, topK)
	for i := 0; i < topK && i < len(scored); i++ {
		results = append(results, scored[i].result)
	}

	return results, nil
}

func tokenize(s string) map[string]bool {
	tokens := make(map[string]bool)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '(' && r != ')'
	})

	stopWords := map[string]bool{
		"the": true, "and": true, "of": true, "in": true, "to": true, "a": true, "is": true, "for": true,
		"or": true, "an": true, "as": true, "at": true, "by": true, "with": true, "on": true, "that": true,
		"this": true, "it": true, "from": true, "can": true, "i": true, "my": true,
	}

	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) > 1 && !stopWords[w] {
			tokens[w] = true
		}
	}
	return tokens
}

func computeTokenOverlap(queryTokens, candTokens map[string]bool) float64 {
	if len(queryTokens) == 0 || len(candTokens) == 0 {
		return 0.0
	}

	matches := 0
	for qt := range queryTokens {
		if candTokens[qt] {
			matches++
		}
	}

	return float64(matches) / float64(len(queryTokens))
}
