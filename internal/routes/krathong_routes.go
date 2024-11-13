package routes

import (
	"loykrathong-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func KrathongRoutes(r *gin.RouterGroup, krathongHandler *handlers.KrathongHandler) {
	r.POST("/krathong", krathongHandler.CreateKrathong)
	r.GET("/krathong", krathongHandler.GetKrathong)
}
