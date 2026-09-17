package rag

import (
	"strings"
)

// LegalTermConcept represents a bilingual statutory legal concept
type LegalTermConcept struct {
	EnglishTerm string   `json:"english_term"`
	HindiTerm   string   `json:"hindi_term"`
	Statute     string   `json:"statute"`
	Category    string   `json:"category"`
	Synonyms    []string `json:"synonyms"`
	Description string   `json:"description"`
}

// AyurvedicFormulationCategory represents the statutory classification of traditional formulations
type AyurvedicFormulationCategory struct {
	FormName          string   `json:"form_name"`
	HindiName         string   `json:"hindi_name"`
	StatutoryCategory string   `json:"statutory_category"` // e.g., "Classical (First Schedule)", "Patent & Proprietary (Rule 158B)"
	ApplicableAct     string   `json:"applicable_act"`
	PatentStatus      string   `json:"patent_status"` // "Barred under Sec 3(p)", "Requires Sec 3(e) synergy", etc.
	Synonyms          []string `json:"synonyms"`
}

// BilingualGlossary manages legal terminology, Ayurvedic taxonomy, and dual-term rendering
type BilingualGlossary struct {
	concepts     map[string]LegalTermConcept
	formulations map[string]AyurvedicFormulationCategory
}

// NewBilingualGlossary initializes the domain glossary with authoritative statutory definitions
func NewBilingualGlossary() *BilingualGlossary {
	bg := &BilingualGlossary{
		concepts:     make(map[string]LegalTermConcept),
		formulations: make(map[string]AyurvedicFormulationCategory),
	}

	bg.initConcepts()
	bg.initFormulations()

	return bg
}

func (bg *BilingualGlossary) initConcepts() {
	list := []LegalTermConcept{
		{
			EnglishTerm: "Patent",
			HindiTerm:   "एकाधिकार / पेटेंट (Patent)",
			Statute:     "The Patents Act, 1970",
			Category:    "PATENT",
			Synonyms:    []string{"patent", "patents", "पेटेंट", "एकाधिकार"},
			Description: "Statutory monopoly grant for novel, non-obvious, industrially applicable inventions.",
		},
		{
			EnglishTerm: "Traditional Knowledge",
			HindiTerm:   "पारंपरिक ज्ञान (Traditional Knowledge)",
			Statute:     "The Patents Act, 1970 (Section 3(p))",
			Category:    "PATENT",
			Synonyms:    []string{"traditional knowledge", "tk", "पारंपरिक ज्ञान", "प्राचीन ज्ञान", "classical knowledge"},
			Description: "Inventions replicating known properties of traditional materials are non-patentable under Section 3(p).",
		},
		{
			EnglishTerm: "Synergistic Effect / Non-Obvious Admixture",
			HindiTerm:   "सहक्रियात्मक प्रभाव (Synergistic Effect)",
			Statute:     "The Patents Act, 1970 (Section 3(e))",
			Category:    "PATENT",
			Synonyms:    []string{"synergy", "synergistic", "mere admixture", "सहक्रियात्मक", "सम्मिश्रण", "admixture"},
			Description: "Formulations combining known components must prove efficacy exceeding the additive sum under Section 3(e).",
		},
		{
			EnglishTerm: "Prior Art",
			HindiTerm:   "पूर्व ज्ञान / पूर्व कला (Prior Art)",
			Statute:     "The Patents Act, 1970 (Section 2(1)(j)) / TKDL",
			Category:    "PATENT",
			Synonyms:    []string{"prior art", "पूर्व ज्ञान", "पूर्व कला", "tkdl"},
			Description: "Publicly available evidence before filing date, including ancient texts documented in TKDL.",
		},
		{
			EnglishTerm: "Biological Diversity & ABS",
			HindiTerm:   "जैविक विविधता एवं पहुंच और लाभ साझाकरण (ABS)",
			Statute:     "Biological Diversity Act, 2023 (Section 3, 6, 7)",
			Category:    "ABS",
			Synonyms:    []string{"abs", "biodiversity", "biological resource", "nba", "sbb", "जैव विविधता", "लाभ साझाकरण"},
			Description: "Access and Benefit Sharing duty on commercial exploitation of Indian biological resources.",
		},
		{
			EnglishTerm: "National Biodiversity Authority (NBA)",
			HindiTerm:   "राष्ट्रीय जैव विविधता प्राधिकरण (NBA)",
			Statute:     "Biological Diversity Act, 2023 (Section 8, 19, 20)",
			Category:    "ABS",
			Synonyms:    []string{"nba", "national biodiversity authority", "राष्ट्रीय जैव विविधता प्राधिकरण"},
			Description: "National statutory authority mandating prior Form III approval for patents and Form I for commercial access.",
		},
		{
			EnglishTerm: "State Biodiversity Board (SBB)",
			HindiTerm:   "राज्य जैव विविधता बोर्ड (SBB)",
			Statute:     "Biological Diversity Act, 2023 (Section 7, 23)",
			Category:    "ABS",
			Synonyms:    []string{"sbb", "state biodiversity board", "राज्य जैव विविधता बोर्ड"},
			Description: "State body requiring prior commercial intimation for access to Indian bio-resources.",
		},
		{
			EnglishTerm: "Registered Practitioner Exemption",
			HindiTerm:   "पंजीकृत वैद्य/चिकित्सक छूट (Section 7 Proviso)",
			Statute:     "Biological Diversity Act, 2023 (Section 7 Proviso)",
			Category:    "ABS",
			Synonyms:    []string{"vaidya exemption", "hakim", "vaid", "स्थानीय वैद्य", "चिकित्सक छूट", "proviso exemption"},
			Description: "Registered local Vaidyas and Hakims practicing traditional healthcare are exempted from ABS intimation and fee.",
		},
		{
			EnglishTerm: "Classical Ayurvedic Medicine",
			HindiTerm:   "शास्त्रीय आयुर्वेदिक औषधि (Classical Medicine)",
			Statute:     "Drugs and Cosmetics Act, 1940 (First Schedule)",
			Category:    "REGULATORY",
			Synonyms:    []string{"classical", "shastriya", "शास्त्रीय", "first schedule", "अनुसूची १"},
			Description: "Medicines manufactured strictly according to recipes in 54 recognized classical Ayurvedic treatises.",
		},
		{
			EnglishTerm: "Patent or Proprietary (P&P) Medicine",
			HindiTerm:   "स्वामित्व / पेटेंट एवं प्रोप्राइटरी औषधि (P&P Medicine)",
			Statute:     "Drugs and Cosmetics Act, 1940 (Rule 158B)",
			Category:    "REGULATORY",
			Synonyms:    []string{"p&p", "proprietary", "rule 158b", "प्रोप्राइटरी", "स्वामित्व औषधि"},
			Description: "Novel or modified Ayurvedic formulations requiring safety and efficacy proof under Rule 158B.",
		},
		{
			EnglishTerm: "Objectionable Advertisement Prohibition",
			HindiTerm:   "आपत्तिजनक विज्ञापन प्रतिषेध (DMRA 1954)",
			Statute:     "Drugs and Magic Remedies (Objectionable Advertisements) Act, 1954",
			Category:    "DMRA",
			Synonyms:    []string{"dmra", "cure claim", "magic remedy", "advertisement", "आपत्तिजनक विज्ञापन", "शर्तिया इलाज", "रामबाण"},
			Description: "Strict criminal prohibition against claiming remedies or cures for 54 scheduled diseases (Section 3).",
		},
		{
			EnglishTerm: "Ayurveda Aahar (Nutraceutical)",
			HindiTerm:   "आयुर्वेद आहार (Ayurveda Aahar)",
			Statute:     "FSSAI (Ayurveda Aahar) Regulations, 2022",
			Category:    "REGULATORY",
			Synonyms:    []string{"ayurveda aahar", "fssai aahar", "nutraceutical", "आयुर्वेद आहार", "आहार"},
			Description: "Food products prepared strictly according to classical Ayurvedic recipe books for dietary use.",
		},
	}

	for _, c := range list {
		bg.concepts[strings.ToLower(c.EnglishTerm)] = c
		for _, s := range c.Synonyms {
			bg.concepts[strings.ToLower(s)] = c
		}
	}
}

func (bg *BilingualGlossary) initFormulations() {
	forms := []AyurvedicFormulationCategory{
		{
			FormName:          "Churna",
			HindiName:         "चूर्ण (Powder)",
			StatutoryCategory: "Classical Formulation (First Schedule)",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940",
			PatentStatus:      "Barred under Section 3(p) if classical recipe; non-patentable aggregation under 3(e)",
			Synonyms:          []string{"churna", "churnam", "powder", "चूर्ण"},
		},
		{
			FormName:          "Kwath / Kashayam",
			HindiName:         "क्वाथ / काढ़ा (Decoction)",
			StatutoryCategory: "Classical Formulation (First Schedule)",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940",
			PatentStatus:      "Barred under Section 3(p); novel extraction technique patentable if demonstrating technical advancement",
			Synonyms:          []string{"kwath", "kashayam", "decoction", "kadha", "क्वाथ", "काढ़ा"},
		},
		{
			FormName:          "Bhasma / Rasashastra",
			HindiName:         "भस्म / रसौषधि (Calcined Mineral/Herbo-mineral)",
			StatutoryCategory: "Classical Herbo-Mineral (First Schedule)",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940 (Heavy Metal Testing mandatory)",
			PatentStatus:      "Classical methods non-patentable (3(p)); novel nanoparticle synthesis process potentially patentable",
			Synonyms:          []string{"bhasma", "rasa", "rasashastra", "swarna bhasma", "भस्म", "रसौषधि"},
		},
		{
			FormName:          "Asava / Arishta",
			HindiName:         "आसव / अरिष्ट (Self-fermented Liquid)",
			StatutoryCategory: "Classical Formulation (First Schedule)",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940 (Rule 158B)",
			PatentStatus:      "Barred under Section 3(p); self-generated bio-fermentation processes require inventive step evidence",
			Synonyms:          []string{"asava", "arishta", "asav", "arisht", "आसव", "अरिष्ट"},
		},
		{
			FormName:          "Taila / Ghrita",
			HindiName:         "तैल / घृत (Medicated Oil/Ghee)",
			StatutoryCategory: "Classical Sneha Kalpana (First Schedule)",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940",
			PatentStatus:      "Classical oil/ghee preparation non-patentable under 3(p); lipid nanoparticle delivery system patentable",
			Synonyms:          []string{"taila", "tailam", "oil", "ghrita", "ghritam", "ghee", "तैल", "तेल", "घृत"},
		},
		{
			FormName:          "Rasayana / Lehyam / Avaleha",
			HindiName:         "रसायन / लेह्य / अवलेह (Confection/Jam)",
			StatutoryCategory: "Classical Rejuvenator (First Schedule)",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940",
			PatentStatus:      "Chyawanprash/Rasayana classical formulations barred under 3(p)",
			Synonyms:          []string{"rasayana", "avaleha", "lehyam", "chyawanprash", "अवलेह", "रसायन", "लेह्य"},
		},
		{
			FormName:          "Vati / Gutika",
			HindiName:         "वटी / गुटिका (Tablet/Pill)",
			StatutoryCategory: "Classical or P&P Tablet Form",
			ApplicableAct:     "Drugs & Cosmetics Act, 1940",
			PatentStatus:      "Compression into tablet is not patentable without novel delayed-release mechanism",
			Synonyms:          []string{"vati", "gutika", "tablet", "pill", "वटी", "गुटिका"},
		},
		{
			FormName:          "Standardized Extract / Fraction",
			HindiName:         "मानकीकृत सत्त / अर्क (Standardized Extract/Phytopharmaceutical)",
			StatutoryCategory: "Patent & Proprietary / Phytopharmaceutical (Rule 122E)",
			ApplicableAct:     "Drugs & Cosmetics Rules (Rule 158B & 122E)",
			PatentStatus:      "Patentable if specific enriched active fractions are isolated and demonstrate non-obvious synergistic efficacy (Sec 3(e))",
			Synonyms:          []string{"extract", "standardized extract", "fraction", "phytopharmaceutical", "सत्त", "अर्क"},
		},
	}

	for _, f := range forms {
		bg.formulations[strings.ToLower(f.FormName)] = f
		for _, s := range f.Synonyms {
			bg.formulations[strings.ToLower(s)] = f
		}
	}
}

// FindConcepts detects all statutory concepts present in the text
func (bg *BilingualGlossary) FindConcepts(text string) []LegalTermConcept {
	lower := strings.ToLower(text)
	var matched []LegalTermConcept
	seen := make(map[string]bool)

	for key, concept := range bg.concepts {
		if strings.Contains(lower, key) {
			if !seen[concept.EnglishTerm] {
				seen[concept.EnglishTerm] = true
				matched = append(matched, concept)
			}
		}
	}

	return matched
}

// FindFormulations identifies Ayurvedic traditional dosage forms and their legal status
func (bg *BilingualGlossary) FindFormulations(text string) []AyurvedicFormulationCategory {
	lower := strings.ToLower(text)
	var matched []AyurvedicFormulationCategory
	seen := make(map[string]bool)

	for key, form := range bg.formulations {
		if strings.Contains(lower, key) {
			if !seen[form.FormName] {
				seen[form.FormName] = true
				matched = append(matched, form)
			}
		}
	}

	return matched
}

// GetDualTermPromptInstructions produces markdown instructions for the LLM to format dual terms in Hindi or English
func (bg *BilingualGlossary) GetDualTermPromptInstructions(language string) string {
	if strings.ToLower(language) == "hi" || strings.ToLower(language) == "hindi" {
		return `### द्विभाषी विधिक शब्दावली नियम (Bilingual Terminology Rendering):
जब आप हिंदी में उत्तर दें, तो कानूनी सटीकता और प्रामाणिकता के लिए निम्नलिखित प्रमुख शब्दों का दोहरा-प्रारूप (Dual-Term) प्रयोग करें:
- एकाधिकार/पेटेंट (Patent)
- पारंपरिक ज्ञान (Traditional Knowledge) [धारा 3(p)]
- सहक्रियात्मक प्रभाव (Synergistic Effect) [धारा 3(e)]
- पूर्व ज्ञान (Prior Art / TKDL)
- जैविक संसाधन (Biological Resource)
- राष्ट्रीय जैव विविधता प्राधिकरण (National Biodiversity Authority - NBA)
- राज्य जैव विविधता बोर्ड (State Biodiversity Board - SBB)
- शास्त्रीय औषधि (Classical Medicine - First Schedule)
- स्वामित्व/पेटेंट एवं प्रोप्राइटरी औषधि (Patent & Proprietary Medicine - Rule 158B)
- आपत्तिजनक विज्ञापन (Objectionable Advertisement - DMRA 1954)
स्पष्टीकरण आम बोलचाल की हिंदी में दें, लेकिन विधिक धाराओं को मूल रूप में उद्धृत करें।`
	}

	return `### Dual-Term Statutory Rendering:
When discussing Indian traditional knowledge, reference recognized statutory provisions alongside their traditional names:
- Traditional Knowledge / Classical texts (e.g. Charaka, Sushruta listed in First Schedule of Drugs & Cosmetics Act)
- Section 3(p) Patent Exclusion for Traditional Knowledge
- Section 3(e) Mere Admixture vs Synergistic Efficacy
- Biological Diversity Act 2023 Form III approval before patent sealing
- DMRA 1954 prohibition for scheduled disease cure assertions`
}
