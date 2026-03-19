package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CtxRequestIDKey = "request_id"
const HeaderRequestID = "X-Request-Id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(HeaderRequestID))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(CtxRequestIDKey, requestID)
		c.Header(HeaderRequestID, requestID)
		c.Next()
	}
}
