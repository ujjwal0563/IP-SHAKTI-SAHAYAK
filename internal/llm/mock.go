package llm

import (
	"context"
	"math"
	"strings"
	"time"
)

// MockClient provides deterministic responses and embeddings for offline testing
type MockClient struct{}

// NewMockClient creates a new mock LLM client
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Generate returns a grounded legal analysis based on query keywords
func (m *MockClient) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	lower := strings.ToLower(userPrompt)

	var answer strings.Builder
	answer.WriteString("### Statutory Analysis & Innovation Guidance\n\n")

	if strings.Contains(lower, "patent") || strings.Contains(lower, "ashwagandha") || strings.Contains(lower, "turmeric") {
		answer.WriteString("1. **Non-Patentability of Traditional Knowledge (Section 3(p))**:\n")
		answer.WriteString("Under **Section 3(p) of the Patents Act, 1970**, an invention which is traditional knowledge or an aggregation/duplication of known properties of traditionally known components is not patentable [^1]. Simply using classical Ayurvedic herbs like Ashwagandha or Curcuma longa does not qualify as an inventive step.\n\n")
		answer.WriteString("2. **Pathway to Patentability (Section 3(e) Synergistic Combinations)**:\n")
		answer.WriteString("To overcome Section 3(p) and Section 3(e), you must demonstrate synergistic efficacy—proving that the combined botanical components produce an unexpected therapeutic effect significantly greater than the sum of individual herbs [^1].\n\n")
	}

	if strings.Contains(lower, "abs") || strings.Contains(lower, "biodiversity") || strings.Contains(lower, "nba") || strings.Contains(lower, "vaidya") {
		answer.WriteString("3. **Biological Diversity Compliance & ABS Clearances**:\n")
		answer.WriteString("Under the **Biological Diversity Act, 2023 (Section 6)**, prior approval of the National Biodiversity Authority (NBA) is mandatory before applying for any IPR based on Indian biological resources [^2].\n")
		answer.WriteString("However, registered local Vaidyas, Hakims, and traditional healthcare practitioners utilizing biological resources for local livelihood are exempted under the **Section 7 Proviso** [^2].\n\n")
	}

	if strings.Contains(lower, "cure") || strings.Contains(lower, "diabetes") || strings.Contains(lower, "cancer") {
		answer.WriteString("4. **DMRA 1954 Advertising Warning**:\n")
		answer.WriteString("⚠️ **Statutory Warning:** Claiming that any Ayurvedic medicine 'cures' diabetes, cancer, or other scheduled conditions is strictly prohibited under **Section 3 of the Drugs and Magic Remedies (Objectionable Advertisements) Act, 1954**. Violations carry criminal penalties under Section 7.\n\n")
	}

	if answer.Len() <= len("### Statutory Analysis & Innovation Guidance\n\n") {
		answer.WriteString("Based on the statutory provisions in the database, Ayurvedic innovations are evaluated across patentability exclusions (Patents Act Section 3), Access & Benefit Sharing obligations (Biological Diversity Act 2023), and Drugs & Cosmetics Act licensing [^1].\n\n")
	}

	answer.WriteString("---\n")
	answer.WriteString("### Citations & Statutory References:\n")
	answer.WriteString("[^1]: The Patents Act, 1970, Section 3(p) & Section 3(e), Indian Patent Office (CGPDTM)\n")
	answer.WriteString("[^2]: The Biological Diversity Act, 2023, Section 6 & Section 7, National Biodiversity Authority\n")

	return answer.String(), nil
}

// GenerateStream simulates real-time token streaming
func (m *MockClient) GenerateStream(ctx context.Context, systemPrompt, userPrompt string, tokenChan chan<- string) error {
	defer close(tokenChan)

	fullText, err := m.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		return err
	}

	words := strings.Split(fullText, " ")
	for _, word := range words {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tokenChan <- word + " ":
			time.Sleep(15 * time.Millisecond)
		}
	}

	return nil
}

// Embed generates a deterministic 1536-dimensional vector from input text
func (m *MockClient) Embed(ctx context.Context, text string) ([]float32, error) {
	vec := make([]float32, 1536)
	if len(text) == 0 {
		return vec, nil
	}

	// Deterministic pseudo-embedding based on character frequencies
	var hash uint32 = 2166136261
	for i := 0; i < len(text); i++ {
		hash ^= uint32(text[i])
		hash *= 16777619
	}

	for i := 0; i < 1536; i++ {
		val := math.Sin(float64(hash) + float64(i)*0.1)
		vec[i] = float32(val * 0.05)
	}

	return vec, nil
}

// EmbedBatch generates embeddings for multiple chunks
func (m *MockClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	batch := make([][]float32, len(texts))
	for i, t := range texts {
		v, err := m.Embed(ctx, t)
		if err != nil {
			return nil, err
		}
		batch[i] = v
	}
	return batch, nil
}
