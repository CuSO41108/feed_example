package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"friend_zone/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/register", h.Register)
	rg.POST("/auth/login", h.Login)
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	result, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		code := "internal_error"
		if errors.Is(err, ErrInvalidInput) {
			status = http.StatusBadRequest
			code = "bad_request"
		}
		response.Error(c, status, code, err.Error())
		return
	}
	response.Created(c, result)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid username or password")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	response.OK(c, result)
}
