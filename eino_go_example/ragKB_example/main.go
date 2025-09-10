package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// ragKB Configuration
	config := &RAGConfig{
		KnowledgeBaseDomain: getEnvOrDefault("RAGKB_DOMAIN", "api-knowledgebase.mlp.cn-hongkong.bytepluses.com"),
		AccountID:           getEnvOrDefault("RAGKB_ACCOUNT_ID", "3000822947"),
		AccessKey:           getEnvOrDefault("RAGKB_ACCESS_KEY", ""),
		SecretKey:           getEnvOrDefault("RAGKB_SECRET_KEY", ""),
		Region:              getEnvOrDefault("RAGKB_REGION", "cn-hongkong"),
		ProjectName:         getEnvOrDefault("RAGKB_PROJECT", "default"),
		CollectionName:      getEnvOrDefault("RAGKB_COLLECTION", "test"),
		// ARK Configuration
		ARKAPIKey:  getEnvOrDefault("ARK_API_KEY", ""),
		ARKBaseURL: getEnvOrDefault("ARK_BASE_URL", "https://ark.cn-beijing.volces.com/api/v3"),
		ChatModel:  getEnvOrDefault("ARK_CHAT_MODEL", "ep-20241211105246-lmqdx"),
	}

	// Validate required configuration
	if config.AccessKey == "" {
		log.Fatal("RAGKB_ACCESS_KEY is required")
	}
	if config.SecretKey == "" {
		log.Fatal("RAGKB_SECRET_KEY is required")
	}
	if config.ARKAPIKey == "" {
		log.Fatal("ARK_API_KEY is required")
	}

	// Initialize RAG service
	ctx := context.Background()
	ragService, err := NewRAGService(ctx, config)
	if err != nil {
		log.Fatal("Failed to initialize RAG service:", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// RAG endpoints
	r.POST("/query", ragService.Query)
	r.POST("/documents", ragService.UploadDocument)
	r.GET("/documents", ragService.ListDocuments)
	r.DELETE("/documents/:id", ragService.DeleteDocument)

	// Start server
	port := getEnvOrDefault("PORT", "8080")
	log.Printf("Starting RAG server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}