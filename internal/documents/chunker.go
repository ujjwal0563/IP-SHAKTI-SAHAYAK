package documents

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// Chunker splits statutory texts along legal boundaries (Sections, Rules, Schedules, Articles)
type Chunker struct {
	MaxTokens   int
	OverlapSize int
}

// NewChunker creates a statutory boundary chunker with standard 512 token sizing
func NewChunker() *Chunker {
	return &Chunker{
		MaxTokens:   512,
		OverlapSize: 64,
	}
}

// RawChunk represents an intermediate chunk before embedding
type RawChunk struct {
	Index            int
	SectionReference string
	PageNumber       int
	Text             string
}

// ChunkDocument performs statutory boundary slicing on parsed document text
func (c *Chunker) ChunkDocument(doc *ParsedDocument) []RawChunk {
	text := doc.RawContent
	if strings.TrimSpace(text) == "" {
		return nil
	}

	// Regex to identify natural statutory section boundaries:
	// - ### Section 3(p)
	// - Section 3. ...
	// - Rule 158B. ...
	// - Article 5. ...
	// - Schedule I ...
	boundaryRegex := regexp.MustCompile(`(?mi)(?:^|\n)(?:#{1,4}\s*)?(Section\s+\d+[a-zA-Z0-9\(\)]*|Rule\s+\d+[a-zA-Z0-9\(\)]*|Article\s+\d+[a-zA-Z0-9\(\)]*|Schedule\s+[IVXLCDM\d]+)[\s:.\-]*(.*)`)

	matches := boundaryRegex.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		// Fallback: If no explicit statutory headings detected, split by double newline paragraphs
		return c.splitByParagraphs(text)
	}

	var rawChunks []RawChunk
	chunkIndex := 1

	for i := 0; i < len(matches); i++ {
		startIdx := matches[i][0]
		var endIdx int
		if i+1 < len(matches) {
			endIdx = matches[i+1][0]
		} else {
			endIdx = len(text)
		}

		chunkSection := strings.TrimSpace(text[startIdx:endIdx])
		if len(chunkSection) == 0 {
			continue
		}

		// Extract section reference name e.g. "Section 3(p)"
		secRef := ""
		submatches := boundaryRegex.FindStringSubmatch(chunkSection)
		if len(submatches) > 1 {
			secRef = strings.TrimSpace(submatches[1])
			if len(submatches) > 2 && submatches[2] != "" {
				firstLineSuffix := strings.TrimSpace(submatches[2])
				if len(firstLineSuffix) < 40 {
					secRef = fmt.Sprintf("%s - %s", secRef, firstLineSuffix)
				}
			}
		}
		if secRef == "" {
			secRef = fmt.Sprintf("Provision %d", chunkIndex)
		}

		// Check if chunk is too long (> 2500 characters / ~600 words)
		if len(chunkSection) > 2500 {
			subChunks := c.splitLargeSection(chunkSection, secRef, chunkIndex)
			for _, sc := range subChunks {
				rawChunks = append(rawChunks, sc)
				chunkIndex++
			}
		} else {
			rawChunks = append(rawChunks, RawChunk{
				Index:            chunkIndex,
				SectionReference: secRef,
				PageNumber:       1,
				Text:             chunkSection,
			})
			chunkIndex++
		}
	}

	return rawChunks
}

// splitLargeSection divides large statutory commentary while maintaining section context
func (c *Chunker) splitLargeSection(text, baseSectionRef string, startIndex int) []RawChunk {
	words := strings.Fields(text)
	var chunks []RawChunk
	wordWindow := 400
	overlap := 50

	currIndex := startIndex
	subPart := 1

	for i := 0; i < len(words); i += (wordWindow - overlap) {
		end := i + wordWindow
		if end > len(words) {
			end = len(words)
		}
		partText := strings.Join(words[i:end], " ")

		chunks = append(chunks, RawChunk{
			Index:            currIndex,
			SectionReference: fmt.Sprintf("%s (Part %d)", baseSectionRef, subPart),
			PageNumber:       1,
			Text:             partText,
		})
		currIndex++
		subPart++

		if end == len(words) {
			break
		}
	}

	return chunks
}

func (c *Chunker) splitByParagraphs(text string) []RawChunk {
	paragraphs := strings.Split(text, "\n\n")
	var chunks []RawChunk
	idx := 1

	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if len(trimmed) < 40 {
			continue // skip empty or trivial lines
		}
		chunks = append(chunks, RawChunk{
			Index:            idx,
			SectionReference: fmt.Sprintf("Paragraph %d", idx),
			PageNumber:       1,
			Text:             trimmed,
		})
		idx++
	}

	return chunks
}

// ConvertToModelChunks converts RawChunks into database models.DocumentChunk
func (c *Chunker) ConvertToModelChunks(documentID string, rawChunks []RawChunk) []models.DocumentChunk {
	var list []models.DocumentChunk
	for _, rc := range rawChunks {
		list = append(list, models.DocumentChunk{
			DocumentID:       documentID,
			ChunkIndex:       rc.Index,
			SectionReference: rc.SectionReference,
			PageNumber:       rc.PageNumber,
			ChunkText:        rc.Text,
		})
	}
	return list
}
