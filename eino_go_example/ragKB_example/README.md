# RAG Backend with Eino Framework

A GoLang backend application for Retrieval-Augmented Generation (RAG) using the Eino framework.

## Features

- Document upload and management
- Semantic search and retrieval
- LLM-powered question answering
- RESTful API endpoints
- Docker support

## Quick Start

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up environment variables:**
   Copy `.env.example` to `.env` and fill in your API keys:
   ```bash
   cp .env .env.local
   # Edit .env.local with your API keys
   ```

3. **Run the application:**
   ```bash
   go run .
   ```

4. **Or use Docker:**
   ```bash
   docker-compose up --build
   ```

## API Endpoints

### Health Check
- `GET /health` - Check service health

### Documents
- `POST /api/v1/documents` - Upload a document
- `GET /api/v1/documents` - List all documents
- `DELETE /api/v1/documents/:id` - Delete a document

### Query
- `POST /api/v1/query` - Query the RAG system

## Example Usage

### Upload a Document
```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Eino is a powerful LLM application development framework for Go.",
    "metadata": {"source": "documentation", "type": "guide"}
  }'
```

### Query the System
```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What is Eino?",
    "top_k": 3
  }'
```

## Configuration

The application uses environment variables for configuration. See `.env` for available options.

## Architecture

This RAG system uses the Eino framework components:
- **Embedding**: Convert text to vector representations
- **Vector Store**: Store and search document embeddings
- **Retriever**: Find relevant documents for queries
- **Chat Model**: Generate answers based on retrieved context

## Development

- Built with Go 1.21+
- Uses Gin for HTTP routing
- Eino framework for LLM operations
- Docker support for easy deployment