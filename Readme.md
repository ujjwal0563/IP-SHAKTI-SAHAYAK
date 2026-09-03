# IP-SAKTI Sahayak

### **Intelligent IP & Regulatory Guidance for Ayurveda Innovation**

**A multilingual, RAG-based (source-cited) AI assistant for Intellectual Property, Traditional Knowledge, Access & Benefit Sharing (ABS), and drug-regulatory guidance in Ayurveda across national and international regimes.**

---

# 1. Executive Summary

**IP-SAKTI Sahayak** is an authoritative, multilingual AI-powered decision-support platform engineered specifically for **Ayurveda practitioners, researchers, AYUSH startups, MSMEs, cultivators, and academic institutions**.

Ayurveda rests on an immense corpus of codified and community-held **Traditional Knowledge (TK)** and therapeutics derived from plant, microbial, and animal biological sources. Bringing an Ayurvedic product to life—from research to commercialization—requires navigating several overlapping legal and statutory regimes simultaneously:

- **Intellectual Property (IPR):** Patents, Geographical Indications (GI), Trademarks, Copyright, Industrial Designs, Trade Secrets, and Protection of Plant Varieties & Farmers' Rights (PPV&FR).
- **Biodiversity Sovereignty & ABS:** Access and Benefit Sharing (ABS) duties flowing from India’s sovereign rights over its biological resources under the National Biodiversity Authority (NBA) and State Biodiversity Boards (SBB).
- **Drug & Food Regulatory Classification:** The statutory drug-regulatory framework determining whether a formulation is a **Classical / Generic Medicine**, a **Patent-or-Proprietary (P&P) Medicine**, a **New Drug**, a **Phytopharmaceutical**, an **Ayurveda-Aahar / Nutraceutical**, or a **Cosmetic**.
- **Dual-Layer Jurisdictions:** National legal frameworks vs. International treaties and export-market regulations, strictly bifurcated via an explicit jurisdiction switch so answers are never conflated.

Built on **Retrieval-Augmented Generation (RAG)** grounded in an authoritative, version-tracked corpus of statutes, rules, pharmacopoeias, treaties, and case records, IP-SAKTI Sahayak provides source-cited, traceable guidance with strict guardrails against legal hallucination.

> **Vision:** Empower the AYUSH ecosystem with authoritative, plain-language guidance to protect legitimate innovation, ensure regulatory and ABS compliance, and prevent the misappropriation of India's traditional knowledge worldwide.

---

# 2. The Problem

Protecting and commercializing Ayurvedic innovation is uniquely challenging because **Intellectual Property protection is inseparable from how the product is regulated and sourced**.

```text
========================================================================================
                               AYURVEDIC INNOVATION JOURNEY
========================================================================================
                                         |
                                         v
                         [ REGULATORY CLASSIFICATION ]
       (Classical Text? Proprietary? Phytopharm? Ayurveda-Aahar? Cosmetic?)
                     /                                       \
                    v                                         v
    [ IP RIGHTS & TK DEFENSE ]                     [ BIODIVERSITY & ABS COMPLIANCE ]
    • Patent bar: Section 3(p) TK                  • Biological Diversity Act (2023 / 2024)
    • TKDL defense vs. New Drug patent             • NBA Form I / SBB Approvals
    • GI, Trademarks, Designs, PPV&FR              • Benefit sharing on bio-resources
                    \                                         /
                     +-------------------+-------------------+
                                         |
                                         v
                            [ JURISDICTION GATEWAY ]
               +-------------------------+-------------------------+
               |                                                   |
               v                                                   v
     [ NATIONAL REGIME (INDIA) ]                      [ INTERNATIONAL REGIME ]
     • Patents Act (2024 Rules)                       • WIPO GRATK Treaty (2024)
     • Drugs & Cosmetics Act (Ch. IV-A)               • Nagoya Protocol / CBD
     • DMRA (Advertising Rules)                       • PCT / Madrid / Budapest
     • FSSAI Ayurveda-Aahar (2022)                    • US FDA Botanicals / EU THMPD
========================================================================================
```

### The Core Crisis

Today's AYUSH innovators face a double-edged crisis:
1. **Under-Protected & Under-Commercialized Innovation:** Legitimate innovations developed by domestic Ayurvedic startups, researchers, and MSMEs remain under-protected due to an inability to navigate complex patent barriers (e.g., Section 3(p) non-patentability of traditional knowledge, Section 3(e) mere admixture bar) and drug-approval pathways.
2. **Vulnerability to Biopiracy Abroad:** India’s vast codified and community-held traditional knowledge remains exposed to predatory patenting and misappropriation in foreign patent offices without fair benefit-sharing.

### Why the Problem is More Urgent Than Ever (2023–2024 Landscape Shifts)

Recent fast-moving statutory and treaty shifts have radically altered the compliance landscape:
- **Patents Rules (2024):** Revised patent prosecution timelines, altered fee structures, and updated disclosure frameworks in India.
- **Biological Diversity (Amendment) Act, 2023 & Biological Diversity Rules, 2024:** Overhauled ABS provisions, granting specific exemptions for registered AYUSH practitioners and codified traditional knowledge while establishing stringent, audited benefit-sharing compliance for commercial manufacturing and IPR approvals via the National Biodiversity Authority (NBA).
- **WIPO GRATK Treaty (Adopted May 2024):** A landmark international treaty on Intellectual Property, Genetic Resources, and Associated Traditional Knowledge mandating patent applicants worldwide to disclose the country of origin of genetic resources and traditional knowledge, closing loopholes for foreign misappropriation.
- **Heightened Advertising & Drug Scrutiny:** Aggressive judicial and regulatory enforcement of the *Drugs and Magic Remedies (Objectionable Advertisements) Act (DMRA)* and the rollout of the *FSSAI Ayurveda-Aahar Regulations (2022)*.

### Critical Pain Points for Stakeholders

**1. Inseparability of Formulation Classification & IP Posture**  
An Ayurvedic innovator cannot evaluate patents in a vacuum. A classical formulation (drawn from a First-Schedule text like *Charaka Samhita* or *Sahasrayogam*) faces the **Section 3(p) patenting bar** and is defended via the **Traditional Knowledge Digital Library (TKDL)**. Conversely, an extracted novel formulation or phytopharmaceutical has patent potential, but triggers rigorous clinical safety and drug efficacy evidence requirements.

**2. Jurisdictional Conflation & Export Barriers**  
National requirements (Chapter IV-A of Drugs & Cosmetics Act, Schedule T GMP, State Licensing) are fundamentally different from international export requirements (US FDA Dietary Supplement cGMP / Botanical Guidance, EU Traditional Herbal Medicinal Products Directive 2004/24/EC). Innovators routinely conflate these regimes, resulting in rejected export shipments or invalid foreign patent applications.

**3. The ABS Compliance Maze**  
Many entrepreneurs are unaware that commercial utilization of biological resources or filing patents on bio-resources requires prior approval from the NBA (Form I / Form III) or intimation to State Biodiversity Boards (SBB), carrying severe statutory penalties if overlooked.

**4. Information Fragmentation & Hallucination Risks**  
Authoritative information is scattered across the Indian Patent Office (InPASS), TKDL, NBA, Ministry of AYUSH, Pharmacopoeial Commission (PCIM&H), CDSCO, and WIPO. Generic AI chatbots hallucinate non-existent patent sections, misquote statutory deadlines, and fail to provide verifiable citations.

**5. Language & Accessibility Divide**  
The vast majority of traditional Vaidyas, grassroots innovators, and cultivators operate primarily in Indian languages (e.g., Hindi), whereas legal statutes, treaty articles, and patent gazettes are drafted in dense English legalese.

# 3. Our Solution

IP-SAKTI Sahayak combines:

```text
MULTILINGUAL INTERFACE
|
v
QUERY UNDERSTANDING
|
v
JURISDICTION + IP CLASSIFIER
|
v
AUTHORITY-AWARE RAG
|
v
POSTGRESQL + PGVECTOR
|
v
RELEVANT SOURCE CHUNKS
|
v
LLM
|
v
EVIDENCE + CITATIONS
|
v
ACTION-ORIENTED RESPONSE
```

The system is designed to answer:

> **"What should I consider, why does it matter, and which source supports that information?"**

rather than simply generating a generic response.

# 4. Key Features

## 4.1 Source-Grounded AI Chat

Users can ask questions in natural language.

Example:

> "Can I protect my Ayurvedic herbal formulation using IP?"

The system retrieves relevant sources before generating its answer.

## 4.2 Evidence-Based Citations

Important answers are accompanied by supporting source information:

```text
Answer
|
+-- Source document
+-- Authority
+-- Page / section
+-- Relevant text
+-- Relevance score
```

This improves **traceability and verification**.

## 4.3 Multilingual Assistance

Initial MVP:

- English
- Hindi

The architecture is designed to support additional Indian languages later.

Example:

```text
Hindi Question
|
v
Language Detection
|
v
Cross-language Retrieval
|
v
Authoritative Sources
|
v
Hindi Answer + Citations
```

## 4.4 IP Pathway Identification

The system can identify potentially relevant protection routes such as:

- Patent
- Trademark
- Geographical Indication
- Industrial Design
- Trade Secret
- Copyright
- Plant-variety-related protection where applicable

The system provides **potential areas for review**, not a binding legal conclusion.

## 4.5 Traditional Knowledge Awareness

The system explicitly considers whether an innovation may intersect with existing traditional knowledge.

This is important for Ayurveda because a claimed innovation may not be evaluated in isolation from known traditional practices and documentation.

## 4.6 Biological Resource / ABS Awareness

Where biological resources are involved, the workflow can surface the need to consider:

- Access requirements
- Benefit-sharing considerations
- Applicable authorities
- Relevant regulatory material

## 4.7 Jurisdiction-Aware Guidance

The same innovation may face different requirements depending on where it is protected or commercialized.

Initial focus:

```text
India
```

Future expansion:

```text
India
|
+-- USA
+-- European Union
+-- Other target jurisdictions
```

## 4.8 Guided IP Assessment

Instead of requiring the user to know the correct legal terminology, the system can ask structured questions:

```text
What are you protecting?
|
v
What is unique?
|
v
Is traditional knowledge involved?
|
v
Are biological resources involved?
|
v
Where will you commercialize?
|
v
Potential IP + regulatory areas
```

# 5. Target Users

### Ayurveda Entrepreneurs
Need help identifying potential IP and regulatory considerations before commercialization.

### Researchers & Innovators
Need to understand how an innovation may interact with existing IP and traditional knowledge.

### Ayurvedic Product Developers
Need a structured starting point for protection and compliance research.

### Students & Academic Institutions
Need accessible explanations of IP concepts related to Ayurveda.

### IP/Regulatory Professionals
Can use the platform as a research and information-retrieval assistant.

# 6. Why RAG?

A general-purpose LLM has a major limitation for this domain:

> **It may generate a plausible answer without providing sufficiently verifiable evidence.**

IP-SAKTI uses **Retrieval-Augmented Generation**.

### Without RAG

```text
User Question
|
v
LLM
|
v
Generated Answer
```

### With RAG

```text
User Question
|
v
Query Analysis
|
v
Embedding
|
v
Vector Search
|
v
Relevant Official Documents
|
v
Context Construction
|
v
LLM
|
v
Answer + Citations
```

This allows the knowledge base to be updated independently of the underlying LLM.

# 7. AI Harness / Orchestration

The AI is not treated as a single chatbot call.

The project uses a modular AI workflow:

```text
USER
|
v
Language Detection
|
v
Intent Classification
|
v
IP Category Classification
|
v
Jurisdiction Detection
|
v
Query Normalization
|
v
Embedding Generation
|
v
Vector Retrieval
|
v
Metadata Filtering
|
v
Reranking
|
v
Context Builder
|
v
LLM Generation
|
v
Citation Mapping
|
v
Evidence / Uncertainty Check
|
v
Multilingual Response
```

The **Go backend orchestrates this pipeline**.

This design makes individual components replaceable and easier to evaluate.

# 8. System Architecture

```text
+----------------------+
|        USER          |
+----------+-----------+
|
v
+----------------------+
| React / Next.js UI   |
+----------+-----------+
|
v
+----------------------+
|      Go Backend      |
|      REST API        |
+----------+-----------+
|
+----------------------+----------------------+
|                      |                      |
v                      v                      v
Query Analysis          RAG Pipeline          Conversation
|                      |                  Service
|                      v
|             +------------------+
|             | PostgreSQL       |
|             | + pgvector       |
|             +--------+---------+
|                      |
|                      v
|               Relevant Chunks
|                      |
+----------------------+ 
|
v
+---------------+
| Context       |
| Builder       |
+-------+-------+
|
v
+---------------+
| LLM API       |
+-------+-------+
|
v
+----------------------+
| Citation / Evidence  |
| Verification         |
+----------+-----------+
|
v
+----------------------+
| Multilingual Answer  |
+----------------------+
```

# 9. Technology Stack

| Layer | Technology | Purpose |
|---|---|---|
| Frontend | React / Next.js | User interface |
| Backend | **Go** | API + orchestration |
| API | Gin / Fiber | REST API |
| Database | **PostgreSQL** | Structured application data |
| Vector Search | **pgvector** | Semantic retrieval |
| AI | LLM API | Answer generation |
| Embeddings | Embedding model/API | Vector representation |
| Documents | PDF/HTML processing | Knowledge ingestion |
| Cache | Redis *(later)* | Performance |
| Containerization | Docker | Reproducible deployment |
| Version Control | Git / GitHub | Collaboration |

### Why PostgreSQL + pgvector?

The application needs both:

**Relational data**

```text
Users
Documents
Sources
Messages
Citations
Metadata
```

and:

**Vector data**

```text
Document chunks
Embeddings
Similarity search
```

Using PostgreSQL + pgvector keeps the MVP architecture simple and reduces unnecessary infrastructure.

A dedicated vector database can be introduced later if scale requires it.

# 10. Knowledge Base Strategy

The quality of the knowledge base is as important as the LLM.

The initial knowledge base should prioritize authoritative sources covering:

```text
data/
├── patent/
├── trademark/
├── gi/
├── copyright/
├── design/
├── traditional_knowledge/
├── ayurveda/
├── biodiversity/
└── international/
```

Each document should contain metadata:

```text
Title
Authority
Source URL
Jurisdiction
Category
Language
Publication Date
Effective Date
Version
Last Verified Date
```

### Source Priority

```text
1. Primary legislation / official government source
2. Official authority guidance
3. International intergovernmental organization
4. Academic / research source
5. Secondary commentary
```

The system should not treat an unverified blog and an official government source as equally authoritative.

# 11. RAG Data Pipeline

```text
Authoritative Source
|
v
Document Collection
|
v
PDF / HTML / Text
|
v
Text Extraction
|
v
Cleaning
|
v
Chunking
|
v
Metadata Assignment
|
v
Embedding Generation
|
v
PostgreSQL + pgvector
```

At query time:

```text
Question
|
v
Embedding
|
v
Vector Search
|
v
Metadata Filtering
|
v
Top-K Relevant Chunks
|
v
Context Construction
|
v
LLM
|
v
Answer + Citations
```

# 12. Database Design

Initial schema:

```text
users
sources
jurisdictions
ip_categories
documents
document_chunks
conversations
messages
citations
```

### Relationship

```text
Users
|
v
Conversations
|
v
Messages
|
v
Citations
|
v
Document Chunks
|
v
Documents
|
+------> Sources
|
+------> Jurisdictions
```

### Important RAG table

```text
document_chunks
----------------
id
document_id
chunk_text
chunk_index
page_number
section
embedding
metadata
created_at
```

# 13. Go Backend Structure

```text
ip-sakti/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   │
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── repository/
│   ├── services/
│   │
│   ├── rag/
│   │   ├── embedding.go
│   │   ├── retrieval.go
│   │   ├── context.go
│   │   └── pipeline.go
│   │
│   ├── llm/
│   │   └── client.go
│   │
│   ├── documents/
│   │   ├── parser.go
│   │   ├── chunker.go
│   │   └── ingestion.go
│   │
│   └── citation/
│
├── migrations/
├── data/
│   ├── raw/
│   └── processed/
│
├── frontend/
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

# 14. API Design

Initial API:

```text
POST /api/v1/chat
POST /api/v1/assessment

GET  /api/v1/sources
GET  /api/v1/documents
GET  /api/v1/conversations

POST /api/v1/documents/upload
```

### Example

```http
POST /api/v1/chat
```

```json
{
"message": "Can I patent my Ayurvedic formulation?",
"language": "en",
"jurisdiction": "IN"
}
```

Response:

```json
{
"answer": "Your formulation should be assessed against the applicable patentability requirements...",
"language": "en",
"jurisdiction": "IN",
"citations": [
{
"document": "Relevant official document",
"page": 12,
"source": "Official authority"
}
]
}
```

# 15. Example End-to-End Scenario

### User

> "I have developed an Ayurvedic herbal formulation using traditional herbs and want to sell it under my own brand."

### IP-SAKTI workflow

```text
USER INPUT
|
v
Language Detection
|
v
Intent Analysis
|
v
IP Classification
|
+-------------+-------------+
|             |             |
v             v             v
Patent       Trademark       GI/TK
|             |             |
+-------------+-------------+
|
v
Biological Resource?
|
v
ABS Review
|
v
Regulatory Review
|
v
Authoritative Retrieval
|
v
LLM Response
|
v
Evidence + Sources
```

The user receives an understandable explanation and can inspect the supporting sources.

# 16. Safety & Trust

IP-SAKTI is designed as a **decision-support and information system**, not an automated legal authority.

The system should:

- Prefer retrieved evidence.
- Avoid unsupported legal claims.
- Never invent statutes, sections, deadlines, or procedures.
- Identify jurisdiction.
- Surface uncertainty.
- Prefer current authoritative sources.
- Preserve source metadata.
- Provide citations for important claims.
- Avoid exposing confidential user information unnecessarily.

### Disclaimer

> **IP-SAKTI Sahayak provides informational and decision-support guidance based on retrieved sources. It does not provide a binding legal opinion and does not replace advice from a qualified IP, legal, or regulatory professional.**

# 17. Handling Freshness

Legal and regulatory information can change.

Therefore, documents should track:

```text
publication_date
effective_date
version
last_verified_at
source_url
```

The ingestion system should allow new versions to be added and outdated material to be identified.

This is important because **an accurate answer based on an outdated rule can still be misleading**.

# 18. Evaluation Strategy

The project should be evaluated as an AI information system, not only as a chatbot.

### Retrieval metrics

- Recall@K
- Precision@K
- Relevant-source retrieval rate

### Generation metrics

- Citation correctness
- Citation completeness
- Faithfulness to retrieved evidence
- Hallucination rate

### System metrics

- Response latency
- Retrieval latency
- Error rate

### User metrics

- Task completion
- Answer usefulness
- Language quality
- User satisfaction

A benchmark dataset should contain realistic Ayurveda/IP questions with expected topics and authoritative sources.

# 19. Development Roadmap

### Phase 1 — Foundation

- Go project
- PostgreSQL
- pgvector
- Database migrations
- Basic REST API

### Phase 2 — Knowledge Base

- Source identification
- Document collection
- Text extraction
- Chunking
- Metadata

### Phase 3 — RAG

- Embeddings
- Vector storage
- Similarity retrieval
- Metadata filtering
- Context construction

### Phase 4 — AI

- LLM integration
- Prompt engineering
- Citation mapping
- Uncertainty handling

### Phase 5 — Multilingual

- English
- Hindi
- Additional Indian languages

### Phase 6 — IP Assessment

- IP classification
- Guided questionnaire
- Potential protection pathways
- Traditional knowledge / ABS considerations

### Phase 7 — Product

- Authentication
- Conversation history
- Admin document management
- Evaluation dashboard
- Deployment

# 20. MVP Definition

The MVP is complete when this workflow works reliably:

```text
User
|
v
React UI
|
v
Go API
|
v
Query Analysis
|
v
PostgreSQL + pgvector
|
v
Relevant Official Sources
|
v
LLM
|
v
Citations
|
v
English / Hindi Response
```

### MVP capabilities

- India-focused
- English + Hindi
- RAG-based chat
- Source citations
- Patent
- Trademark
- GI
- Traditional Knowledge
- Biological resources / ABS
- Ayurveda-related regulatory information
- Basic IP assessment

# 21. Future Vision

The long-term platform can evolve from an AI Q&A system into an **Ayurveda IP intelligence and decision-support platform**.

Potential future capabilities:

```text
IP-SAKTI
|
+-------------+-------------+
|             |             |
v             v             v
IP Assistant   TK Analysis   Regulatory
|             |             |
v             v             v
Prior-Art       Traditional    Compliance
Research       Knowledge       Workflow
|             |             |
+-------------+-------------+
|
v
Global Expansion
|
+----------+----------+
|          |          |
USA        EU       Other
```

Future modules may include:

- Prior-art search
- Patent document analysis
- Trademark similarity search
- Regulatory checklist generation
- Document upload and analysis
- Export/commercialization guidance
- More Indian languages
- More international jurisdictions
- Source-update monitoring
- Expert referral

# 22. Project Principles

### **Evidence First**
Important answers should be grounded in retrieved sources.

### **Authority Aware**
Official and authoritative sources receive higher priority.

### **Jurisdiction Aware**
IP and regulatory requirements are not universally identical.

### **Multilingual by Design**
Language should not become a barrier to accessing information.

### **Uncertainty Aware**
The system should say when available evidence is insufficient.

### **Human in the Loop**
The platform supports decision-making; it does not replace qualified professionals.

### **Modular AI**
The RAG and AI pipeline should be replaceable and independently evaluable.

# 23. Quick Start

## Prerequisites

Install:

- Go
- PostgreSQL
- pgvector
- Git
- Docker
- Node.js

## Clone

```bash
git clone <repository-url>
cd ip-sakti
```

## Initialize Go

```bash
go mod init ip-sakti
```

## Environment

```bash
cp .env.example .env
```

Example:

```env
APP_ENV=development
PORT=8080

DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/ip_sakti

LLM_API_KEY=your_key
EMBEDDING_API_KEY=your_key
```

**Never commit real API keys to Git.**

# 24. First Development Milestone

The first milestone is intentionally small:

```text
Go Server
|
v
PostgreSQL
|
v
/health
```

Expected response:

```json
{
"status": "ok",
"database": "connected"
}
```

Only after this works should document ingestion, embeddings, RAG, and LLM integration be added.

# 25. Recommended Engineering Approach

Do not build the entire platform at once.

Build incrementally:

```text
Go API
↓
PostgreSQL
↓
Document Management
↓
Document Ingestion
↓
Embeddings
↓
pgvector Retrieval
↓
RAG
↓
LLM
↓
Citations
↓
Hindi
↓
IP Assessment
↓
International Expansion
```

Avoid premature:

- Microservices
- Kubernetes
- Multiple databases
- Multiple vector databases
- Custom LLM training
- Complex agent frameworks

The first objective is a **reliable, demonstrable end-to-end RAG system**.

# 26. Project Status

| Component | Status |
|---|---|
| Problem definition | ✅ Defined |
| Architecture | ✅ Designed |
| Technology stack | ✅ Selected |
| Go backend | 🚧 Development |
| PostgreSQL schema | 🚧 Development |
| Knowledge base | 🚧 Planning |
| RAG pipeline | 🚧 Planning |
| LLM integration | ⏳ Planned |
| Citation engine | ⏳ Planned |
| English support | ⏳ Planned |
| Hindi support | ⏳ Planned |
| IP assessment | ⏳ Planned |
| International support | ⏳ Future |

# 27. Why IP-SAKTI Matters

The objective is not to create **another generic AI chatbot**.

The objective is to create a **trusted interface between Ayurveda innovation and complex IP/regulatory information**.

```text
AYURVEDA INNOVATOR
|
v
"What should I do?"
|
v
IP-SAKTI
|
+-----------------+-----------------+
|                 |                 |
v                 v                 v
IP ROUTES       TK / ABS CHECKS     REGULATORY
|                 |                 |
+-----------------+-----------------+
|
v
AUTHORITATIVE SOURCES
|
v
EXPLAINED + CITED
|
v
INFORMED NEXT STEPS
```

**IP-SAKTI Sahayak aims to make the first layer of IP and regulatory research for Ayurveda faster, more accessible, evidence-grounded, and multilingual — while keeping qualified human professionals in the loop for final decisions.**

## License

Add the project's chosen license before public release.

## Disclaimer

This project is an AI-powered research and decision-support prototype. It does not constitute legal, patent, regulatory, medical, or professional advice.