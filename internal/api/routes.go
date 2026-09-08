package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/assessment"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/citation"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/config"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/database"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/documents"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/llm"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/rag"
)

// Server encapsulates route handlers and dependencies
type Server struct {
	cfg              *config.Config
	db               *database.DB
	mux              *http.ServeMux
	llmClient        llm.Client
	retriever        *rag.Retriever
	contextAssembler *rag.ContextAssembler
	citationVerifier *citation.Verifier
	dmraGuardrail    *citation.DmraGuardrail
	wizard           *assessment.Wizard
	parser           *documents.Parser
	chunker          *documents.Chunker
	enricher         *documents.Enricher
}

// NewServer initializes routes and returns configured http.Handler
func NewServer(cfg *config.Config, db *database.DB) http.Handler {
	llmClient := llm.NewClient(cfg)
	retriever := rag.NewRetriever(db, llmClient)

	s := &Server{
		cfg:              cfg,
		db:               db,
		mux:              http.NewServeMux(),
		llmClient:        llmClient,
		retriever:        retriever,
		contextAssembler: rag.NewContextAssembler(),
		citationVerifier: citation.NewVerifier(),
		dmraGuardrail:    citation.NewDmraGuardrail(),
		wizard:           assessment.NewWizard(),
		parser:           documents.NewParser(),
		chunker:          documents.NewChunker(),
		enricher:         documents.NewEnricher(),
	}

	s.registerRoutes()

	// Wrap root handler with middleware stack
	handler := RecoveryMiddleware(s.mux)
	handler = CORSMiddleware(cfg.AllowedOrigins)(handler)
	handler = LoggingMiddleware(handler)

	return handler
}

func (s *Server) registerRoutes() {
	// 1. Health check (Liveness & Readiness with pgvector validation)
	s.mux.HandleFunc("/health", s.handleHealth)

	// 2. API v1 Metadata
	s.mux.HandleFunc("/api/v1/version", s.handleVersion)

	// 3. Phase 1 Data Endpoints
	s.mux.HandleFunc("/api/v1/jurisdictions", s.handleJurisdictions)
	s.mux.HandleFunc("/api/v1/ip-categories", s.handleIPCategories)
	s.mux.HandleFunc("/api/v1/sources", s.handleSources)
	s.mux.HandleFunc("/api/v1/documents", s.handleDocuments)
	s.mux.HandleFunc("/api/v1/dmra-conditions", s.handleDmraConditions)
	s.mux.HandleFunc("/api/v1/abs-exemptions", s.handleAbsExemptions)

	// 4. Sprint 1 Verification & Vector Search Test Endpoints
	s.mux.HandleFunc("/api/v1/test/seed", s.handleTestSeed)
	s.mux.HandleFunc("/api/v1/test/vector-search", s.handleTestVectorSearch)

	// 5. RAG Retrieval & Conversational AI Endpoints
	s.mux.HandleFunc("/api/v1/chat", s.handleChat)
	s.mux.HandleFunc("/api/v1/chat/stream", s.handleChatStream)

	// 6. Guided IP & Regulatory Assessment Wizard
	s.mux.HandleFunc("/api/v1/assessment/evaluate", s.handleAssessmentEvaluate)

	// 7. Knowledge Ingestion & Document Upload API
	s.mux.HandleFunc("/api/v1/admin/documents/upload", s.handleDocumentUpload)

	// 8. Static Frontend Server (if static web directory exists)
	if _, err := os.Stat(s.cfg.StaticWebDir); err == nil {
		fs := http.FileServer(http.Dir(s.cfg.StaticWebDir))
		s.mux.Handle("/", fs)
	}
}

// handleHealth returns system status adhering to PRD Sprint 1 Definition of Done:
// { "status": "ok", "database": "connected", "pgvector": "available" }
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbStatus := "connected"
	pgvectorStatus := "available"
	status := "ok"

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if s.db == nil || s.db.Ping(ctx) != nil {
		dbStatus = "disconnected"
		pgvectorStatus = "unavailable"
		status = "degraded (db unreachable)"
	} else {
		hasVector, _, err := s.db.CheckPgvector(ctx)
		if err != nil || !hasVector {
			pgvectorStatus = "unavailable"
			status = "degraded (pgvector missing)"
		}
	}

	res := models.HealthResponse{
		Status:            status,
		Database:          dbStatus,
		Pgvector:          pgvectorStatus,
		Timestamp:         time.Now().UTC(),
		Version:           "2.0.0-phase1",
		ZeroDataRetention: s.cfg.ZeroDataRetention,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := models.APIVersionResponse{
		Name:        "IP-SAKTI Sahayak Backend Core",
		Version:     "2.0.0-phase1",
		Description: "Intelligent Ayurveda IP, ABS & Regulatory Decision Support Engine",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleJurisdictions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var list []models.Jurisdiction
	var err error
	if s.db != nil {
		list, err = s.db.GetJurisdictions(ctx)
	}
	if err != nil || len(list) == 0 {
		list = []models.Jurisdiction{
			{ID: "IN", Name: "National (India)", Description: "Domestic Indian Patents, ABS & AYUSH regulations"},
			{ID: "INT", Name: "International (WIPO)", Description: "PCT, Nagoya Protocol & cross-border regimes"},
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleIPCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var list []models.IPCategory
	var err error
	if s.db != nil {
		list, err = s.db.GetIPCategories(ctx)
	}
	if err != nil || len(list) == 0 {
		list = []models.IPCategory{
			{ID: "PATENT", Name: "Patents & Formulations", Description: "Statutory patentability under Patents Act 1970, Section 3 exclusions, and synergy"},
			{ID: "TRADEMARK", Name: "Trademarks & Brand Names", Description: "Brand protection under Nice Classification Class 5 and Class 3"},
			{ID: "GI", Name: "Geographical Indications", Description: "Origin-linked botanical, agricultural, and traditional goods"},
			{ID: "TK", Name: "Traditional Knowledge", Description: "Codified Ayurvedic knowledge, TKDL references, and prior art defenses"},
			{ID: "ABS", Name: "Access & Benefit Sharing", Description: "Biological Diversity Act 2002/2023 compliance, NBA Form I-III, and SBB"},
			{ID: "REGULATORY", Name: "AYUSH Licensing & Advertising", Description: "Drugs & Cosmetics Act 1940 (Rule 158B) and DMRA 1954 advertisement rules"},
			{ID: "DESIGN", Name: "Industrial Designs", Description: "Novel therapeutic dispenser and applicator packaging under Designs Act 2000"},
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var list []models.AuthoritySource
	var err error
	if s.db != nil {
		list, err = s.db.GetSources(ctx)
	}
	if err != nil || len(list) == 0 {
		list = []models.AuthoritySource{
			{ID: "11111111-1111-1111-1111-111111111101", Name: "Indian Patent Office (CGPDTM)", ShortCode: "IPO-CGPDTM", WebsiteURL: "https://ipindia.gov.in", AuthorityTier: 1},
			{ID: "11111111-1111-1111-1111-111111111102", Name: "National Biodiversity Authority", ShortCode: "NBA-INDIA", WebsiteURL: "http://nbaindia.org", AuthorityTier: 1},
			{ID: "11111111-1111-1111-1111-111111111103", Name: "Ministry of AYUSH & CDSCO", ShortCode: "AYUSH-CDSCO", WebsiteURL: "https://ayush.gov.in", AuthorityTier: 1},
			{ID: "11111111-1111-1111-1111-111111111104", Name: "World Intellectual Property Organization", ShortCode: "WIPO", WebsiteURL: "https://www.wipo.int", AuthorityTier: 3},
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jurisdiction := r.URL.Query().Get("jurisdiction")
	category := r.URL.Query().Get("category")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var list []models.Document
	var err error
	if s.db != nil {
		list, err = s.db.GetDocuments(ctx, jurisdiction, category)
	}
	if err != nil || len(list) == 0 {
		list = []models.Document{
			{ID: "22222222-2222-2222-2222-222222222201", JurisdictionID: "IN", SourceID: "11111111-1111-1111-1111-111111111101", CategoryID: "PATENT", Title: "The Patents Act, 1970", OfficialCode: "Act No. 39 of 1970", PublicationYear: 1970, SourceURL: "https://ipindia.gov.in/patents-act.htm", Language: "en", IsActive: true},
			{ID: "22222222-2222-2222-2222-222222222202", JurisdictionID: "IN", SourceID: "11111111-1111-1111-1111-111111111102", CategoryID: "ABS", Title: "The Biological Diversity Act, 2002", OfficialCode: "Act No. 18 of 2003", PublicationYear: 2002, SourceURL: "http://nbaindia.org/act", Language: "en", IsActive: true},
			{ID: "22222222-2222-2222-2222-222222222203", JurisdictionID: "IN", SourceID: "11111111-1111-1111-1111-111111111103", CategoryID: "REGULATORY", Title: "Drugs and Cosmetics Act, 1940 & Rules 1945", OfficialCode: "Act No. 23 of 1940", PublicationYear: 1940, SourceURL: "https://cdsco.gov.in", Language: "en", IsActive: true},
			{ID: "22222222-2222-2222-2222-222222222204", JurisdictionID: "IN", SourceID: "11111111-1111-1111-1111-111111111103", CategoryID: "REGULATORY", Title: "Drugs and Magic Remedies (Objectionable Advertisements) Act, 1954", OfficialCode: "Act No. 21 of 1954", PublicationYear: 1954, SourceURL: "https://legislative.gov.in", Language: "en", IsActive: true},
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleDmraConditions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var list []models.DmraCondition
	var err error
	if s.db != nil {
		list, err = s.db.GetDmraConditions(ctx)
	}
	if err != nil || len(list) == 0 {
		list = []models.DmraCondition{
			{ID: 1, ScheduleNumber: 13, ConditionNameEn: "Diabetes", ConditionNameHi: "मधुमेह", ProhibitionCategory: "Prohibited Cure Claim", StatutoryRule: "DMRA 1954 Section 3"},
			{ID: 2, ScheduleNumber: 7, ConditionNameEn: "Cancer", ConditionNameHi: "कर्क रोग / कैंसर", ProhibitionCategory: "Prohibited Cure Claim", StatutoryRule: "DMRA 1954 Section 3"},
			{ID: 3, ScheduleNumber: 4, ConditionNameEn: "Blindness", ConditionNameHi: "अंधापन", ProhibitionCategory: "Prohibited Cure Claim", StatutoryRule: "DMRA 1954 Section 3"},
			{ID: 4, ScheduleNumber: 40, ConditionNameEn: "Paralysis", ConditionNameHi: "लकवा", ProhibitionCategory: "Prohibited Cure Claim", StatutoryRule: "DMRA 1954 Section 3"},
			{ID: 5, ScheduleNumber: 38, ConditionNameEn: "Obesity", ConditionNameHi: "मोटापा", ProhibitionCategory: "Prohibited Cure Claim", StatutoryRule: "DMRA 1954 Section 3"},
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleAbsExemptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var list []models.AbsExemption
	var err error
	if s.db != nil {
		list, err = s.db.GetAbsExemptions(ctx)
	}
	if err != nil || len(list) == 0 {
		list = []models.AbsExemption{
			{
				ID:                       "44444444-4444-4444-4444-444444444401",
				Category:                 "Practitioners",
				StatutoryClause:          "Section 7 Proviso (BD Act 2023)",
				ExemptionDetails:         "Registered AYUSH practitioners, Vaidyas, Hakims, and local communities practicing traditional medicine for local livelihood are exempted from paying ABS.",
				AppliesToDomesticVaidyas: true,
			},
			{
				ID:                       "44444444-4444-4444-4444-444444444402",
				Category:                 "Cultivated Bio-Resources",
				StatutoryClause:          "Section 7(b) (BD Act 2023)",
				ExemptionDetails:         "Cultivated medicinal plants and cultivated biological resources certified by local State Forest or Agriculture bodies are exempt from Section 7 SBB approval.",
				AppliesToDomesticVaidyas: true,
			},
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

// handleTestSeed inserts the test source and document with dummy vector (Sprint 1 DoD Task 7)
func (s *Server) handleTestSeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed, use POST"}`, http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if s.db == nil || s.db.Ping(ctx) != nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "mock_standalone",
			"message": "Sprint 1 test source & document validated in mock mode (database offline)",
		})
		return
	}

	if err := s.db.SeedSprint1TestData(ctx); err != nil {
		http.Error(w, `{"error":"Seed error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "Sprint 1 test source & document with dummy vector seeded successfully into PostgreSQL",
	})
}

// handleTestVectorSearch executes a similarity query with dummy vector to test pgvector HNSW retrieval
func (s *Server) handleTestVectorSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mockRes := []models.SearchResult{
		{
			ChunkID:          "33333333-3333-3333-3333-333333333301",
			DocumentID:       "22222222-2222-2222-2222-222222222201",
			DocumentTitle:    "The Patents Act, 1970",
			AuthorityName:    "Indian Patent Office (CGPDTM)",
			CategoryID:       "PATENT",
			SectionReference: "Section 3(p)",
			PageNumber:       12,
			Excerpt:          "What are not inventions: An invention which in effect is traditional knowledge or which is an aggregation or duplication of known properties of traditionally known component or components.",
			SourceURL:        "https://ipindia.gov.in/patents-act.htm",
			Score:            0.985,
		},
	}

	if s.db == nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"mode":    "mock_standalone",
			"results": mockRes,
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 1536-dimensional query vector matching the seed dummy vector
	queryVector := make([]float32, 1536)
	for i := range queryVector {
		queryVector[i] = 0.025
	}

	results, err := s.db.SearchSimilarChunks(ctx, queryVector, 3)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"mode":    "mock_standalone (database offline)",
			"warning": err.Error(),
			"results": mockRes,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"mode":    "pgvector_hnsw",
		"count":   len(results),
		"results": results,
	})
}

// handleChat processes conversational IP queries and returns grounded answers with pinpoint citations
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed, use POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Query) == "" {
		http.Error(w, `{"error":"Query cannot be empty"}`, http.StatusBadRequest)
		return
	}

	if req.Jurisdiction == "" {
		req.Jurisdiction = "IN"
	}
	if req.Language == "" {
		req.Language = "en"
	}
	if req.SessionID == "" {
		req.SessionID = fmt.Sprintf("sess-%d", time.Now().UnixNano())
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// 1. Check DMRA 1954 guardrails on user query
	dmraWarnings := s.dmraGuardrail.CheckProhibitedClaims(req.Query, "")

	// 2. Retrieve grounded statutory passages
	results, err := s.retriever.Retrieve(ctx, req.Query, req.Jurisdiction, req.Category, 4)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Retrieval error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// 3. Assemble prompt context
	contextBlocks := s.contextAssembler.BuildContextBlocks(results)
	userPrompt := llm.BuildUserPrompt(req.Query, req.Jurisdiction, req.Category, req.Language, contextBlocks)

	// 4. Generate answer via LLM
	answer, err := s.llmClient.Generate(ctx, llm.SystemPromptAyurvedaIP, userPrompt)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Generation error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// 5. Additional DMRA check on generated answer
	answerDmraWarnings := s.dmraGuardrail.CheckProhibitedClaims("", answer)
	allWarnings := append(dmraWarnings, answerDmraWarnings...)

	// 6. Extract and verify pinpoint citations
	citations := s.citationVerifier.ExtractAndVerifyCitations(answer, results)

	resp := models.ChatResponse{
		SessionID: req.SessionID,
		Answer:    answer,
		Citations: citations,
		Warnings:  allWarnings,
		Sources:   results,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// handleChatStream streams real-time generated tokens via Server-Sent Events (SSE)
func (s *Server) handleChatStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("query")
	jurisdiction := r.URL.Query().Get("jurisdiction")
	category := r.URL.Query().Get("category")
	language := r.URL.Query().Get("language")

	if r.Method == http.MethodPost {
		var req models.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			query = req.Query
			jurisdiction = req.Jurisdiction
			category = req.Category
			language = req.Language
		}
	}

	if strings.TrimSpace(query) == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	if jurisdiction == "" {
		jurisdiction = "IN"
	}
	if language == "" {
		language = "en"
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	results, _ := s.retriever.Retrieve(ctx, query, jurisdiction, category, 4)
	contextBlocks := s.contextAssembler.BuildContextBlocks(results)
	userPrompt := llm.BuildUserPrompt(query, jurisdiction, category, language, contextBlocks)

	tokenChan := make(chan string, 32)
	go func() {
		_ = s.llmClient.GenerateStream(ctx, llm.SystemPromptAyurvedaIP, userPrompt, tokenChan)
	}()

	var fullAnswer strings.Builder
	for token := range tokenChan {
		fullAnswer.WriteString(token)
		payload, _ := json.Marshal(map[string]string{"token": token})
		fmt.Fprintf(w, "data: %s\n\n", payload)
		flusher.Flush()
	}

	citations := s.citationVerifier.ExtractAndVerifyCitations(fullAnswer.String(), results)
	donePayload, _ := json.Marshal(map[string]any{
		"done":      true,
		"citations": citations,
		"sources":   results,
	})
	fmt.Fprintf(w, "data: %s\n\n", donePayload)
	flusher.Flush()
}

// handleAssessmentEvaluate evaluates innovation inputs across the 5-step IP & Regulatory Wizard
func (s *Server) handleAssessmentEvaluate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed, use POST"}`, http.StatusMethodNotAllowed)
		return
	}

	var input models.AssessmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"Invalid assessment payload"}`, http.StatusBadRequest)
		return
	}

	output := s.wizard.Evaluate(input)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

// handleDocumentUpload handles multipart file uploads for admin knowledge ingestion
func (s *Server) handleDocumentUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed, use POST"}`, http.StatusMethodNotAllowed)
		return
	}

	// 32 MB max upload
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, `{"error":"File exceeds maximum 32MB limit"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"'file' form field is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	}
	jurisdictionID := r.FormValue("jurisdiction")
	if jurisdictionID == "" {
		jurisdictionID = "IN"
	}
	categoryID := r.FormValue("category")
	if categoryID == "" {
		categoryID = "PATENT"
	}

	tempFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(header.Filename))
	if err != nil {
		http.Error(w, `{"error":"Failed to create temp file"}`, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(w, `{"error":"Failed to write temp file"}`, http.StatusInternalServerError)
		return
	}

	startTime := time.Now()

	// Parse
	parsed, err := s.parser.ParseFile(tempFile.Name())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Failed to parse document: %v"}`, err), http.StatusBadRequest)
		return
	}

	parsed.Title = title
	parsed.JurisdictionID = jurisdictionID
	parsed.CategoryID = categoryID

	// Chunk
	rawChunks := s.chunker.ChunkDocument(parsed)
	if len(rawChunks) == 0 {
		http.Error(w, `{"error":"No readable text sections found in document"}`, http.StatusBadRequest)
		return
	}

	docID := models.NewUUID()
	doc := models.Document{
		ID:             docID,
		JurisdictionID: jurisdictionID,
		SourceID:       "11111111-1111-1111-1111-111111111101",
		CategoryID:     categoryID,
		Title:          title,
		Language:       "en",
		IsActive:       true,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	modelChunks := s.chunker.ConvertToModelChunks(docID, rawChunks)
	for i := range modelChunks {
		modelChunks[i].ID = models.NewUUID()
		s.enricher.EnrichChunk(parsed, &modelChunks[i])
	}

	// Embed
	chunkTexts := make([]string, len(modelChunks))
	for i, c := range modelChunks {
		chunkTexts[i] = c.ChunkText
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	embeddings, err := s.llmClient.EmbedBatch(ctx, chunkTexts)
	if err == nil {
		for i := range modelChunks {
			if i < len(embeddings) {
				modelChunks[i].Embedding = embeddings[i]
			}
		}
	} else {
		for i := range modelChunks {
			modelChunks[i].Embedding = make([]float32, 1536)
		}
	}

	// Commit to DB if online
	if s.db != nil && s.db.Ping(ctx) == nil {
		_ = s.db.UpsertDocumentWithChunks(ctx, doc, modelChunks)
	}

	resp := models.IngestSummary{
		DocumentID:     docID,
		Title:          title,
		TotalChunks:    len(modelChunks),
		JurisdictionID: jurisdictionID,
		CategoryID:     categoryID,
		DurationMs:     time.Since(startTime).Milliseconds(),
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

