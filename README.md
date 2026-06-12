# RAG Assistant

A Retrieval-Augmented Generation (RAG) assistant built in Go using Ollama for embeddings and LLM inference.

## Architecture

```
rag-assistant/
├── cmd/server/        # HTTP server entry point
├── internal/
│   ├── ingest/        # Document loading from disk
│   ├── chunk/         # Text chunking with overlap
│   ├── embed/         # Embedding via Ollama API
│   ├── retrieve/      # Cosine similarity search
│   ├── llm/           # LLM chat via Ollama API
│   ├── storage/       # Index persistence (JSON)
│   └── rag/           # Orchestrator
├── data/docs/         # Source markdown documents
├── index/             # Pre-built vector index
└── docker-compose.yml # Ollama + rag-server
```

## Prerequisites

- Go 1.26+
- Ollama running locally with `nomic-embed-text` and `tinyllama` models

## Quick Start

```bash
# Run with local Ollama
go run ./cmd/server

# Or with Docker
docker compose up --build
```

## API

### Query
```bash
curl -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"query": "talk to me about scaling systems"}'
```

### Rebuild Index
```bash
curl -X POST http://localhost:8080/api/ingest
```

### Health
```bash
curl http://localhost:8080/health
```
