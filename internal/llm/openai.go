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

// OpenAIClient interacts with OpenAI REST API
type OpenAIClient struct {
	apiKey     string
	model      string
	embedModel string
	httpClient *http.Client
}

// NewOpenAIClient creates an OpenAI client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:     apiKey,
		model:      "gpt-4o",
		embedModel: "text-embedding-3-small",
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

// Generate sends chat completion request to OpenAI
func (o *OpenAIClient) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	reqBody := map[string]any{
		"model": o.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.2,
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
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai status %d: %s", resp.StatusCode, string(b))
	}

	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) > 0 {
		return res.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty choices from openai")
}

// GenerateStream streams tokens from OpenAI chat completions
func (o *OpenAIClient) GenerateStream(ctx context.Context, systemPrompt, userPrompt string, tokenChan chan<- string) error {
	defer close(tokenChan)

	url := "https://api.openai.com/v1/chat/completions"

	reqBody := map[string]any{
		"model": o.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.2,
		"stream":      true,
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
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("openai streaming error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("openai streaming status %d: %s", resp.StatusCode, string(b))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("data: ")) {
			line = line[6:]
		}

		if string(line) == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(line, &chunk); err == nil {
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case tokenChan <- chunk.Choices[0].Delta.Content:
				}
			}
		}
	}

	return nil
}

// Embed generates 1536-dim vector via text-embedding-3-small
func (o *OpenAIClient) Embed(ctx context.Context, text string) ([]float32, error) {
	url := "https://api.openai.com/v1/embeddings"

	reqBody := map[string]any{
		"model": o.embedModel,
		"input": text,
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
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai embedding status %d: %s", resp.StatusCode, string(b))
	}

	var res struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Data) > 0 {
		return res.Data[0].Embedding, nil
	}

	return nil, fmt.Errorf("empty embedding from openai")
}

// EmbedBatch creates embeddings for multiple inputs
func (o *OpenAIClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var batch [][]float32
	for _, t := range texts {
		v, err := o.Embed(ctx, t)
		if err != nil {
			return nil, err
		}
		batch = append(batch, v)
	}
	return batch, nil
}
