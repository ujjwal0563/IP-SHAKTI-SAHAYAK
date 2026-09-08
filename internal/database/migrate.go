package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// DefaultSchemaSQL contains the fallback schema if file is not on disk
const DefaultSchemaSQL = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

CREATE TABLE IF NOT EXISTS jurisdictions (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ip_categories (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL UNIQUE,
    short_code VARCHAR(50) NOT NULL UNIQUE,
    website_url VARCHAR(500),
    authority_tier INT NOT NULL DEFAULT 1 CHECK (authority_tier BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jurisdiction_id VARCHAR(10) NOT NULL REFERENCES jurisdictions(id) ON DELETE RESTRICT,
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE RESTRICT,
    category_id VARCHAR(20) REFERENCES ip_categories(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    official_code VARCHAR(100),
    publication_year INT,
    effective_date DATE,
    source_url VARCHAR(1000),
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_documents_lookup ON documents(jurisdiction_id, category_id, is_active);

CREATE TABLE IF NOT EXISTS document_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INT NOT NULL DEFAULT 0,
    section_reference VARCHAR(200) NOT NULL,
    page_number INT,
    chunk_text TEXT NOT NULL,
    tsv_content tsvector GENERATED ALWAYS AS (to_tsvector('english', chunk_text)) STORED,
    embedding vector(1536),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chunks_embedding_hnsw ON document_chunks USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);
CREATE INDEX IF NOT EXISTS idx_chunks_tsv ON document_chunks USING gin (tsv_content);
CREATE INDEX IF NOT EXISTS idx_chunks_document_id ON document_chunks(document_id);

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id VARCHAR(100) NOT NULL,
    jurisdiction_id VARCHAR(10) REFERENCES jurisdictions(id) DEFAULT 'IN',
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_role VARCHAR(20) NOT NULL CHECK (sender_role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);

CREATE TABLE IF NOT EXISTS citations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    chunk_id UUID NOT NULL REFERENCES document_chunks(id) ON DELETE CASCADE,
    citation_index INT NOT NULL DEFAULT 1,
    query_text TEXT,
    relevance_score FLOAT8 NOT NULL DEFAULT 0.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_citations_chunk_id ON citations(chunk_id);

CREATE TABLE IF NOT EXISTS dmra_conditions (
    id SERIAL PRIMARY KEY,
    schedule_number INT NOT NULL,
    condition_name_en VARCHAR(100) NOT NULL,
    condition_name_hi VARCHAR(100) NOT NULL,
    prohibition_category VARCHAR(100) NOT NULL,
    statutory_rule VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS abs_exemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(100) NOT NULL,
    statutory_clause VARCHAR(100) NOT NULL,
    exemption_details TEXT NOT NULL,
    applies_to_domestic_vaidyas BOOLEAN NOT NULL DEFAULT TRUE
);
`

// Migrate runs schema and seed migrations on the connected database
func (db *DB) Migrate(ctx context.Context) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection pool not initialized")
	}

	// 1. Read or use schema migration
	schemaSQL := loadMigrationSQL("000001_init.up.sql", DefaultSchemaSQL)
	log.Println("[Database Migrate] Applying migration 000001_init.up.sql...")
	if _, err := db.Pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("failed executing 000001_init.up.sql: %w", err)
	}

	// 2. Read seed migration
	seedSQL := loadMigrationSQL("000002_seed_authoritative_data.up.sql", "")
	if seedSQL != "" {
		log.Println("[Database Migrate] Applying migration 000002_seed_authoritative_data.up.sql...")
		if _, err := db.Pool.Exec(ctx, seedSQL); err != nil {
			return fmt.Errorf("failed executing 000002_seed_authoritative_data.up.sql: %w", err)
		}
	}

	// 3. Ensure Sprint 1 test data is seeded
	if err := db.SeedSprint1TestData(ctx); err != nil {
		log.Printf("[Database Warning] SeedSprint1TestData error: %v", err)
	}

	log.Println("[Database Migrate] All database migrations and authoritative seeds applied successfully!")
	return nil
}

func loadMigrationSQL(filename, fallback string) string {
	possibleDirs := []string{
		"migrations",
		"./migrations",
		"../migrations",
		"../../migrations",
		"/app/migrations",
	}

	for _, dir := range possibleDirs {
		p := filepath.Join(dir, filename)
		if data, err := os.ReadFile(p); err == nil && len(data) > 0 {
			return string(data)
		}
	}

	return fallback
}
