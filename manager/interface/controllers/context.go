package controllers

import (
	"net/http"
	"time"
)

// Context is interface
type Context interface {
	Param(string) string
	GetHeader(key string) string
	Bind(interface{}) error
	ShouldBindQuery(obj interface{}) error
	ShouldBind(obj interface{}) error
	ShouldBindHeader(obj interface{}) error

	Status(int)
	String(int, string, ...interface{})
	JSON(int, interface{})
	InternalError(error)
	BindError(error)
	RequestID() string
	Request() *http.Request

	/************************************/
	/***** GOLANG.ORG/X/NET/CONTEXT *****/
	/************************************/
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type HTTPHandler interface {
	Handler(c Context)
}
