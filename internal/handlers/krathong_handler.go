package handlers

import (
	"log"
	"net/http"
	"time"

	"loykrathong-api/internal/models"
	"loykrathong-api/pkg/kafka"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KrathongHandler handles requests for Krathong-related operations
type KrathongHandler struct {
	DB            *gorm.DB
	KafkaProducer *kafka.Producer // Using KafkaProducer instead of KafkaWriter
}

// Struct for the Krathong data response
type KrathongDataResponse struct {
	ResponseCode    string          `json:"response_code"`
	ResponseMessage string          `json:"response_message"`
	Data            models.Krathong `json:"data"`
}

// Struct for a successful list of Krathong items
type GetKrathongListResponse struct {
	ResponseCode    string            `json:"response_code"`
	ResponseMessage string            `json:"response_message"`
	Data            []models.Krathong `json:"data"`
}

// CreateKrathong handles the creation of a new Krathong entry
// @Summary Create Krathong
// @Description Create a new Krathong entry
// @Tags Krathong
// @Accept  json
// @Produce  json
// @Param krathong body models.Krathong true "Krathong Data"
// @Success 201 {object} KrathongDataResponse
// @Failure 400 {object} KrathongDataResponse
// @Failure 500 {object} KrathongDataResponse
// @Router /krathong [post]
func (h *KrathongHandler) CreateKrathong(c *gin.Context) {
	var krathong models.Krathong

	// Bind JSON input to the Krathong model
	if err := c.ShouldBindJSON(&krathong); err != nil {
		c.JSON(http.StatusBadRequest, KrathongDataResponse{
			ResponseCode:    "0001",
			ResponseMessage: "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Set the CreatedAt field
	krathong.CreatedAt = time.Now()

	// Create the Krathong record using GORM's Create method
	if err := h.DB.Create(&krathong).Error; err != nil {
		c.JSON(http.StatusInternalServerError, KrathongDataResponse{
			ResponseCode:    "0004",
			ResponseMessage: "Failed to create Krathong: " + err.Error(),
		})
		return
	}

	// Publish the created Krathong to Kafka
	// message, err := json.Marshal(krathong)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, KrathongDataResponse{
	// 		ResponseCode:    "0004",
	// 		ResponseMessage: "Failed to encode Krathong data: " + err.Error(),
	// 	})
	// 	return
	// }

	// Log before publishing to Kafka
	// log.Println("Attempting to publish to Kafka...")

	// if err := h.KafkaProducer.PublishMessage("krathong_key", message); err != nil {
	// 	log.Printf("Failed to publish to Kafka: %v", err)
	// 	c.JSON(http.StatusInternalServerError, KrathongDataResponse{
	// 		ResponseCode:    "0004",
	// 		ResponseMessage: "Failed to publish message to Kafka: " + err.Error(),
	// 	})
	// 	return
	// }

	// Log success message
	log.Println("Message successfully published to Kafka")

	c.JSON(http.StatusCreated, KrathongDataResponse{
		ResponseCode:    "0000",
		ResponseMessage: "Krathong created successfully",
		Data:            krathong,
	})
}

// GetKrathong retrieves the latest 50 Krathong entries ordered by creation date
// @Summary Get Krathongs
// @Description Retrieve the latest 50 Krathong entries
// @Tags Krathong
// @Produce json
// @Success 200 {object} GetKrathongListResponse
// @Failure 500 {object} KrathongDataResponse
// @Router /krathong [get]
func (h *KrathongHandler) GetKrathong(c *gin.Context) {
	var krathongs []models.Krathong

	// Use raw SQL to retrieve the latest 50 records
	query := "SELECT * FROM krathongs ORDER BY created_at DESC LIMIT 50"
	if err := h.DB.Raw(query).Scan(&krathongs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, GetKrathongListResponse{
			ResponseCode:    "0004",
			ResponseMessage: "Failed to retrieve Krathongs - Database query error",
		})
		return
	}

	c.JSON(http.StatusOK, GetKrathongListResponse{
		ResponseCode:    "0000",
		ResponseMessage: "Krathongs retrieved successfully",
		Data:            krathongs,
	})
}
