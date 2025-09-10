package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	// Initialize logger
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Create RAG configuration from environment variables
	config := &RAGConfig{
		VikingDBHost:      getEnvOrDefault("VIKINGDB_HOST", "vikingdb.volces.com"),
		VikingDBRegion:    getEnvOrDefault("VIKINGDB_REGION", "cn-beijing"),
		VikingDBAK:        os.Getenv("VIKINGDB_AK"),
		VikingDBSK:        os.Getenv("VIKINGDB_SK"),
		CollectionName:    getEnvOrDefault("VIKINGDB_COLLECTION", "rag_collection"),
		IndexName:         getEnvOrDefault("VIKINGDB_INDEX", "rag_index"),
		ModelName:         getEnvOrDefault("VIKINGDB_MODEL", "bge-m3"),
		TopK:              getEnvAsInt("VIKINGDB_TOP_K", 5),
		ScoreThreshold:    getEnvAsFloat("VIKINGDB_SCORE_THRESHOLD", 0.7),
		// ARK Configuration
		ARKAPIKey:  os.Getenv("ARK_API_KEY"),
		ARKBaseURL: getEnvOrDefault("ARK_BASE_URL", "https://ark.ap-southeast.bytepluses.com/api/v3"),
		ChatModel:  getEnvOrDefault("CHAT_MODEL", "seed-1-6-250615"),
	}

	// Validate required environment variables
	if config.VikingDBAK == "" || config.VikingDBSK == "" {
		log.Fatal("VIKINGDB_AK and VIKINGDB_SK environment variables are required")
	}
	if config.ARKAPIKey == "" {
		log.Fatal("ARK_API_KEY environment variable is required")
	}

	// Initialize RAG service
	ctx := context.Background()
	ragService, err := NewRAGService(ctx, config)
	if err != nil {
		log.Fatal("Failed to initialize RAG service:", err)
	}

	// Initialize router
	router := setupRouter(ragService)

	// Start server
	port := getEnvOrDefault("PORT", "8080")
	logrus.Infof("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(ragService *RAGService) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/documents", ragService.UploadDocument)
		api.POST("/query", ragService.Query)
		api.GET("/documents", ragService.ListDocuments)
		api.DELETE("/documents/:id", ragService.DeleteDocument)
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}