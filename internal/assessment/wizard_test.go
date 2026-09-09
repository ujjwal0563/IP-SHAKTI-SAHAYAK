package assessment

import (
	"testing"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

func TestWizardClassicalFormulation(t *testing.T) {
	wizard := NewWizard()

	input := models.AssessmentInput{
		InnovationType:        "classical_ayurveda",
		UsesIndianBioResource: true,
		TargetMarket:          "domestic_only",
		IntendsToPatent:       true,
		HasDiseaseClaims:      false,
	}

	res := wizard.Evaluate(input)

	if len(res.Steps) != 4 {
		t.Fatalf("expected 4 assessment steps, got %d", len(res.Steps))
	}

	for _, s := range res.Steps {
		if s.Domain == "PATENT" {
			if s.Status != "PROHIBITED" {
				t.Errorf("expected classical formulation to be PROHIBITED for patenting, got %s", s.Status)
			}
		}
		if s.Domain == "ABS_BIODIVERSITY" {
			if s.Status != "ACTION_REQUIRED" {
				t.Errorf("expected NBA Form III action required for bio resource patenting, got %s", s.Status)
			}
		}
	}
}

func TestWizardDeviceAllowed(t *testing.T) {
	wizard := NewWizard()

	input := models.AssessmentInput{
		InnovationType:        "device",
		UsesIndianBioResource: false,
		TargetMarket:          "global",
		IntendsToPatent:       true,
		HasDiseaseClaims:      false,
	}

	res := wizard.Evaluate(input)

	for _, s := range res.Steps {
		if s.Domain == "PATENT" {
			if s.Status != "ALLOWED" {
				t.Errorf("expected device patenting to be ALLOWED, got %s", s.Status)
			}
		}
	}
}
