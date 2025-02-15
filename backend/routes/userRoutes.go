package routes

import (
	"Pinspire/backend/controllers"
	"Pinspire/backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.POST("/register", controllers.RegisterUser)
		user.POST("/login", controllers.LoginUser)
		user.GET("/logout", middleware.AuthMiddleware() /* controllers.LogOutUser */, nil)
		user.GET("/me", middleware.AuthMiddleware(), controllers.MyProfile)
		user.GET("/:id", middleware.AuthMiddleware(), controllers.UserProfile)
		user.POST("/follow/:id", middleware.AuthMiddleware() /* controllers.FollowAndUnfollowUser */, nil)
	}
}
