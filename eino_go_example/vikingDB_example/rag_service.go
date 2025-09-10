package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudwego/eino-ext/components/retriever/volc_vikingdb"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type RAGService struct {
	retriever retriever.Retriever
	chain     compose.Runnable[string, []*schema.Document]
}

type RAGConfig struct {
	VikingDBHost   string
	VikingDBRegion string
	VikingDBAK     string
	VikingDBSK     string
	CollectionName string
	IndexName      string
	ModelName      string
	TopK           int
	ScoreThreshold float64
}

// HTTP request/response types
type QueryRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  *int   `json:"top_k,omitempty"`
}

type QueryResponse struct {
	Documents []*DocumentResponse `json:"documents"`
	Count     int                 `json:"count"`
}

type DocumentResponse struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score,omitempty"`
}

type UploadRequest struct {
	Content  string                 `json:"content" binding:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type UploadResponse struct {
	Message    string `json:"message"`
	DocumentID string `json:"document_id,omitempty"`
}

func NewRAGService(ctx context.Context, config *RAGConfig) (*RAGService, error) {
	// Configure VikingDB retriever
	vikingConfig := &volc_vikingdb.RetrieverConfig{
		Host:              config.VikingDBHost,
		Region:            config.VikingDBRegion,
		AK:                config.VikingDBAK,
		SK:                config.VikingDBSK,
		Scheme:            "https",
		ConnectionTimeout: 0,
		Collection:        config.CollectionName,
		Index:             config.IndexName,
		EmbeddingConfig: volc_vikingdb.EmbeddingConfig{
			UseBuiltin:  true,
			ModelName:   config.ModelName, // e.g., "bge-m3"
			UseSparse:   true,
			DenseWeight: 0.4,
		},
		Partition:      "",
		TopK:           &config.TopK,
		ScoreThreshold: &config.ScoreThreshold,
		FilterDSL:      nil,
	}

	// Create VikingDB retriever
	vikingRetriever, err := volc_vikingdb.NewRetriever(ctx, vikingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create VikingDB retriever: %w", err)
	}

	// Create retrieval chain
	// Create and compile the chain
	chain := compose.NewChain[string, []*schema.Document]()
	chain.AppendRetriever(vikingRetriever)

	compiledChain, err := chain.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to compile chain: %w", err)
	}

	return &RAGService{
		retriever: vikingRetriever,
		chain:     compiledChain,
	}, nil
}

// Core retrieval methods
func (r *RAGService) QueryDocuments(ctx context.Context, query string) ([]*schema.Document, error) {
	// Use the retriever directly
	docs, err := r.retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	log.Printf("VikingDB retrieve success, query=%v, found %d docs", query, len(docs))
	return docs, nil
}

func (r *RAGService) QueryWithChain(ctx context.Context, query string) ([]*schema.Document, error) {
	// Use the compiled chain
	docs, err := r.chain.Invoke(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("chain invocation failed: %w", err)
	}

	log.Printf("Chain retrieve success, query=%v, found %d docs", query, len(docs))
	return docs, nil
}

func (r *RAGService) AddDocument(ctx context.Context, doc *schema.Document) error {
	// Note: Document ingestion typically requires separate VikingDB data management APIs
	// This would involve using VikingDB's data insertion APIs directly
	log.Printf("Document ingestion not implemented in this example: %+v", doc)
	return fmt.Errorf("document ingestion requires VikingDB data management APIs")
}

// HTTP Handlers
func (r *RAGService) Query(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query documents
	docs, err := r.QueryDocuments(c.Request.Context(), req.Query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query documents"})
		return
	}

	// Convert to response format
	docResponses := make([]*DocumentResponse, len(docs))
	for i, doc := range docs {
		docResponses[i] = &DocumentResponse{
			ID:       doc.ID,
			Content:  doc.Content,
			Metadata: doc.MetaData,
			Score:    doc.Score(),
		}
	}

	response := QueryResponse{
		Documents: docResponses,
		Count:     len(docs),
	}

	c.JSON(http.StatusOK, response)
}

func (r *RAGService) UploadDocument(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create document
	doc := &schema.Document{
		ID:       fmt.Sprintf("doc_%d", len(req.Content)), // Simple ID generation
		Content:  req.Content,
		MetaData: req.Metadata,
	}

	// Try to add document (will return error as not implemented)
	err := r.AddDocument(c.Request.Context(), doc)
	if err != nil {
		log.Printf("Document upload failed: %v", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"error":   "Document upload not implemented",
			"message": "This requires VikingDB data management APIs",
		})
		return
	}

	response := UploadResponse{
		Message:    "Document uploaded successfully",
		DocumentID: doc.ID,
	}

	c.JSON(http.StatusCreated, response)
}

func (r *RAGService) ListDocuments(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// This would typically query VikingDB for document listing
	// For now, return a not implemented response
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Document listing not implemented",
		"message": "This requires VikingDB data management APIs",
		"limit":   limit,
		"offset":  offset,
	})
}

func (r *RAGService) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	// This would typically delete from VikingDB
	// For now, return a not implemented response
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":       "Document deletion not implemented",
		"message":     "This requires VikingDB data management APIs",
		"document_id": documentID,
	})
}
