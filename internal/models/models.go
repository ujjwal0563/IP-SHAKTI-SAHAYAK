package models

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// IsValidUUID checks whether a string conforms to the 8-4-4-4-12 hex UUID format
func IsValidUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

// NewUUID generates an RFC 4122 compliant UUID v4 string
func NewUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// DeterministicUUID generates a deterministic RFC 4122 compliant UUID (v5 style) from namespace and name
func DeterministicUUID(namespace, name string) string {
	h := sha1.Sum([]byte(namespace + ":" + name))
	h[6] = (h[6] & 0x0f) | 0x50 // Version 5
	h[8] = (h[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", h[0:4], h[4:6], h[6:8], h[8:10], h[10:16])
}

// Jurisdiction represents a legal territory ("IN" for India, "INT" for International)
type Jurisdiction struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// IPCategory represents statutory IP categories (PATENT, TRADEMARK, GI, TK, ABS, REGULATORY, DESIGN)
type IPCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// AuthoritySource represents an issuing statutory authority (e.g. CGPDTM, NBA, CDSCO)
type AuthoritySource struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ShortCode     string    `json:"short_code"`
	WebsiteURL    string    `json:"website_url,omitempty"`
	AuthorityTier int       `json:"authority_tier"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

// Document represents an official legal act, rule, or gazette
type Document struct {
	ID              string     `json:"id"`
	JurisdictionID  string     `json:"jurisdiction_id"`
	SourceID        string     `json:"source_id"`
	CategoryID      string     `json:"category_id,omitempty"`
	Title           string     `json:"title"`
	OfficialCode    string     `json:"official_code,omitempty"`
	PublicationYear int        `json:"publication_year,omitempty"`
	EffectiveDate   *time.Time `json:"effective_date,omitempty"`
	SourceURL       string     `json:"source_url,omitempty"`
	Language        string     `json:"language,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// DocumentChunk represents a discrete, atomic statutory provision
type DocumentChunk struct {
	ID               string          `json:"id"`
	DocumentID       string          `json:"document_id"`
	ChunkIndex       int             `json:"chunk_index"`
	SectionReference string          `json:"section_reference"`
	PageNumber       int             `json:"page_number,omitempty"`
	ChunkText        string          `json:"chunk_text"`
	Embedding        []float32       `json:"embedding,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

// SearchResult represents an enriched document chunk returned from hybrid vector/keyword search
type SearchResult struct {
	ChunkID          string          `json:"chunk_id"`
	DocumentID       string          `json:"document_id"`
	DocumentTitle    string          `json:"document_title"`
	AuthorityName    string          `json:"authority_name"`
	CategoryID       string          `json:"category_id,omitempty"`
	SectionReference string          `json:"section_reference"`
	PageNumber       int             `json:"page_number,omitempty"`
	Excerpt          string          `json:"excerpt"`
	SourceURL        string          `json:"source_url,omitempty"`
	Score            float64         `json:"score"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
}

// Citation records a verified claim grounding
type Citation struct {
	ID            string    `json:"id,omitempty"`
	CitationID    int       `json:"citation_id,omitempty"`
	MessageID     string    `json:"message_id,omitempty"`
	ChunkID       string    `json:"chunk_id,omitempty"`
	CitationIndex int       `json:"citation_index,omitempty"`
	Statute       string    `json:"statute,omitempty"`
	Section       string    `json:"section,omitempty"`
	Authority     string    `json:"authority,omitempty"`
	URL           string    `json:"url,omitempty"`
	Page          int       `json:"page,omitempty"`
	Excerpt       string    `json:"excerpt,omitempty"`
	QueryText     string    `json:"query_text,omitempty"`
	Relevance     float64   `json:"relevance,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

// DmraCondition represents a scheduled disease prohibited under DMRA 1954
type DmraCondition struct {
	ID                  int    `json:"id"`
	ScheduleNumber      int    `json:"schedule_number"`
	ConditionNameEn     string `json:"condition_name_en"`
	ConditionNameHi     string `json:"condition_name_hi"`
	ProhibitionCategory string `json:"prohibition_category"`
	StatutoryRule       string `json:"statutory_rule"`
}

// AbsExemption represents an exemption category under the Biological Diversity Act
type AbsExemption struct {
	ID                       string `json:"id"`
	Category                 string `json:"category"`
	StatutoryClause          string `json:"statutory_clause"`
	ExemptionDetails         string `json:"exemption_details"`
	AppliesToDomesticVaidyas bool   `json:"applies_to_domestic_vaidyas"`
}

// HealthResponse represents system liveness and database pool status
type HealthResponse struct {
	Status            string    `json:"status"`
	Database          string    `json:"database"`
	Pgvector          string    `json:"pgvector"`
	Timestamp         time.Time `json:"timestamp"`
	Version           string    `json:"version"`
	ZeroDataRetention bool      `json:"zero_data_retention"`
}

// APIVersionResponse represents API version info
type APIVersionResponse struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// User represents a system user
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Conversation represents a chat session
type Conversation struct {
	ID             string    `json:"id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id,omitempty"`
	JurisdictionID string    `json:"jurisdiction_id,omitempty"`
	Language       string    `json:"language,omitempty"`
	Title          string    `json:"title,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// Message represents a single chat message
type Message struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	SenderRole     string     `json:"sender_role,omitempty"`
	Role           string     `json:"role,omitempty"` // "user" or "assistant"
	Content        string     `json:"content"`
	Citations      []Citation `json:"citations,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ChatRequest represents the incoming user query payload
type ChatRequest struct {
	SessionID    string `json:"session_id,omitempty"`
	Query        string `json:"query"`
	Jurisdiction string `json:"jurisdiction,omitempty"` // "IN" (default) or "INT"
	Category     string `json:"category,omitempty"`     // "PATENT", "ABS", "TRADEMARK", etc.
	Language     string `json:"language,omitempty"`     // "en" (default) or "hi"
}

// ChatResponse represents the generated grounded answer with citations and guardrail warnings
type ChatResponse struct {
	SessionID string         `json:"session_id"`
	Answer    string         `json:"answer"`
	Citations []Citation     `json:"citations"`
	Warnings  []string       `json:"warnings,omitempty"`
	Sources   []SearchResult `json:"sources,omitempty"`
}

// AssessmentInput represents input to the 5-step guided IP & Regulatory Wizard
type AssessmentInput struct {
	InnovationType       string `json:"innovation_type"`       // "classical_ayurveda", "novel_extract", "synergistic_combo", "device"
	UsesIndianBioResource bool  `json:"uses_indian_bio_resource"`
	TargetMarket         string `json:"target_market"`         // "domestic_only", "export_us", "export_eu", "global"
	IntendsToPatent      bool   `json:"intends_to_patent"`
	HasDiseaseClaims     bool   `json:"has_disease_claims"`
	ClaimedDiseases      []string `json:"claimed_diseases,omitempty"`
}

// AssessmentStep summarizes compliance requirements for a specific legal domain
type AssessmentStep struct {
	Domain       string   `json:"domain"`       // "PATENT", "ABS_BIODIVERSITY", "AYUSH_REGULATORY", "DMRA"
	Status       string   `json:"status"`       // "ALLOWED", "RESTRICTED", "PROHIBITED", "ACTION_REQUIRED"
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ActionPoints []string `json:"action_points"`
	KeyStatute   string   `json:"key_statute"`
}

// AssessmentOutput is the generated compliance & strategy roadmap report
type AssessmentOutput struct {
	Summary       string           `json:"summary"`
	ReadinessScore int             `json:"readiness_score"` // 0 to 100
	Steps         []AssessmentStep `json:"steps"`
	Warnings      []string         `json:"warnings,omitempty"`
	GeneratedAt   time.Time        `json:"generated_at"`
}

// IngestSummary reports the result of ingesting documents into PostgreSQL
type IngestSummary struct {
	DocumentID      string `json:"document_id"`
	Title           string `json:"title"`
	TotalChunks     int    `json:"total_chunks"`
	JurisdictionID  string `json:"jurisdiction_id"`
	CategoryID      string `json:"category_id"`
	DurationMs      int64  `json:"duration_ms"`
}

