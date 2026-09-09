package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ParsedDocument contains the extracted text and detected header metadata
type ParsedDocument struct {
	Title          string            `json:"title"`
	OfficialCode   string            `json:"official_code"`
	Authority      string            `json:"authority"`
	JurisdictionID string            `json:"jurisdiction_id"`
	CategoryID     string            `json:"category_id"`
	RawContent     string            `json:"raw_content"`
	Metadata       map[string]string `json:"metadata"`
}

// Parser handles loading and cleaning documents from various formats (.md, .txt, .json, .pdf)
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile inspects file extension and extracts content into a ParsedDocument
func (p *Parser) ParseFile(filePath string) (*ParsedDocument, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return p.parseJSON(data, filePath)
	case ".pdf":
		return p.parsePDF(data, filePath)
	case ".md", ".txt":
		return p.parseTextOrMarkdown(string(data), filePath)
	default:
		// Default to treating as plain text
		return p.parseTextOrMarkdown(string(data), filePath)
	}
}

func (p *Parser) parseTextOrMarkdown(content, filePath string) (*ParsedDocument, error) {
	doc := &ParsedDocument{
		Title:          strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)),
		JurisdictionID: "IN",
		CategoryID:     detectCategoryFromPath(filePath),
		Metadata:       make(map[string]string),
	}

	// Extract YAML-like or Markdown header metadata if present
	// e.g. **Official Code:** Act No. 39 of 1970
	lines := strings.Split(content, "\n")
	var bodyLines []string

	titleRegex := regexp.MustCompile(`(?i)^#\s+(.+)$`)
	metaRegex := regexp.MustCompile(`(?i)^\*\*(.+?):\*\*\s*(.+)$`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if m := titleRegex.FindStringSubmatch(trimmed); len(m) > 1 && doc.Title == strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)) {
			doc.Title = strings.TrimSpace(m[1])
			continue
		}
		if m := metaRegex.FindStringSubmatch(trimmed); len(m) > 2 {
			key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(m[1]), " ", "_"))
			val := strings.TrimSpace(m[2])
			switch key {
			case "official_code":
				doc.OfficialCode = val
			case "authority":
				doc.Authority = val
			case "jurisdiction":
				doc.JurisdictionID = strings.ToUpper(val)
			case "category":
				doc.CategoryID = strings.ToUpper(val)
			default:
				doc.Metadata[key] = val
			}
			continue
		}
		bodyLines = append(bodyLines, line)
	}

	doc.RawContent = p.cleanText(strings.Join(bodyLines, "\n"))
	return doc, nil
}

func (p *Parser) parseJSON(data []byte, filePath string) (*ParsedDocument, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid json in %s: %w", filePath, err)
	}

	doc := &ParsedDocument{
		Title:          getString(raw, "title", strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))),
		OfficialCode:   getString(raw, "official_code", ""),
		Authority:      getString(raw, "authority", ""),
		JurisdictionID: getString(raw, "jurisdiction_id", "IN"),
		CategoryID:     getString(raw, "category_id", detectCategoryFromPath(filePath)),
		RawContent:     p.cleanText(getString(raw, "content", getString(raw, "text", ""))),
		Metadata:       make(map[string]string),
	}

	return doc, nil
}

// parsePDF extracts text from standard digital PDF files using ledongthuc/pdf
func (p *Parser) parsePDF(data []byte, filePath string) (*ParsedDocument, error) {
	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("parsing pdf %s: %w", filePath, err)
	}

	b, err := reader.GetPlainText()
	if err != nil {
		return nil, fmt.Errorf("extracting text from pdf %s: %w", filePath, err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(b)
	if err != nil {
		return nil, fmt.Errorf("reading extracted text from pdf %s: %w", filePath, err)
	}

	textContent := buf.String()

	return p.parseTextOrMarkdown(textContent, filePath)
}

func (p *Parser) cleanText(text string) string {
	// Remove carriage returns
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Collapse repeated empty lines (> 2) into double newline
	reEmptyLines := regexp.MustCompile(`\n{3,}`)
	text = reEmptyLines.ReplaceAllString(text, "\n\n")

	// Strip common PDF page footer noise: "Page X of Y" or "--- Page X ---"
	rePageFooters := regexp.MustCompile(`(?i)(page\s+\d+\s+of\s+\d+|---+\s*page\s+\d+\s*---+|\bpg\.\s*\d+\b)`)
	text = rePageFooters.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

func filterPrintableText(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		if (b >= 32 && b <= 126) || b == '\n' || b == '\t' {
			sb.WriteByte(b)
		}
	}
	return sb.String()
}

func detectCategoryFromPath(path string) string {
	lower := strings.ToLower(path)
	switch {
	case strings.Contains(lower, "patent"):
		return "PATENT"
	case strings.Contains(lower, "trademark") || strings.Contains(lower, "trade_mark"):
		return "TRADEMARK"
	case strings.Contains(lower, "biological") || strings.Contains(lower, "biodiversity") || strings.Contains(lower, "abs"):
		return "ABS"
	case strings.Contains(lower, "geographical") || strings.Contains(lower, "gi"):
		return "GI"
	case strings.Contains(lower, "design"):
		return "DESIGN"
	case strings.Contains(lower, "copyright"):
		return "COPYRIGHT"
	case strings.Contains(lower, "ayush") || strings.Contains(lower, "regulatory") || strings.Contains(lower, "dmra"):
		return "REGULATORY"
	default:
		return "PATENT"
	}
}

func getString(m map[string]any, key, fallback string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok && s != "" {
			return s
		}
	}
	return fallback
}
