package citation

import (
	"fmt"
	"strings"
)

// DmraGuardrail scans for prohibited disease cure claims under DMRA 1954 Schedule
type DmraGuardrail struct {
	prohibitedDiseases map[string]string
}

// NewDmraGuardrail initializes the 54 scheduled disease scanner
func NewDmraGuardrail() *DmraGuardrail {
	diseases := map[string]string{
		"appendicitis":        "Appendicitis (Schedule Item 1)",
		"arteriosclerosis":    "Arteriosclerosis (Schedule Item 2)",
		"blindness":           "Blindness (Schedule Item 3)",
		"blood poisoning":     "Blood poisoning (Schedule Item 4)",
		"bright's disease":    "Bright's disease (Schedule Item 5)",
		"cancer":              "Cancer (Schedule Item 6)",
		"cataract":            "Cataract (Schedule Item 7)",
		"deafness":            "Deafness (Schedule Item 8)",
		"diabetes":            "Diabetes (Schedule Item 9)",
		"brain disease":       "Diseases and disorders of brain (Schedule Item 10)",
		"optical disease":     "Diseases and disorders of the optical system (Schedule Item 11)",
		"uterus disease":      "Diseases and disorders of the uterus (Schedule Item 12)",
		"menstrual disorder":  "Disorders of menstrual flow (Schedule Item 13)",
		"nervous disorder":    "Disorders of the nervous system (Schedule Item 14)",
		"dropsy":              "Dropsy (Schedule Item 15)",
		"epilepsy":            "Epilepsy (Schedule Item 16)",
		"female disease":      "Female diseases in general (Schedule Item 17)",
		"fever":               "Fevers in general (Schedule Item 18)",
		"fits":                "Fits (Schedule Item 19)",
		"female bust":         "Form and structure of the female bust (Schedule Item 20)",
		"gall stone":          "Gall stones, kidney stones and bladder stones (Schedule Item 21)",
		"kidney stone":        "Gall stones, kidney stones and bladder stones (Schedule Item 21)",
		"bladder stone":       "Gall stones, kidney stones and bladder stones (Schedule Item 21)",
		"gangrene":            "Gangrene (Schedule Item 22)",
		"glaucoma":            "Glaucoma (Schedule Item 23)",
		"goitre":              "Goitre (Schedule Item 24)",
		"heart disease":       "Heart diseases (Schedule Item 25)",
		"high blood pressure": "High/low blood pressure (Schedule Item 26)",
		"low blood pressure":  "High/low blood pressure (Schedule Item 26)",
		"hydrocele":           "Hydrocele (Schedule Item 27)",
		"hysteria":            "Hysteria (Schedule Item 28)",
		"infantile paralysis": "Infantile paralysis (Schedule Item 29)",
		"insanity":            "Insanity (Schedule Item 30)",
		"leprosy":             "Leprosy (Schedule Item 31)",
		"leucoderma":          "Leucoderma (Schedule Item 32)",
		"lockjaw":             "Lockjaw (Schedule Item 33)",
		"locomotor ataxia":    "Locomotor ataxia (Schedule Item 34)",
		"lupus":               "Lupus (Schedule Item 35)",
		"nervous debility":    "Nervous debility (Schedule Item 36)",
		"obesity":             "Obesity (Schedule Item 37)",
		"paralysis":           "Paralysis (Schedule Item 38)",
		"plague":              "Plague (Schedule Item 39)",
		"pleurisy":            "Pleurisy (Schedule Item 40)",
		"pneumonia":           "Pneumonia (Schedule Item 41)",
		"rheumatism":          "Rheumatism (Schedule Item 42)",
		"rupture":             "Ruptures (Schedule Item 43)",
		"sexual impotence":    "Sexual impotence (Schedule Item 44)",
		"smallpox":            "Smallpox (Schedule Item 45)",
		"stature":             "Stature of persons (Schedule Item 46)",
		"sterility":           "Sterility in women (Schedule Item 47)",
		"trachoma":            "Trachoma (Schedule Item 48)",
		"tuberculosis":        "Tuberculosis (Schedule Item 49)",
		"tumour":              "Tumours (Schedule Item 50)",
		"typhoid":             "Typhoid fever (Schedule Item 51)",
		"ulcer":               "Ulcers of the gastro-intestinal tract (Schedule Item 52)",
		"venereal disease":    "Venereal diseases (Schedule Item 53)",
		"baldness":            "Baldness / Hair loss (Schedule Item 54)",
		"madhumeha":           "Diabetes / Madhumeha (Schedule Item 9)",
		"kark rog":            "Cancer / Kark Rog (Schedule Item 6)",
	}

	return &DmraGuardrail{prohibitedDiseases: diseases}
}

// CheckProhibitedClaims scans text and returns statutory warnings if prohibited claims are found
func (dg *DmraGuardrail) CheckProhibitedClaims(query, answer string) []string {
	combined := strings.ToLower(query + " " + answer)
	var warnings []string

	// Check if cure or treatment is promised
	hasCureIntent := strings.Contains(combined, "cure") ||
		strings.Contains(combined, "treat") ||
		strings.Contains(combined, "heal") ||
		strings.Contains(combined, "remedy") ||
		strings.Contains(combined, "claim") ||
		strings.Contains(combined, "advertise")

	for keyword, scheduleRef := range dg.prohibitedDiseases {
		if strings.Contains(combined, keyword) {
			if hasCureIntent {
				warning := fmt.Sprintf(
					"⚠️ DMRA 1954 STATUTORY WARNING: Advertising or claiming a 'cure' for %s is strictly prohibited under Section 3 of the Drugs and Magic Remedies (Objectionable Advertisements) Act, 1954. Violations are cognizable offenses punishable by imprisonment under Section 7.",
					scheduleRef,
				)
				warnings = append(warnings, warning)
			}
		}
	}

	return warnings
}
