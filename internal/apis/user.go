package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersAPIHandler struct {
}

func UsersAPI() *UsersAPIHandler {
	return &UsersAPIHandler{}
}

func (h *UsersAPIHandler) getUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world",
	})
}

func UserRoutes(router *gin.Engine, h *UsersAPIHandler) {
	userRoutes := router.Group("api/v1/users")
	{
		userRoutes.GET("/", h.getUser)
	}
}
