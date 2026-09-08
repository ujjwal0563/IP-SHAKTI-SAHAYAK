package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/config"
	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:              "8080",
		AllowedOrigins:    "*",
		ZeroDataRetention: true,
	}

	handler := NewServer(cfg, nil) // standalone mode with nil DB

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res models.HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if res.Database != "disconnected" {
		t.Errorf("expected database 'disconnected' in standalone mode, got %s", res.Database)
	}

	if res.Pgvector != "unavailable" {
		t.Errorf("expected pgvector 'unavailable' in standalone mode, got %s", res.Pgvector)
	}
}

func TestIPCategoriesEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		AllowedOrigins: "*",
	}

	handler := NewServer(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ip-categories", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var categories []models.IPCategory
	if err := json.NewDecoder(rec.Body).Decode(&categories); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if len(categories) == 0 {
		t.Error("expected fallback IP categories, got 0")
	}

	// Verify key statutory categories exist
	foundPatent := false
	foundABS := false
	for _, c := range categories {
		if c.ID == "PATENT" {
			foundPatent = true
		}
		if c.ID == "ABS" {
			foundABS = true
		}
	}

	if !foundPatent || !foundABS {
		t.Errorf("expected PATENT and ABS categories, got %+v", categories)
	}
}

func TestSourcesEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		AllowedOrigins: "*",
	}

	handler := NewServer(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sources", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var sources []models.AuthoritySource
	if err := json.NewDecoder(rec.Body).Decode(&sources); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if len(sources) == 0 {
		t.Error("expected fallback sources, got 0")
	}
}

func TestDmraConditionsEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		AllowedOrigins: "*",
	}

	handler := NewServer(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dmra-conditions", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var conditions []models.DmraCondition
	if err := json.NewDecoder(rec.Body).Decode(&conditions); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if len(conditions) == 0 {
		t.Error("expected fallback DMRA conditions, got 0")
	}
}

func TestAbsExemptionsEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		AllowedOrigins: "*",
	}

	handler := NewServer(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/abs-exemptions", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var exemptions []models.AbsExemption
	if err := json.NewDecoder(rec.Body).Decode(&exemptions); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if len(exemptions) == 0 {
		t.Error("expected fallback ABS exemptions, got 0")
	}
}

func TestChatEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		AllowedOrigins: "*",
		LLMProvider:    "mock",
	}

	handler := NewServer(cfg, nil)

	payload := `{"query":"Can I patent an Ayurvedic herbal formulation with Ashwagandha?","jurisdiction":"IN","category":"PATENT"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var chatResp models.ChatResponse
	if err := json.NewDecoder(rec.Body).Decode(&chatResp); err != nil {
		t.Fatalf("failed to decode ChatResponse: %v", err)
	}

	if chatResp.Answer == "" {
		t.Error("expected non-empty answer")
	}

	if len(chatResp.Citations) == 0 {
		t.Error("expected citations in chat response, got 0")
	}

	// Verify Section 3(p) patent bar warning is mentioned in citations or answer
	hasPatentsRef := strings.Contains(chatResp.Answer, "3(p)") || strings.Contains(chatResp.Answer, "Patents Act")
	if !hasPatentsRef {
		t.Errorf("expected Section 3(p) reference in answer, got: %s", chatResp.Answer)
	}
}

func TestAssessmentEvaluateEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:           "8080",
		AllowedOrigins: "*",
	}

	handler := NewServer(cfg, nil)

	payload := `{
		"innovation_type": "classical_ayurveda",
		"uses_indian_bio_resource": true,
		"target_market": "domestic_only",
		"intends_to_patent": true,
		"has_disease_claims": true,
		"claimed_diseases": ["diabetes"]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/assessment/evaluate", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var assessment models.AssessmentOutput
	if err := json.NewDecoder(rec.Body).Decode(&assessment); err != nil {
		t.Fatalf("failed to decode AssessmentOutput: %v", err)
	}

	if len(assessment.Steps) != 4 {
		t.Errorf("expected 4 assessment steps, got %d", len(assessment.Steps))
	}

	// Classical formulation + patent intent must result in PROHIBITED under Section 3(p)
	var patentStep *models.AssessmentStep
	for _, s := range assessment.Steps {
		if s.Domain == "PATENT" {
			patentStep = &s
			break
		}
	}

	if patentStep == nil || patentStep.Status != "PROHIBITED" {
		t.Errorf("expected classical formulation patenting to be PROHIBITED under 3(p), got %+v", patentStep)
	}

	// Diabetes cure claim must trigger DMRA warning
	if len(assessment.Warnings) == 0 {
		t.Error("expected DMRA warning for diabetes cure claim, got none")
	}
}

