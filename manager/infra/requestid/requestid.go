package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const headerXRequestID = "X-Request-ID"

func newUUID() *string {
	result := xid.New().String()
	return &result
}

// New is function
func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = *newUUID()
		}
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

// Get is function
func Get(c *gin.Context) string {
	return c.Writer.Header().Get(headerXRequestID)
}
