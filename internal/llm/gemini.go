package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GeminiClient interacts with Google Gemini API via official REST v1beta
type GeminiClient struct {
	apiKey     string
	model      string
	embedModel string
	httpClient *http.Client
}

// NewGeminiClient creates a client configured with Google AI Studio key
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey:     apiKey,
		model:      "gemini-1.5-flash",
		embedModel: "text-embedding-004",
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

// Generate calls Gemini generateContent
func (g *GeminiClient) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

	reqBody := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]string{
				{"text": systemPrompt},
			},
		},
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": userPrompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature": 0.2, // Low temperature for high factual accuracy
			"maxOutputTokens": 2048,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini api error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", fmt.Errorf("decoding gemini response: %w", err)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("empty response from gemini")
}

// GenerateStream calls streamGenerateContent?alt=sse and emits words to channel
func (g *GeminiClient) GenerateStream(ctx context.Context, systemPrompt, userPrompt string, tokenChan chan<- string) error {
	defer close(tokenChan)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:streamGenerateContent?alt=sse&key=%s", g.model, g.apiKey)

	reqBody := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]string{
				{"text": systemPrompt},
			},
		},
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": userPrompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature": 0.2, // Low temperature for high factual accuracy
			"maxOutputTokens": 2048,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gemini streaming api error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gemini streaming api returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// We read the streaming response. The format is a JSON array of events, or SSE lines starting with 'data: '
	// Wait, Gemini streamGenerateContent?alt=sse returns Server-Sent Events.
	// Each event starts with "data: " followed by JSON, or sometimes just raw JSON depending on the API.
	// Let's implement a simple scanner.
	
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading gemini stream: %w", err)
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("data: ")) {
			line = line[6:]
		}

		// Skip comments or irrelevant lines
		if bytes.HasPrefix(line, []byte(":")) {
			continue
		}

		var chunk struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}

		if err := json.Unmarshal(line, &chunk); err == nil {
			if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case tokenChan <- chunk.Candidates[0].Content.Parts[0].Text:
					// sent chunk successfully
				}
			}
		}
	}

	return nil
}

// Embed generates embeddings using Gemini text-embedding-004
func (g *GeminiClient) Embed(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent?key=%s", g.embedModel, g.apiKey)

	reqBody := map[string]any{
		"model": fmt.Sprintf("models/%s", g.embedModel),
		"content": map[string]any{
			"parts": []map[string]string{
				{"text": text},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini embedding error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	// Pad or resize to 1536 if needed to match PostgreSQL column
	vals := res.Embedding.Values
	if len(vals) < 1536 {
		padded := make([]float32, 1536)
		copy(padded, vals)
		return padded, nil
	} else if len(vals) > 1536 {
		return vals[:1536], nil
	}

	return vals, nil
}

// EmbedBatch embeds multiple texts
func (g *GeminiClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var batch [][]float32
	for _, t := range texts {
		v, err := g.Embed(ctx, t)
		if err != nil {
			return nil, err
		}
		batch = append(batch, v)
	}
	return batch, nil
}
