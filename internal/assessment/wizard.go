package assessment

import (
	"strings"
	"time"

	"github.com/ujjwal0563/ip-shakti-sahayak/internal/models"
)

// Wizard evaluates user questionnaire inputs across Indian IP, ABS, and AYUSH regulatory domains
type Wizard struct{}

// NewWizard creates a new assessment wizard
func NewWizard() *Wizard {
	return &Wizard{}
}

// Evaluate runs the 5-step decision tree logic and returns a comprehensive readiness report
func (w *Wizard) Evaluate(in models.AssessmentInput) models.AssessmentOutput {
	var steps []models.AssessmentStep
	var warnings []string
	score := 100

	// Step 1: Patentability Assessment (Section 3(p) & 3(e))
	patentStep := models.AssessmentStep{
		Domain:     "PATENT",
		KeyStatute: "Patents Act, 1970 (Section 3(p), 3(e), 3(d))",
	}

	if in.IntendsToPatent {
		switch in.InnovationType {
		case "classical_ayurveda":
			patentStep.Status = "PROHIBITED"
			patentStep.Title = "Direct Patenting Bar under Section 3(p)"
			patentStep.Description = "Inventions that directly replicate classical Ayurvedic texts (e.g. Charaka Samhita) are non-patentable as Traditional Knowledge."
			patentStep.ActionPoints = []string{
				"Do NOT file a standard formulation patent for classical recipes; it will be rejected by the Patent Office using TKDL citations.",
				"Consider protecting the brand name via Trademark (Nice Class 5) or packaging via Industrial Design Act 2000.",
				"Focus on patenting novel extraction equipment or delivery mechanisms rather than the formulation itself.",
			}
			score -= 30
		case "novel_extract", "synergistic_combo":
			patentStep.Status = "RESTRICTED"
			patentStep.Title = "Patentable with Synergistic Proof (Section 3(e))"
			patentStep.Description = "Formulations combining known herbs can only be patented if you demonstrate synergistic efficacy exceeding the additive sum of individual components."
			patentStep.ActionPoints = []string{
				"Conduct comparative bio-activity or in-vitro/in-vivo synergy studies to rebut Section 3(e) mere admixture objections.",
				"Explicitly draft claims around standardized fraction markers and specific ratio ranges rather than broad botanical names.",
				"Ensure prior biological resource clearance under the Biological Diversity Act before patent grant.",
			}
			score -= 10
		case "device":
			patentStep.Status = "ALLOWED"
			patentStep.Title = "Patentable Mechanical / Therapeutic Device"
			patentStep.Description = "Ayurvedic therapy equipment (e.g. Panchakarma dispensers, steam chambers) are free from Section 3(p) botanical bars."
			patentStep.ActionPoints = []string{
				"Focus claims on mechanical novelties, automated dosage controls, and temperature regulators.",
				"Consider concurrent Industrial Design filing under Designs Act 2000 for aesthetics.",
			}
		default:
			patentStep.Status = "ACTION_REQUIRED"
			patentStep.Title = "Patentability Evaluation Required"
			patentStep.Description = "Further evaluation required based on specific chemical and formulation details."
			patentStep.ActionPoints = []string{"Consult with a registered Patent Agent."}
		}
	} else {
		patentStep.Status = "ALLOWED"
		patentStep.Title = "Patent Protection Not Sought"
		patentStep.Description = "Innovator is relying on Trade Secrets, Trademarks, and Brand Protection."
		patentStep.ActionPoints = []string{
			"Protect product brand name under Trademark Nice Class 5 (Pharmaceuticals) and Class 3 (Cosmetics).",
			"Execute strict Non-Disclosure Agreements (NDAs) with manufacturing partners.",
		}
	}
	steps = append(steps, patentStep)

	// Step 2: Biological Diversity & ABS Compliance (BD Act 2023)
	absStep := models.AssessmentStep{
		Domain:     "ABS_BIODIVERSITY",
		KeyStatute: "Biological Diversity Act, 2023 (Section 3, 6, 7 & Rules 2024)",
	}

	if in.UsesIndianBioResource {
		if in.IntendsToPatent {
			absStep.Status = "ACTION_REQUIRED"
			absStep.Title = "Mandatory NBA Form III Approval"
			absStep.Description = "Under Section 6 of the Biological Diversity Act, prior approval from the National Biodiversity Authority is MANDATORY before sealing any patent based on Indian bio-resources."
			absStep.ActionPoints = []string{
				"File NBA Form III before patent grant to avoid invalidation or criminal penalties under Section 55.",
				"Maintain traceability records of whether biological resources were sourced from certified cultivators or wild collection.",
			}
			score -= 15
		} else {
			absStep.Status = "RESTRICTED"
			absStep.Title = "State Biodiversity Board (SBB) Intimation"
			absStep.Description = "Commercial utilization of biological resources requires prior intimation to the concerned State Biodiversity Board under Section 7."
			absStep.ActionPoints = []string{
				"Check if your practice qualifies under the Section 7 Proviso exemption (Registered local Vaidyas and Hakims for local healthcare).",
				"If purchasing commercial raw materials, ensure suppliers possess valid SBB access permits.",
			}
		}
	} else {
		absStep.Status = "ALLOWED"
		absStep.Title = "No Indian Biological Resources Utilized"
		absStep.Description = "Biological Diversity Act ABS requirements are not triggered."
		absStep.ActionPoints = []string{"Maintain import certificates or synthesis proof to verify non-Indian sourcing."}
	}
	steps = append(steps, absStep)

	// Step 3: AYUSH Regulatory & Licensing Pathway
	regStep := models.AssessmentStep{
		Domain:     "AYUSH_REGULATORY",
		KeyStatute: "Drugs & Cosmetics Act, 1940 (Chapter IV-A) & Rules 1945 (Rule 158B)",
	}

	if in.InnovationType == "classical_ayurveda" {
		regStep.Status = "ALLOWED"
		regStep.Title = "Classical Ayurvedic Medicine License"
		regStep.Description = "Eligible for streamlined licensing under Schedule I recognized Ayurvedic texts without clinical trials."
		regStep.ActionPoints = []string{
			"Manufacture strictly according to textual formulations cited in the First Schedule of the Drugs & Cosmetics Act.",
			"Obtain State AYUSH Licensing Authority (SLA) manufacturing approval.",
		}
	} else if in.InnovationType == "novel_extract" || in.InnovationType == "synergistic_combo" {
		regStep.Status = "ACTION_REQUIRED"
		regStep.Title = "Patent & Proprietary (P&P) Medicine (Rule 158B)"
		regStep.Description = "Requires safety and efficacy proof under Rule 158B of the Drugs and Cosmetics Rules, 1945."
		regStep.ActionPoints = []string{
			"Conduct pilot clinical studies or published safety documentation as required by Rule 158B.",
			"Submit formulation recipe, shelf-life stability data, and heavy-metal testing certificates to State Licensing Authority.",
		}
		score -= 15
	} else {
		regStep.Status = "ALLOWED"
		regStep.Title = "Medical Device / Applicator Standards"
		regStep.Description = "Comply with Medical Devices Rules, 2017 standards."
		regStep.ActionPoints = []string{"Register with CDSCO under applicable Medical Device Classification."}
	}
	steps = append(steps, regStep)

	// Step 4: DMRA 1954 Advertising Compliance
	dmraStep := models.AssessmentStep{
		Domain:     "DMRA",
		KeyStatute: "Drugs and Magic Remedies (Objectionable Advertisements) Act, 1954",
	}

	if in.HasDiseaseClaims || len(in.ClaimedDiseases) > 0 {
		prohibitedFound := false
		for _, d := range in.ClaimedDiseases {
			lower := strings.ToLower(d)
			if strings.Contains(lower, "diabetes") || strings.Contains(lower, "cancer") ||
				strings.Contains(lower, "obesity") || strings.Contains(lower, "paralysis") ||
				strings.Contains(lower, "cure") {
				prohibitedFound = true
				break
			}
		}

		if prohibitedFound {
			dmraStep.Status = "PROHIBITED"
			dmraStep.Title = "Severe DMRA Violation Risk"
			dmraStep.Description = "Claiming a 'cure' for scheduled diseases (Diabetes, Cancer, Obesity, etc.) is a criminal offense under Section 3 of DMRA 1954."
			dmraStep.ActionPoints = []string{
				"Immediately remove all 'cure' or 'reversal' claims from marketing materials, packaging, and digital ads.",
				"Reframe claims to supportive lifestyle or wellness terminology permitted under AYUSH guidelines (e.g. 'Supports healthy blood sugar metabolism').",
			}
			warnings = append(warnings, "DMRA 1954 Warning: Curing claims for scheduled diseases carry imprisonment penalties under Section 7.")
			score -= 30
		} else {
			dmraStep.Status = "RESTRICTED"
			dmraStep.Title = "Marketing Claims Require AYUSH Review"
			dmraStep.Description = "Promotional claims must remain within therapeutic indications approved on your SLA product label."
			dmraStep.ActionPoints = []string{"Ensure packaging text matches the exact indications approved in your license."}
			score -= 5
		}
	} else {
		dmraStep.Status = "ALLOWED"
		dmraStep.Title = "Low DMRA Advertising Risk"
		dmraStep.Description = "No prohibited scheduled disease cure claims reported."
		dmraStep.ActionPoints = []string{"Maintain factual, evidence-backed labeling."}
	}
	steps = append(steps, dmraStep)

	// Ensure score is clamped between 0 and 100
	if score < 0 {
		score = 0
	}

	return models.AssessmentOutput{
		Summary:        "Comprehensive multi-step assessment across Indian Patent, ABS Biodiversity, and AYUSH drug regulatory statutes.",
		ReadinessScore: score,
		Steps:          steps,
		Warnings:       warnings,
		GeneratedAt:    time.Now().UTC(),
	}
}
