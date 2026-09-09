package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/config"
)

// Client abstracts LLM generation and embedding capabilities
type Client interface {
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	GenerateStream(ctx context.Context, systemPrompt, userPrompt string, tokenChan chan<- string) error
	Embed(ctx context.Context, text string) ([]float32, error)
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

// FailoverClient wraps multiple clients and cascades on failure
type FailoverClient struct {
	clients []Client
}

func (f *FailoverClient) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	var lastErr error
	for _, c := range f.clients {
		res, err := c.Generate(ctx, systemPrompt, userPrompt)
		if err == nil {
			return res, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func (f *FailoverClient) GenerateStream(ctx context.Context, systemPrompt, userPrompt string, tokenChan chan<- string) error {
	var lastErr error
	for i, c := range f.clients {
		// We need a local channel for each attempt to avoid closing the main channel prematurely
		// Actually, streaming failover is tricky. Let's do a simple generate check first or just try the stream.
		// If the stream errs out early, we try the next.
		
		// To properly handle tokenChan closure, we shouldn't pass the original tokenChan directly if we want to retry.
		// However, for simplicity, if the request fails immediately, GenerateStream will return an error before sending anything.
		
		// If it's the last client, pass the tokenChan directly.
		if i == len(f.clients)-1 {
			return c.GenerateStream(ctx, systemPrompt, userPrompt, tokenChan)
		}
		
		// For non-last clients, use a dummy channel to capture immediate failures.
		// But GenerateStream already delegates to the implementation which closes the channel.
		// We'll have to use a wrapper channel.
		ch := make(chan string)
		
		go func(c Client) {
			_ = c.GenerateStream(ctx, systemPrompt, userPrompt, ch)
		}(c)
		
		// Read first token to see if it succeeds
		token, ok := <-ch
		if !ok {
			// Failed or empty, try next
			lastErr = fmt.Errorf("stream failed on client %d", i)
			continue
		}
		
		// Succeeded, forward the first token and pump the rest
		go func() {
			defer close(tokenChan)
			tokenChan <- token
			for t := range ch {
				tokenChan <- t
			}
		}()
		return nil
	}
	close(tokenChan)
	return lastErr
}

func (f *FailoverClient) Embed(ctx context.Context, text string) ([]float32, error) {
	var lastErr error
	for _, c := range f.clients {
		res, err := c.Embed(ctx, text)
		if err == nil {
			return res, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (f *FailoverClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var lastErr error
	for _, c := range f.clients {
		res, err := c.EmbedBatch(ctx, texts)
		if err == nil {
			return res, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

// NewClient initializes the configured LLM provider with failover to Mock
func NewClient(cfg *config.Config) Client {
	var primary Client
	
	provider := strings.ToLower(cfg.LLMProvider)
	switch provider {
	case "gemini":
		if cfg.LLMAPIKey != "" {
			primary = NewGeminiClient(cfg.LLMAPIKey)
		}
	case "openai":
		if cfg.LLMAPIKey != "" {
			primary = NewOpenAIClient(cfg.LLMAPIKey)
		}
	}
	
	mock := NewMockClient()
	
	if primary != nil {
		return &FailoverClient{clients: []Client{primary, mock}}
	}

	return mock
}
