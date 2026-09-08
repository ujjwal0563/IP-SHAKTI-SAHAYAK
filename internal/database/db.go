package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// DB encapsulates the pgxpool connection pool and query helpers
type DB struct {
	Pool *pgxpool.Pool
}

// New initializes a production-tuned pgx connection pool
func New(ctx context.Context, connString string) (*DB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify database connectivity with 3-second ping timeout
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		log.Printf("[Database Warning] Ping failed (%v). Continuing in standalone mode.", err)
		return &DB{Pool: pool}, err
	}

	log.Println("[Database] Connected successfully to PostgreSQL")
	return &DB{Pool: pool}, nil
}

// Close gracefully closes all idle connections in the pool
func (db *DB) Close() {
	if db != nil && db.Pool != nil {
		db.Pool.Close()
		log.Println("[Database] Connection pool closed")
	}
}

// Ping checks whether database is responsive
func (db *DB) Ping(ctx context.Context) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection pool not initialized")
	}
	return db.Pool.Ping(ctx)
}

// CheckPgvector checks whether pgvector extension is installed and available
func (db *DB) CheckPgvector(ctx context.Context) (bool, string, error) {
	if db == nil || db.Pool == nil {
		return false, "", fmt.Errorf("database connection pool not initialized")
	}

	var extVersion string
	err := db.Pool.QueryRow(ctx, "SELECT extversion FROM pg_extension WHERE extname = 'vector'").Scan(&extVersion)
	if err != nil {
		return false, "", err
	}
	return true, extVersion, nil
}

// FormatVectorString serializes a float32 slice into a PostgreSQL vector literal "[0.01, 0.02, ...]"
func FormatVectorString(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range vec {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%g", v))
	}
	sb.WriteByte(']')
	return sb.String()
}

// GetJurisdictions retrieves all supported jurisdictions
func (db *DB) GetJurisdictions(ctx context.Context) ([]models.Jurisdiction, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, COALESCE(description, ''), created_at 
		FROM jurisdictions 
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Jurisdiction
	for rows.Next() {
		var j models.Jurisdiction
		if err := rows.Scan(&j.ID, &j.Name, &j.Description, &j.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, j)
	}

	return list, rows.Err()
}

// GetIPCategories retrieves all registered statutory IP categories
func (db *DB) GetIPCategories(ctx context.Context) ([]models.IPCategory, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, COALESCE(description, ''), created_at 
		FROM ip_categories 
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.IPCategory
	for rows.Next() {
		var cat models.IPCategory
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, cat)
	}

	return list, rows.Err()
}

// GetSources retrieves all statutory authorities
func (db *DB) GetSources(ctx context.Context) ([]models.AuthoritySource, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, short_code, COALESCE(website_url, ''), authority_tier, created_at 
		FROM sources 
		ORDER BY authority_tier ASC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AuthoritySource
	for rows.Next() {
		var s models.AuthoritySource
		if err := rows.Scan(&s.ID, &s.Name, &s.ShortCode, &s.WebsiteURL, &s.AuthorityTier, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, rows.Err()
}

// GetDocuments retrieves documents, optionally filtered by jurisdiction or category
func (db *DB) GetDocuments(ctx context.Context, jurisdictionID, categoryID string) ([]models.Document, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	query := `
		SELECT id, jurisdiction_id, source_id, COALESCE(category_id, ''), title, 
		       COALESCE(official_code, ''), COALESCE(publication_year, 0), effective_date, 
		       COALESCE(source_url, ''), language, is_active, created_at, updated_at
		FROM documents
		WHERE ($1 = '' OR jurisdiction_id = $1)
		  AND ($2 = '' OR category_id = $2)
		ORDER BY publication_year DESC, title ASC
	`

	rows, err := db.Pool.Query(ctx, query, jurisdictionID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(
			&d.ID, &d.JurisdictionID, &d.SourceID, &d.CategoryID, &d.Title,
			&d.OfficialCode, &d.PublicationYear, &d.EffectiveDate,
			&d.SourceURL, &d.Language, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, rows.Err()
}

// GetDocumentByID retrieves a single document by UUID
func (db *DB) GetDocumentByID(ctx context.Context, id string) (*models.Document, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	var d models.Document
	query := `
		SELECT id, jurisdiction_id, source_id, COALESCE(category_id, ''), title, 
		       COALESCE(official_code, ''), COALESCE(publication_year, 0), effective_date, 
		       COALESCE(source_url, ''), language, is_active, created_at, updated_at
		FROM documents
		WHERE id = $1
	`
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.JurisdictionID, &d.SourceID, &d.CategoryID, &d.Title,
		&d.OfficialCode, &d.PublicationYear, &d.EffectiveDate,
		&d.SourceURL, &d.Language, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetDocumentChunks retrieves all atomic chunks for a given document
func (db *DB) GetDocumentChunks(ctx context.Context, documentID string) ([]models.DocumentChunk, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, document_id, chunk_index, section_reference, COALESCE(page_number, 0), chunk_text, COALESCE(metadata, '{}'::jsonb), created_at
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY chunk_index ASC
	`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DocumentChunk
	for rows.Next() {
		var c models.DocumentChunk
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.ChunkIndex, &c.SectionReference, &c.PageNumber, &c.ChunkText, &c.Metadata, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// SearchSimilarChunks queries pgvector using HNSW cosine distance operator (<=>)
func (db *DB) SearchSimilarChunks(ctx context.Context, queryVector []float32, limit int) ([]models.SearchResult, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}
	if limit <= 0 {
		limit = 5
	}

	vectorStr := FormatVectorString(queryVector)

	query := `
		SELECT 
			c.id, c.document_id, d.title, s.name, COALESCE(d.category_id, ''), 
			c.section_reference, COALESCE(c.page_number, 0), c.chunk_text, 
			COALESCE(d.source_url, ''), (1.0 - (c.embedding <=> $1::vector)) AS score,
			COALESCE(c.metadata, '{}'::jsonb)
		FROM document_chunks c
		JOIN documents d ON c.document_id = d.id
		JOIN sources s ON d.source_id = s.id
		WHERE c.embedding IS NOT NULL
		ORDER BY c.embedding <=> $1::vector ASC
		LIMIT $2
	`

	rows, err := db.Pool.Query(ctx, query, vectorStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		if err := rows.Scan(
			&r.ChunkID, &r.DocumentID, &r.DocumentTitle, &r.AuthorityName, &r.CategoryID,
			&r.SectionReference, &r.PageNumber, &r.Excerpt, &r.SourceURL, &r.Score, &r.Metadata,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

// SearchHybrid combines full-text tsvector keyword search and pgvector cosine distance
func (db *DB) SearchHybrid(ctx context.Context, keyword string, queryVector []float32, limit int) ([]models.SearchResult, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}
	if limit <= 0 {
		limit = 5
	}

	vectorStr := FormatVectorString(queryVector)

	query := `
		SELECT 
			c.id, c.document_id, d.title, s.name, COALESCE(d.category_id, ''), 
			c.section_reference, COALESCE(c.page_number, 0), c.chunk_text, 
			COALESCE(d.source_url, ''),
			(
				0.5 * (1.0 - (c.embedding <=> $1::vector)) + 
				0.5 * ts_rank_cd(c.tsv_content, plainto_tsquery('english', $2))
			) AS score,
			COALESCE(c.metadata, '{}'::jsonb)
		FROM document_chunks c
		JOIN documents d ON c.document_id = d.id
		JOIN sources s ON d.source_id = s.id
		WHERE c.embedding IS NOT NULL
		   OR c.tsv_content @@ plainto_tsquery('english', $2)
		ORDER BY score DESC
		LIMIT $3
	`

	rows, err := db.Pool.Query(ctx, query, vectorStr, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		if err := rows.Scan(
			&r.ChunkID, &r.DocumentID, &r.DocumentTitle, &r.AuthorityName, &r.CategoryID,
			&r.SectionReference, &r.PageNumber, &r.Excerpt, &r.SourceURL, &r.Score, &r.Metadata,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

// SearchHybridWithFilters combines keyword matching, vector similarity, and jurisdiction/category filters
func (db *DB) SearchHybridWithFilters(ctx context.Context, keyword string, queryVector []float32, jurisdictionID, categoryID string, limit int) ([]models.SearchResult, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}
	if limit <= 0 {
		limit = 5
	}

	vectorStr := FormatVectorString(queryVector)

	query := `
		SELECT 
			c.id, c.document_id, d.title, s.name, COALESCE(d.category_id, ''), 
			c.section_reference, COALESCE(c.page_number, 0), c.chunk_text, 
			COALESCE(d.source_url, ''),
			(
				0.5 * (1.0 - (c.embedding <=> $1::vector)) + 
				0.5 * ts_rank_cd(COALESCE(c.tsv_content, to_tsvector('english', c.chunk_text)), plainto_tsquery('english', $2))
			) AS score,
			COALESCE(c.metadata, '{}'::jsonb)
		FROM document_chunks c
		JOIN documents d ON c.document_id = d.id
		JOIN sources s ON d.source_id = s.id
		WHERE (c.embedding IS NOT NULL OR c.tsv_content @@ plainto_tsquery('english', $2))
		  AND ($3 = '' OR d.jurisdiction_id = $3)
		  AND ($4 = '' OR d.category_id = $4)
		ORDER BY score DESC
		LIMIT $5
	`

	rows, err := db.Pool.Query(ctx, query, vectorStr, keyword, jurisdictionID, categoryID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		if err := rows.Scan(
			&r.ChunkID, &r.DocumentID, &r.DocumentTitle, &r.AuthorityName, &r.CategoryID,
			&r.SectionReference, &r.PageNumber, &r.Excerpt, &r.SourceURL, &r.Score, &r.Metadata,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

// UpsertDocumentWithChunks transactionally saves a document and its atomic chunks
func (db *DB) UpsertDocumentWithChunks(ctx context.Context, doc models.Document, chunks []models.DocumentChunk) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection pool not initialized")
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Ensure jurisdiction exists
	if doc.JurisdictionID != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO jurisdictions (id, name)
			VALUES ($1, $2)
			ON CONFLICT (id) DO NOTHING;
		`, doc.JurisdictionID, doc.JurisdictionID)
		if err != nil {
			return fmt.Errorf("ensuring jurisdiction: %w", err)
		}
	}

	// 2. Ensure category exists
	if doc.CategoryID != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO ip_categories (id, name)
			VALUES ($1, $2)
			ON CONFLICT (id) DO NOTHING;
		`, doc.CategoryID, doc.CategoryID)
		if err != nil {
			return fmt.Errorf("ensuring category: %w", err)
		}
	}

	// Ensure valid UUID for document
	if !models.IsValidUUID(doc.ID) {
		doc.ID = models.DeterministicUUID("doc", doc.Title)
	}

	// 3. Upsert document
	docQuery := `
		INSERT INTO documents (
			id, jurisdiction_id, source_id, category_id, title, 
			official_code, publication_year, effective_date, source_url, language, is_active
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, $10, $11
		) ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			official_code = EXCLUDED.official_code,
			publication_year = EXCLUDED.publication_year,
			effective_date = EXCLUDED.effective_date,
			source_url = EXCLUDED.source_url,
			updated_at = NOW();
	`
	_, err = tx.Exec(ctx, docQuery,
		doc.ID, doc.JurisdictionID, doc.SourceID, doc.CategoryID, doc.Title,
		doc.OfficialCode, doc.PublicationYear, doc.EffectiveDate, doc.SourceURL, doc.Language, doc.IsActive,
	)
	if err != nil {
		return fmt.Errorf("upserting document: %w", err)
	}

	// 4. Upsert chunks
	chunkQuery := `
		INSERT INTO document_chunks (
			id, document_id, chunk_index, section_reference, page_number, 
			chunk_text, embedding, tsv_content, metadata
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7::vector, to_tsvector('english', $6), $8::jsonb
		) ON CONFLICT (id) DO UPDATE SET
			chunk_text = EXCLUDED.chunk_text,
			embedding = EXCLUDED.embedding,
			tsv_content = to_tsvector('english', EXCLUDED.chunk_text),
			metadata = EXCLUDED.metadata;
	`

	for i, c := range chunks {
		chunkID := c.ID
		if !models.IsValidUUID(chunkID) {
			chunkID = models.DeterministicUUID("chunk", fmt.Sprintf("%s:%d:%s", doc.ID, i+1, c.SectionReference))
		}

		vecStr := FormatVectorString(c.Embedding)
		metaJSON := string(c.Metadata)
		if metaJSON == "" {
			metaJSON = "{}"
		}

		_, err = tx.Exec(ctx, chunkQuery,
			chunkID, doc.ID, c.ChunkIndex, c.SectionReference, c.PageNumber,
			c.ChunkText, vecStr, metaJSON,
		)
		if err != nil {
			return fmt.Errorf("upserting chunk %s: %w", c.SectionReference, err)
		}
	}

	return tx.Commit(ctx)
}

// GetDmraConditions retrieves all scheduled diseases under DMRA 1954
func (db *DB) GetDmraConditions(ctx context.Context) ([]models.DmraCondition, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, schedule_number, condition_name_en, condition_name_hi, prohibition_category, statutory_rule 
		FROM dmra_conditions 
		ORDER BY schedule_number ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DmraCondition
	for rows.Next() {
		var c models.DmraCondition
		if err := rows.Scan(&c.ID, &c.ScheduleNumber, &c.ConditionNameEn, &c.ConditionNameHi, &c.ProhibitionCategory, &c.StatutoryRule); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, rows.Err()
}

// GetAbsExemptions retrieves all Biological Diversity Act statutory exemptions
func (db *DB) GetAbsExemptions(ctx context.Context) ([]models.AbsExemption, error) {
	if db == nil || db.Pool == nil {
		return nil, fmt.Errorf("database connection pool not initialized")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, category, statutory_clause, exemption_details, applies_to_domestic_vaidyas 
		FROM abs_exemptions 
		ORDER BY category ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AbsExemption
	for rows.Next() {
		var a models.AbsExemption
		if err := rows.Scan(&a.ID, &a.Category, &a.StatutoryClause, &a.ExemptionDetails, &a.AppliesToDomesticVaidyas); err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	return list, rows.Err()
}

// InsertCitation registers a verified statutory citation log
func (db *DB) InsertCitation(ctx context.Context, c models.Citation) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection pool not initialized")
	}

	query := `
		INSERT INTO citations (chunk_id, query_text, relevance_score)
		VALUES ($1, $2, $3)
	`
	_, err := db.Pool.Exec(ctx, query, c.ChunkID, c.QueryText, c.Relevance)
	return err
}

// SeedSprint1TestData satisfies PRD Sprint 1 DoD Task 7:
// "Insert 1 test source (Indian Patent Office) and 1 test document (Section 3 of the Patents Act, 1970) with a dummy vector."
func (db *DB) SeedSprint1TestData(ctx context.Context) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection pool not initialized")
	}

	// 1. Ensure 'IN' Jurisdiction
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO jurisdictions (id, name, description)
		VALUES ('IN', 'National (India)', 'Domestic Indian statutory regime')
		ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("seeding jurisdiction IN: %w", err)
	}

	// 2. Ensure 'PATENT' IP Category
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO ip_categories (id, name, description)
		VALUES ('PATENT', 'Patents & Formulations', 'Statutory patentability under Patents Act 1970')
		ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("seeding ip_category PATENT: %w", err)
	}

	// 3. Ensure test source (Indian Patent Office)
	sourceID := "11111111-1111-1111-1111-111111111101"
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO sources (id, name, short_code, website_url, authority_tier)
		VALUES ($1, 'Indian Patent Office (CGPDTM)', 'IPO-CGPDTM', 'https://ipindia.gov.in', 1)
		ON CONFLICT (id) DO NOTHING;
	`, sourceID)
	if err != nil {
		return fmt.Errorf("seeding test source: %w", err)
	}

	// 4. Ensure test document (Section 3 of the Patents Act, 1970)
	docID := "22222222-2222-2222-2222-222222222201"
	_, err = db.Pool.Exec(ctx, `
		INSERT INTO documents (id, jurisdiction_id, source_id, category_id, title, official_code, publication_year, effective_date, source_url, language, is_active)
		VALUES ($1, 'IN', $2, 'PATENT', 'The Patents Act, 1970', 'Act No. 39 of 1970', 1970, '2024-03-15', 'https://ipindia.gov.in/patents-act.htm', 'en', TRUE)
		ON CONFLICT (id) DO NOTHING;
	`, docID, sourceID)
	if err != nil {
		return fmt.Errorf("seeding test document: %w", err)
	}

	// 5. Ensure test document chunk with dummy 1536-dimensional vector
	chunkID := "33333333-3333-3333-3333-333333333301"
	dummyVector := make([]float32, 1536)
	for i := range dummyVector {
		dummyVector[i] = 0.025
	}
	vectorStr := FormatVectorString(dummyVector)

	chunkMeta, _ := json.Marshal(map[string]any{
		"category":       "PATENT",
		"act":            "Patents Act 1970",
		"statutory_bar":  true,
		"sprint_1_test":  true,
	})

	_, err = db.Pool.Exec(ctx, `
		INSERT INTO document_chunks (id, document_id, chunk_index, section_reference, page_number, chunk_text, embedding, metadata)
		VALUES (
			$1, $2, 1, 'Section 3(p)', 12, 
			'What are not inventions: An invention which in effect is traditional knowledge or which is an aggregation or duplication of known properties of traditionally known component or components.',
			$3::vector, $4::jsonb
		)
		ON CONFLICT (id) DO UPDATE SET
			embedding = EXCLUDED.embedding,
			metadata = EXCLUDED.metadata;
	`, chunkID, docID, vectorStr, string(chunkMeta))
	if err != nil {
		return fmt.Errorf("seeding test chunk with dummy vector: %w", err)
	}

	log.Println("[Database] Sprint 1 test data seeded successfully (1 test source, 1 test document, 1 dummy vector chunk)")
	return nil
}
