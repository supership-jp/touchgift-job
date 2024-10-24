package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
)

type AppContext struct {
	*gin.Context
}

const headerXRequestID = "x-request-id"

func NewContext(c *gin.Context) *AppContext {
	return &AppContext{c}
}

func (c *AppContext) InternalError(err error) {
	c.Error(err).SetType(gin.ErrorTypePrivate) //nolint:errcheck // ミドルウェアでハンドリングするため
}

func (c *AppContext) SetResponseHeader() {
	c.Header(headerXRequestID, c.RequestID())
}

func (c *AppContext) BindError(err error) {
	c.Error(err).SetType(gin.ErrorTypeBind) //nolint:errcheck // ミドルウェアでハンドリングするため
}

func (c *AppContext) RequestID() string {
	requestID := c.Writer.Header().Get(headerXRequestID)
	if requestID == "" {
		// クエリにidが含まれる場合はidをrequestIDとして扱う
		requestID = c.Query("id")
		if requestID == "" {
			// 含まれていない場合は新規生成する
			requestID = xid.New().String()
		}
		c.Header(headerXRequestID, requestID)
	}
	return requestID
}

func (c *AppContext) Request() *http.Request {
	return c.Context.Request
}
