package infra

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorSupport struct {
	logger           *Logger
	validatorSupport *ValidatorSupport
}

func NewErrorSupport(logger *Logger, validatorSupport *ValidatorSupport) *ErrorSupport {
	instance := ErrorSupport{
		logger:           logger,
		validatorSupport: validatorSupport,
	}
	return &instance
}

func (e *ErrorSupport) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.ByType(gin.ErrorTypeBind).Last()
			if err != nil {
				e.bindError(c, err)
				return
			}
			err = c.Errors.Last()
			if err != nil {
				e.etcError(c, err)
				return
			}
			e.etcError(c, c.Errors.Last())
		}
	}
}

func (e *ErrorSupport) bindError(c *gin.Context, err *gin.Error) {
	var verr validator.ValidationErrors
	if errors.As(err.Err, &verr) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": e.validatorSupport.Response(verr)})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.JSON())
	}
}

func (e *ErrorSupport) etcError(c *gin.Context, err *gin.Error) {
	c.AbortWithStatus(http.StatusInternalServerError)
}
