package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type RAGService struct {
	chatModel model.ChatModel
	config    *RAGConfig
}

type RAGConfig struct {
	// ragKB Configuration
	KnowledgeBaseDomain string
	AccountID           string
	AccessKey           string
	SecretKey           string
	Region              string
	ProjectName         string
	CollectionName      string
	// ARK Configuration
	ARKAPIKey  string
	ARKBaseURL string
	ChatModel  string
}

// ragKB API request/response types
type SearchKnowledgeRequest struct {
	Project        string         `json:"project"`
	Name           string         `json:"name"`
	Query          string         `json:"query"`
	Limit          int            `json:"limit"`
	PreProcessing  PreProcessing  `json:"pre_processing"`
	DenseWeight    float64        `json:"dense_weight"`
	PostProcessing PostProcessing `json:"post_processing"`
}

type PreProcessing struct {
	NeedInstruction  bool      `json:"need_instruction"`
	Rewrite          bool      `json:"rewrite"`
	ReturnTokenUsage bool      `json:"return_token_usage"`
	Messages         []Message `json:"messages"`
}

type PostProcessing struct {
	GetAttachmentLink   bool `json:"get_attachment_link"`
	ChunkGroup          bool `json:"chunk_group"`
	RerankOnlyChunk     bool `json:"rerank_only_chunk"`
	RerankSwitch        bool `json:"rerank_switch"`
	ChunkDiffusionCount int  `json:"chunk_diffusion_count"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SearchKnowledgeResponse struct {
	Code int                 `json:"code"`
	Data SearchKnowledgeData `json:"data"`
}

type SearchKnowledgeData struct {
	ResultList []KnowledgePoint `json:"result_list"`
}

type KnowledgePoint struct {
	Content          string            `json:"content"`
	OriginalQuestion string            `json:"original_question,omitempty"`
	ChunkTitle       string            `json:"chunk_title,omitempty"`
	DocInfo          DocInfo           `json:"doc_info"`
	ChunkAttachment  []ChunkAttachment `json:"chunk_attachment,omitempty"`
	TableChunkFields []TableChunkField `json:"table_chunk_fields,omitempty"`
	Score            float64           `json:"score,omitempty"`
}

type DocInfo struct {
	DocName string `json:"doc_name"`
	Title   string `json:"title"`
}

type ChunkAttachment struct {
	Link string `json:"link"`
}

type TableChunkField struct {
	FieldName  string `json:"field_name"`
	FieldValue string `json:"field_value"`
}

// HTTP request/response types
type QueryRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  *int   `json:"top_k,omitempty"`
}

type QueryResponse struct {
	Documents []*DocumentResponse `json:"documents"`
	Count     int                 `json:"count"`
	Answer    string              `json:"answer,omitempty"`
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
	// Create ARK chat model
	arkConfig := &ark.ChatModelConfig{
		APIKey:  config.ARKAPIKey,
		BaseURL: config.ARKBaseURL,
		Model:   config.ChatModel,
	}

	chatModel, err := ark.NewChatModel(ctx, arkConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ARK chat model: %w", err)
	}

	return &RAGService{
		chatModel: chatModel,
		config:    config,
	}, nil
}

// Sign request using AWS Signature Version 4
func (r *RAGService) signRequest(req *http.Request, body []byte) error {
	// Set timestamp
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dateStamp := timestamp[:8]

	// Calculate content hash
	hasher := sha256.New()
	hasher.Write(body)
	contentSha256 := hex.EncodeToString(hasher.Sum(nil))

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", req.Host)
	req.Header.Set("X-Content-Sha256", contentSha256)
	req.Header.Set("X-Date", timestamp)
	req.Header.Set("V-Account-Id", r.config.AccountID)

	// Create canonical headers (must be sorted)
	canonicalHeaders := "content-type:application/json\n" +
		"host:" + req.Host + "\n" +
		"x-content-sha256:" + contentSha256 + "\n" +
		"x-date:" + timestamp + "\n"

	signedHeaders := "content-type;host;x-content-sha256;x-date"

	// Create canonical request
	canonicalRequest := req.Method + "\n" +
		req.URL.Path + "\n" +
		req.URL.RawQuery + "\n" +
		canonicalHeaders + "\n" +
		signedHeaders + "\n" +
		contentSha256

	// Hash canonical request
	hasher = sha256.New()
	hasher.Write([]byte(canonicalRequest))
	canonicalRequestHash := hex.EncodeToString(hasher.Sum(nil))

	// Create credential scope
	credentialScope := dateStamp + "/" + r.config.Region + "/air/request"

	// Create string to sign
	stringToSign := "HMAC-SHA256\n" +
		timestamp + "\n" +
		credentialScope + "\n" +
		canonicalRequestHash

	// Calculate signature
	kDate := hmacSHA256([]byte(r.config.SecretKey), dateStamp)
	kRegion := hmacSHA256(kDate, r.config.Region)
	kService := hmacSHA256(kRegion, "air")
	kSigning := hmacSHA256(kService, "request")
	signature := hex.EncodeToString(hmacSHA256(kSigning, stringToSign))

	// Create authorization header
	authorization := "HMAC-SHA256 Credential=" + r.config.AccessKey + "/" + credentialScope +
		", SignedHeaders=" + signedHeaders + ", Signature=" + signature

	req.Header.Set("Authorization", authorization)

	return nil
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Search knowledge using ragKB API
func (r *RAGService) SearchKnowledge(ctx context.Context, query string, limit int) ([]*schema.Document, error) {
	// Prepare request payload
	reqPayload := SearchKnowledgeRequest{
		Project: r.config.ProjectName,
		Name:    r.config.CollectionName,
		Query:   query,
		Limit:   limit,
		PreProcessing: PreProcessing{
			NeedInstruction:  true,
			Rewrite:          false,
			ReturnTokenUsage: true,
			Messages: []Message{
				{Role: "system", Content: ""},
				{Role: "user", Content: query},
			},
		},
		DenseWeight: 0.5,
		PostProcessing: PostProcessing{
			GetAttachmentLink:   true,
			ChunkGroup:          true,
			RerankOnlyChunk:     false,
			RerankSwitch:        false,
			ChunkDiffusionCount: 0,
		},
	}

	log.Printf("ragKB request: project=%s, collection=%s, query=%s", r.config.ProjectName, r.config.CollectionName, query)

	// Marshal request body
	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := "http://" + r.config.KnowledgeBaseDomain + "/api/knowledge/collection/search_knowledge"
	log.Printf("ragKB API URL: %s", url)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Host", r.config.KnowledgeBaseDomain)
	req.Header.Set("V-Account-Id", r.config.AccountID)

	// Sign the request
	if err := r.signRequest(req, body); err != nil {
		log.Printf("Failed to sign request: %v", err)
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("ragKB API response status: %d", resp.StatusCode)

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("ragKB API response body: %s", string(respBody))

	// Check HTTP status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ragKB API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var searchResp SearchKnowledgeResponse
	if err := json.Unmarshal(respBody, &searchResp); err != nil {
		log.Printf("Failed to unmarshal response: %v, body: %s", err, string(respBody))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if searchResp.Code != 0 {
		log.Printf("ragKB API error: code %d, response: %+v", searchResp.Code, searchResp)
		return nil, fmt.Errorf("ragKB API error: code %d", searchResp.Code)
	}

	// Convert to schema.Document format
	docs := make([]*schema.Document, len(searchResp.Data.ResultList))
	for i, point := range searchResp.Data.ResultList {
		// Create metadata from doc info and other fields
		metadata := map[string]interface{}{
			"doc_name":    point.DocInfo.DocName,
			"title":       point.DocInfo.Title,
			"chunk_title": point.ChunkTitle,
		}

		if point.OriginalQuestion != "" {
			metadata["original_question"] = point.OriginalQuestion
		}

		// Store score in metadata since schema.Document doesn't have SetScore method
		if point.Score > 0 {
			metadata["score"] = point.Score
		}

		docs[i] = &schema.Document{
			ID:       fmt.Sprintf("ragkb_%d", i),
			Content:  point.Content,
			MetaData: metadata,
		}
	}

	log.Printf("ragKB search success, query=%v, found %d docs", query, len(docs))
	return docs, nil
}

// Core retrieval methods
func (r *RAGService) QueryDocuments(ctx context.Context, query string) ([]*schema.Document, error) {
	return r.SearchKnowledge(ctx, query, 10)
}

func (r *RAGService) QueryWithChain(ctx context.Context, query string) ([]*schema.Document, error) {
	// For compatibility, use the same search method
	return r.SearchKnowledge(ctx, query, 10)
}

// RAG with chat model
func (r *RAGService) QueryWithRAG(ctx context.Context, query string) (string, []*schema.Document, error) {
	// First, retrieve relevant documents using ragKB
	docs, err := r.SearchKnowledge(ctx, query, 10)
	if err != nil {
		return "", nil, fmt.Errorf("document retrieval failed: %w", err)
	}

	// Build context from retrieved documents
	var contextParts []string
	for _, doc := range docs {
		contextParts = append(contextParts, doc.Content)
	}
	context := strings.Join(contextParts, "\n\n")

	// Create prompt with context and query
	prompt := fmt.Sprintf(`Based on the following context, please answer the question.

Context:
%s

Question: %s

Answer:`, context, query)

	// Generate response using chat model
	messages := []*schema.Message{
		schema.SystemMessage("You are a helpful assistant that answers questions based on the provided context."),
		schema.UserMessage(prompt),
	}

	response, err := r.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", docs, fmt.Errorf("chat model generation failed: %w", err)
	}

	return response.Content, docs, nil
}

func (r *RAGService) AddDocument(ctx context.Context, doc *schema.Document) error {
	// Note: Document ingestion requires ragKB data management APIs
	log.Printf("Document ingestion not implemented in this example: %+v", doc)
	return fmt.Errorf("document ingestion requires ragKB data management APIs")
}

// HTTP Handlers
func (r *RAGService) Query(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if we should use RAG (generate answer) or just retrieve documents
	useRAG := c.Query("rag") == "true"

	if useRAG {
		// Use RAG to generate answer
		answer, docs, err := r.QueryWithRAG(c.Request.Context(), req.Query)
		if err != nil {
			log.Printf("RAG query failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process RAG query"})
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
			Answer:    answer,
		}

		c.JSON(http.StatusOK, response)
	} else {
		// Just retrieve documents
		docs, err := r.QueryDocuments(c.Request.Context(), req.Query)
		if err != nil {
			log.Printf("Query documents failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query documents: %v", err)})
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
			"message": "This requires ragKB data management APIs",
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

	// This would typically query ragKB for document listing
	// For now, return a not implemented response
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Document listing not implemented",
		"message": "This requires ragKB data management APIs",
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

	// This would typically delete from ragKB
	// For now, return a not implemented response
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":       "Document deletion not implemented",
		"message":     "This requires ragKB data management APIs",
		"document_id": documentID,
	})
}
