package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware(header string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(header)
		if reqID == "" {
			reqID = uuid.New().String()
			c.Request.Header.Set(header, reqID)
		}

		c.Set(header, reqID)
		c.Writer.Header().Set(header, reqID)

		c.Next()
	}
}
