package routes

import (
	"Pinspire/backend/controllers"
	"Pinspire/backend/middleware"
	"github.com/gin-gonic/gin"
)

func PinRoutes(rg *gin.RouterGroup) {
	pin := rg.Group("/pin", middleware.AuthMiddleware())
	{
		pin.POST("/new", controllers.CreatePin)
		pin.GET("/all", controllers.GetAllPins)               // Implement GetAllPins in controller
		pin.GET("/:id", controllers.GetSinglePin)             // Implement GetSinglePin in controller
		pin.PUT("/:id", controllers.UpdatePin)                // Implement UpdatePin in controller
		pin.DELETE("/:id", controllers.DeletePin)             // Implement DeletePin in controller
		pin.POST("/comment/:id", controllers.CommentOnPin)    // Implement CommentOnPin in controller
		pin.DELETE("/comment/:id", controllers.DeleteComment) // Implement DeleteComment in controller
	}
}
