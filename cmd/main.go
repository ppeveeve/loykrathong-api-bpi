package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"loykrathong-api/config"
	"loykrathong-api/pkg/kafka"

	"loykrathong-api/internal/handlers"
	"loykrathong-api/internal/middleware"
	"loykrathong-api/internal/routes"
	pkg "loykrathong-api/pkg/database"

	_ "loykrathong-api/docs"
)

// @title Loy Krathong API
// @version 1.0
// @description This is a API for the Loy Krathong festival.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 10.144.130.86:8080
// @BasePath /api/v1
func main() {

	log.Println("Starting application...")

	// Load the application environment variable
	env := getEnv("APP_ENV", "dev") // Use a helper function to get the environment

	// Load configuration based on the environment
	cfg := config.LoadConfig(env)

	// Check required config values are set
	requiredConfig := []struct {
		value string
		name  string
	}{
		{cfg.DBUser, "DBUser"},
		{cfg.DBPassword, "DBPassword"},
		{cfg.DBName, "DBName"},
		{cfg.DBHost, "DBHost"},
		{cfg.DBPort, "DBPort"},
		{cfg.KafkaBroker, "KafkaBroker"},
		{cfg.KafkaTopic, "KafkaTopic"},
	}

	for _, config := range requiredConfig {
		if config.value == "" {
			log.Fatalf("Missing required configuration: %s", config.name)
		}
	}

	// Connect to the database
	db, err := pkg.DatabaseConnect(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Println("Successfully connected to the database")
	}

	// Initialize Kafka producer from the custom kafka package
	kafkaProducer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	log.Println("Kafka producer initialized")

	defer kafkaProducer.Writer.Close()

	// Initialize and run the server
	runServer(cfg.AppPort, db, kafkaProducer)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// runServer initializes the Gin router and starts the server
func runServer(port string, db *gorm.DB, kafkaProducer *kafka.Producer) {
	r := gin.Default()

	// Enable CORS for all origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow any origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Create KrathongHandler as a pointer
	krathongHandler := &handlers.KrathongHandler{
		DB:            db,
		KafkaProducer: kafkaProducer,
	}

	// Add logging middleware
	r.Use(middleware.Logger())

	// Define the API routes
	api := r.Group("/api/v1")
	// Pass krathongHandler as a pointer to the routes
	routes.KrathongRoutes(api, krathongHandler)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
