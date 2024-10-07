package handlers

import (
	"github.com/jonesrussell/loggo"
)

type Handler struct {
	logger loggo.LoggerInterface
}

func NewHandler(logger loggo.LoggerInterface) *Handler {
	return &Handler{logger: logger}
}
