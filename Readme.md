# IP-SAKTI Sahayak

### **Intelligent IP & Regulatory Guidance for Ayurveda Innovation**

**A multilingual, source-grounded RAG platform for navigating Intellectual Property, Traditional Knowledge, biological-resource considerations, and regulatory pathways for Ayurveda.**

## 1. Executive Summary

**IP-SAKTI Sahayak** is a multilingual AI-powered decision-support platform designed to help **Ayurveda innovators, researchers, entrepreneurs, practitioners, and organizations** navigate the complex intersection of:

- Intellectual Property (IP)
- Traditional Knowledge (TK)
- Biological resources and Access & Benefit Sharing (ABS)
- Ayurveda-related regulation
- National and international IP/regulatory frameworks

The core challenge is not a lack of information. The challenge is that **relevant information is fragmented across different laws, authorities, portals, regulations, documents, and jurisdictions**.

IP-SAKTI Sahayak addresses this fragmentation through a **Retrieval-Augmented Generation (RAG)** architecture. Instead of relying only on an LLM's internal knowledge, the system retrieves relevant information from a curated knowledge base and generates an answer that is **grounded in evidence and linked to source documents**.

> **Vision:** Make complex IP and regulatory information more accessible, traceable, multilingual, and actionable for Ayurveda innovation.

# 2. The Problem

Ayurveda is supported by a vast body of codified and community-held traditional knowledge and frequently involves biological resources.

When an Ayurvedic product or innovation moves from **idea → research → protection → commercialization**, an innovator may need to answer several questions simultaneously:

```text
AYURVEDA INNOVATION
|
+-----------------+-----------------+
|                 |                 |
v                 v                 v
IP RIGHTS       TRADITIONAL        REGULATORY
KNOWLEDGE          REQUIREMENTS
|                 |                 |
v                 v                 v
Patent / TM /      Existing TK /      Product /
GI / Design /      Prior knowledge    market rules
Trade Secret
|
v
BIOLOGICAL RESOURCES
|
v
ABS / Compliance
|
v
INTERNATIONAL MARKET
```

### Existing pain points

**1. Fragmented information**  
Relevant information is distributed across multiple official sources.

**2. Complex terminology**  
Patent, prior art, traditional knowledge, GI, ABS, regulatory and jurisdiction-specific concepts can be difficult for non-specialists.

**3. Traditional Knowledge considerations**  
An apparently "new" Ayurvedic claim may intersect with knowledge already documented or held traditionally.

**4. Multiple legal/regulatory layers**  
IP protection may be only one part of the commercialization journey.

**5. Language barrier**  
Important information is often available primarily in English, while many users are more comfortable using Indian languages.

**6. Source-traceability problem**  
A generic AI answer may not clearly show which authoritative document supports an important claim.

**7. Changing rules**  
Laws, regulations, procedures, notifications, and official guidance can change over time.

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