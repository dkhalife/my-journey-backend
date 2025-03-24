package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Printf("IP:%s UA:%q Route:%s Status:%d\n",
			c.ClientIP(), c.Request.UserAgent(), c.Request.URL.Path, c.Writer.Status())
	}
}
