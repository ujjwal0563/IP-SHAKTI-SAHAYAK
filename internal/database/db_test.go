package database

import (
	"context"
	"strings"
	"testing"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

func TestFormatVectorString(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected string
	}{
		{
			name:     "empty slice",
			input:    []float32{},
			expected: "[]",
		},
		{
			name:     "single element",
			input:    []float32{0.5},
			expected: "[0.5]",
		},
		{
			name:     "multiple elements",
			input:    []float32{0.01, 0.02, 0.03},
			expected: "[0.01,0.02,0.03]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FormatVectorString(tt.input)
			if actual != tt.expected {
				t.Errorf("FormatVectorString() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestDefaultSchemaSQL(t *testing.T) {
	// Verify that the fallback schema contains essential tables and extensions
	requiredKeywords := []string{
		"CREATE EXTENSION IF NOT EXISTS \"vector\"",
		"CREATE TABLE IF NOT EXISTS jurisdictions",
		"CREATE TABLE IF NOT EXISTS ip_categories",
		"CREATE TABLE IF NOT EXISTS sources",
		"CREATE TABLE IF NOT EXISTS documents",
		"CREATE TABLE IF NOT EXISTS document_chunks",
		"embedding vector(1536)",
		"idx_chunks_embedding_hnsw",
		"dmra_conditions",
		"abs_exemptions",
	}

	for _, kw := range requiredKeywords {
		if !strings.Contains(DefaultSchemaSQL, kw) {
			t.Errorf("DefaultSchemaSQL missing required keyword: %q", kw)
		}
	}
}

func TestNilDBHandling(t *testing.T) {
	var db *DB
	ctx := context.Background()

	if err := db.Ping(ctx); err == nil {
		t.Error("expected error for Ping on nil DB, got nil")
	}

	if _, err := db.GetJurisdictions(ctx); err == nil {
		t.Error("expected error for GetJurisdictions on nil DB, got nil")
	}

	if _, err := db.GetIPCategories(ctx); err == nil {
		t.Error("expected error for GetIPCategories on nil DB, got nil")
	}

	if _, err := db.GetSources(ctx); err == nil {
		t.Error("expected error for GetSources on nil DB, got nil")
	}

	if _, err := db.GetDocuments(ctx, "", ""); err == nil {
		t.Error("expected error for GetDocuments on nil DB, got nil")
	}

	if _, err := db.GetDmraConditions(ctx); err == nil {
		t.Error("expected error for GetDmraConditions on nil DB, got nil")
	}

	if _, err := db.GetAbsExemptions(ctx); err == nil {
		t.Error("expected error for GetAbsExemptions on nil DB, got nil")
	}

	if err := db.Migrate(ctx); err == nil {
		t.Error("expected error for Migrate on nil DB, got nil")
	}
}

func TestModelStructures(t *testing.T) {
	j := models.Jurisdiction{
		ID:          "IN",
		Name:        "National (India)",
		Description: "Domestic Indian Law",
	}
	if j.ID != "IN" {
		t.Errorf("expected jurisdiction ID IN, got %s", j.ID)
	}

	c := models.IPCategory{
		ID:          "PATENT",
		Name:        "Patents & Formulations",
		Description: "Section 3 exclusions",
	}
	if c.ID != "PATENT" {
		t.Errorf("expected category ID PATENT, got %s", c.ID)
	}

	chunk := models.DocumentChunk{
		ChunkIndex:       1,
		SectionReference: "Section 3(p)",
		Embedding:        []float32{0.01, 0.02},
	}
	if chunk.ChunkIndex != 1 || chunk.SectionReference != "Section 3(p)" {
		t.Errorf("unexpected document chunk fields: %+v", chunk)
	}
}
