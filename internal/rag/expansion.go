package rag

import (
	"context"
	"strings"
)

// ExpandedQuery contains the original query, decomposed sub-queries, and statutory expansion terms
type ExpandedQuery struct {
	OriginalQuery     string   `json:"original_query"`
	NormalizedQuery   string   `json:"normalized_query"`
	SubQueries        []string `json:"sub_queries"`
	ExpansionKeywords []string `json:"expansion_keywords"`
	IdentifiedDomains []string `json:"identified_domains"`
	DetectedForms     []string `json:"detected_forms"`
}

// QueryExpander analyzes queries, identifies multi-part regulatory questions, and enriches search terms
type QueryExpander struct {
	glossary *BilingualGlossary
}

// NewQueryExpander creates a query expander with statutory glossary backing
func NewQueryExpander(glossary *BilingualGlossary) *QueryExpander {
	if glossary == nil {
		glossary = NewBilingualGlossary()
	}
	return &QueryExpander{
		glossary: glossary,
	}
}

// Expand analyzes, decomposes, and enriches the user query for hybrid retrieval
func (qe *QueryExpander) Expand(ctx context.Context, query, language string) ExpandedQuery {
	cleanQuery := strings.TrimSpace(query)
	lower := strings.ToLower(cleanQuery)

	res := ExpandedQuery{
		OriginalQuery: cleanQuery,
		SubQueries:    []string{cleanQuery},
	}

	// 1. Identify concepts and Ayurvedic dosage forms via Glossary
	concepts := qe.glossary.FindConcepts(cleanQuery)
	forms := qe.glossary.FindFormulations(cleanQuery)

	domainMap := make(map[string]bool)
	var keywords []string
	var detectedForms []string

	for _, c := range concepts {
		domainMap[c.Category] = true
		keywords = append(keywords, c.Statute, c.EnglishTerm)
	}

	for _, f := range forms {
		detectedForms = append(detectedForms, f.FormName)
		keywords = append(keywords, f.StatutoryCategory, f.ApplicableAct)
		// Ayurvedic formulations naturally trigger Patent Section 3(p) and Regulatory considerations
		domainMap["PATENT"] = true
		domainMap["REGULATORY"] = true
	}

	// 2. Identify cross-domain keywords
	if strings.Contains(lower, "cure") || strings.Contains(lower, "diabetes") ||
		strings.Contains(lower, "cancer") || strings.Contains(lower, "obesity") ||
		strings.Contains(lower, "advertis") || strings.Contains(lower, "इलाज") ||
		strings.Contains(lower, "विज्ञापन") {
		domainMap["DMRA"] = true
		keywords = append(keywords, "DMRA 1954 Section 3", "prohibited objectionable advertisements")
	}

	if strings.Contains(lower, "herb") || strings.Contains(lower, "plant") ||
		strings.Contains(lower, "forest") || strings.Contains(lower, "cultivat") ||
		strings.Contains(lower, "nba") || strings.Contains(lower, "sbb") ||
		strings.Contains(lower, "जड़ी-बूटी") || strings.Contains(lower, "बायो") {
		domainMap["ABS"] = true
		keywords = append(keywords, "Biological Diversity Act 2023 Section 6 NBA Form III", "Section 7 SBB intimation")
	}

	if strings.Contains(lower, "export") || strings.Contains(lower, "foreign") ||
		strings.Contains(lower, "us fda") || strings.Contains(lower, "america") ||
		strings.Contains(lower, "विदेश") || strings.Contains(lower, "निर्यात") {
		domainMap["INTERNATIONAL"] = true
		keywords = append(keywords, "US FDA Botanical Guidance", "WIPO GRATK Treaty", "Nagoya Protocol")
	}

	// 3. Cross-Lingual Hindi to English Legal Query Mapping
	// If query is in Hindi (contains Devanagari characters), inject English statutory query
	if isHindiText(cleanQuery) {
		hindiExpansion := qe.buildCrossLingualEnglishQuery(cleanQuery, concepts, forms)
		if hindiExpansion != "" {
			res.SubQueries = append(res.SubQueries, hindiExpansion)
			keywords = append(keywords, "Cross-lingual translated legal search vector")
		}
	}

	// 4. Multi-Intent Query Decomposition
	decomposed := qe.decomposeMultiIntent(cleanQuery, domainMap, forms)
	if len(decomposed) > 1 {
		res.SubQueries = append(res.SubQueries, decomposed...)
	}

	// Remove duplicate subqueries
	res.SubQueries = deduplicateStrings(res.SubQueries)
	res.ExpansionKeywords = deduplicateStrings(keywords)
	res.DetectedForms = detectedForms

	for d := range domainMap {
		res.IdentifiedDomains = append(res.IdentifiedDomains, d)
	}

	return res
}

func (qe *QueryExpander) decomposeMultiIntent(query string, domains map[string]bool, forms []AyurvedicFormulationCategory) []string {
	var subQueries []string
	lower := strings.ToLower(query)

	// Case 1: Composite Patent + Advertising/Disease Cure Query
	if (domains["PATENT"] || len(forms) > 0) && (domains["DMRA"] || strings.Contains(lower, "cure") || strings.Contains(lower, "diabetes")) {
		subQueries = append(subQueries, "Patents Act 1970 Section 3(p) Traditional Knowledge exclusion and Section 3(e) synergy")
		subQueries = append(subQueries, "Drugs and Magic Remedies Act 1954 Section 3 prohibition on scheduled disease cure claims")
	}

	// Case 2: Composite Patent + Biodiversity/ABS Query
	if domains["PATENT"] && (domains["ABS"] || strings.Contains(lower, "biological resource") || strings.Contains(lower, "herb")) {
		subQueries = append(subQueries, "Biological Diversity Act 2023 Section 6 mandatory NBA Form III approval before patent grant")
	}

	// Case 3: Composite Domestic Regulatory + Export Query
	if (domains["REGULATORY"] || len(forms) > 0) && (domains["INTERNATIONAL"] || strings.Contains(lower, "export") || strings.Contains(lower, "us")) {
		subQueries = append(subQueries, "AYUSH licensing Rule 158B compared with US FDA Botanical Drug Guidance")
	}

	return subQueries
}

func (qe *QueryExpander) buildCrossLingualEnglishQuery(hindiQuery string, concepts []LegalTermConcept, forms []AyurvedicFormulationCategory) string {
	var parts []string

	lower := strings.ToLower(hindiQuery)
	if strings.Contains(lower, "पेटेंट") || strings.Contains(lower, "एकाधिकार") {
		parts = append(parts, "Patentability Patents Act 1970 Section 3(p) traditional knowledge")
	}
	if strings.Contains(lower, "चूर्ण") || strings.Contains(lower, "क्वाथ") || strings.Contains(lower, "काढ़ा") || strings.Contains(lower, "भस्म") {
		parts = append(parts, "Classical Ayurvedic formulation First Schedule Rule 158B non-patentable")
	}
	if strings.Contains(lower, "जड़ी") || strings.Contains(lower, "बूटी") || strings.Contains(lower, "पादप") || strings.Contains(lower, "बायो") {
		parts = append(parts, "Biological Diversity Act 2023 Section 6 NBA Form III access benefit sharing")
	}
	if strings.Contains(lower, "इलाज") || strings.Contains(lower, "रोग") || strings.Contains(lower, "दवा") {
		parts = append(parts, "Drugs and Magic Remedies Act 1954 DMRA prohibited disease claims")
	}

	if len(parts) == 0 {
		return "Ayurveda IP and regulatory guidance statutory provisions"
	}

	return strings.Join(parts, " ")
}

func isHindiText(text string) bool {
	// Unicode range for Devanagari script: 0900–097F
	for _, r := range text {
		if r >= 0x0900 && r <= 0x097F {
			return true
		}
	}
	return false
}

func deduplicateStrings(items []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, item := range items {
		clean := strings.TrimSpace(item)
		if clean != "" && !seen[clean] {
			seen[clean] = true
			out = append(out, clean)
		}
	}
	return out
}
