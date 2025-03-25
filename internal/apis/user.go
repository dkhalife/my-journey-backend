package apis

import (
	"dkhalife.com/journey/internal/models"
	"dkhalife.com/journey/internal/utils"
	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt/v2"
)

type UsersAPIHandler struct {
}

func UsersAPI() *UsersAPIHandler {
	return &UsersAPIHandler{}
}

func (h *UsersAPIHandler) signUp(c *gin.Context) {
}

func (h *UsersAPIHandler) CreateAppToken(c *gin.Context) {
}

func (h *UsersAPIHandler) GetAllUserToken(c *gin.Context) {
}

func (h *UsersAPIHandler) DeleteUserToken(c *gin.Context) {
}

func UserRoutes(router *gin.Engine, h *UsersAPIHandler, auth *jwt.GinJWTMiddleware) {
	userRoutes := router.Group("api/v1/users")
	userRoutes.Use(auth.MiddlewareFunc())
	{
		userRoutes.POST("/tokens", utils.ScopeMiddleware(models.ApiTokenScopeTokenWrite), h.CreateAppToken)
		userRoutes.GET("/tokens", utils.ScopeMiddleware(models.ApiTokenScopeTokenWrite), h.GetAllUserToken)
		userRoutes.DELETE("/tokens/:id", utils.ScopeMiddleware(models.ApiTokenScopeTokenWrite), h.DeleteUserToken)
	}

	authRoutes := router.Group("api/v1/auth")
	{
		authRoutes.POST("/", h.signUp)
		authRoutes.POST("login", auth.LoginHandler)
		authRoutes.GET("refresh", auth.RefreshHandler)
	}
}
