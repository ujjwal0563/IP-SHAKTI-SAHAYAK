package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/config"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/database"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/documents"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/llm"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

func main() {
	dirFlag := flag.String("dir", "data/raw", "Path to directory containing statutory legal texts (.md, .txt, .json, .pdf)")
	flag.Parse()

	log.Println("================================================================================")
	log.Println("            IP-SAKTI Sahayak: Statutory Knowledge Ingestion Engine             ")
	log.Println("================================================================================")

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Initialize DB connection
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("[Warning] Database connection not available: %v. Running in verification mode.\n", err)
	} else {
		defer db.Close()
	}

	// Initialize LLM / Embedder Client
	llmClient := llm.NewClient(cfg)
	parser := documents.NewParser()
	chunker := documents.NewChunker()
	enricher := documents.NewEnricher()

	// Find all eligible files in directory
	var files []string
	err = filepath.Walk(*dirFlag, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".md" || ext == ".txt" || ext == ".json" || ext == ".pdf" {
				// Ignore general README files
				if filepath.Base(path) != "README.md" {
					files = append(files, path)
				}
			}
		}
		return nil
	})

	if err != nil || len(files) == 0 {
		log.Printf("No legal corpus files found in directory %s\n", *dirFlag)
		return
	}

	log.Printf("Discovered %d statutory files for ingestion in '%s':\n", len(files), *dirFlag)
	for i, f := range files {
		log.Printf("  [%d] %s\n", i+1, f)
	}

	totalChunksIngested := 0
	startTime := time.Now()

	for _, file := range files {
		log.Printf("\n--- Processing: %s ---", file)

		// 1. Parse
		parsed, err := parser.ParseFile(file)
		if err != nil {
			log.Printf("❌ Failed to parse %s: %v\n", file, err)
			continue
		}

		log.Printf("  Title: %s | Jurisdiction: %s | Category: %s\n", parsed.Title, parsed.JurisdictionID, parsed.CategoryID)

		// 2. Statutory Chunking
		rawChunks := chunker.ChunkDocument(parsed)
		if len(rawChunks) == 0 {
			log.Printf("  ⚠️ No chunks created from %s (empty content)\n", file)
			continue
		}
		log.Printf("  Extracted %d statutory boundary chunks.\n", len(rawChunks))

		docID := fmt.Sprintf("doc-%d", time.Now().UnixNano())
		doc := models.Document{
			ID:             docID,
			JurisdictionID: parsed.JurisdictionID,
			SourceID:       "11111111-1111-1111-1111-111111111101", // Default IPO source
			CategoryID:     parsed.CategoryID,
			Title:          parsed.Title,
			OfficialCode:   parsed.OfficialCode,
			Language:       "en",
			IsActive:       true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}

		modelChunks := chunker.ConvertToModelChunks(docID, rawChunks)

		// 3. Enrich each chunk with legal tags
		for i := range modelChunks {
			modelChunks[i].ID = fmt.Sprintf("chunk-%s-%d", docID, i+1)
			enricher.EnrichChunk(parsed, &modelChunks[i])
		}

		// 4. Batch Embeddings
		chunkTexts := make([]string, len(modelChunks))
		for i, c := range modelChunks {
			chunkTexts[i] = c.ChunkText
		}

		log.Printf("  Generating 1536-dim vector embeddings for %d chunks...\n", len(chunkTexts))
		embeddings, err := llmClient.EmbedBatch(ctx, chunkTexts)
		if err != nil {
			log.Printf("  ⚠️ Embedding batch failed: %v. Using fallback coordinates.\n", err)
			for i := range modelChunks {
				modelChunks[i].Embedding = make([]float32, 1536)
			}
		} else {
			for i := range modelChunks {
				if i < len(embeddings) {
					modelChunks[i].Embedding = embeddings[i]
				}
			}
		}

		// 5. Upsert into PostgreSQL if online
		if db != nil && db.Pool != nil && db.Ping(ctx) == nil {
			if err := db.UpsertDocumentWithChunks(ctx, doc, modelChunks); err != nil {
				log.Printf("  ❌ Database upsert error: %v\n", err)
			} else {
				log.Printf("  ✅ Successfully committed %d chunks into PostgreSQL (pgvector HNSW + tsvector)\n", len(modelChunks))
			}
		} else {
			log.Printf("  ℹ️ Database offline; verified and structured %d chunks in memory.\n", len(modelChunks))
		}

		totalChunksIngested += len(modelChunks)
	}

	duration := time.Since(startTime)
	log.Printf("\n================================================================================")
	log.Printf("Ingestion Complete: Processed %d files, generated %d statutory chunks in %v\n", len(files), totalChunksIngested, duration)
	log.Printf("================================================================================\n")
}
