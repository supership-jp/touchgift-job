package controllers

import (
	"net/http"
	"touchgift-job-manager/manager/usecase"
)

type ping struct {
	logger usecase.Logger
}

func NewPing(logger usecase.Logger) HTTPHandler {
	instance := ping{
		logger: logger,
	}
	return &instance
}

func (p *ping) Handler(c Context) {
	c.String(http.StatusOK, "pong")
}
