package post

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"friend_zone/internal/middleware"
	"friend_zone/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/posts", h.Create)
	rg.DELETE("/posts/:content_id", h.Delete)
	rg.GET("/posts/:content_id", h.Get)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	result, err := h.service.Create(c.Request.Context(), middleware.CurrentUserID(c), req)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.Created(c, result)
}

func (h *Handler) Delete(c *gin.Context) {
	contentID, ok := parseIDParam(c, "content_id")
	if !ok {
		return
	}
	err := h.service.Delete(c.Request.Context(), middleware.CurrentUserID(c), contentID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, gin.H{"content_id": contentID, "deleted": true})
}

func (h *Handler) Get(c *gin.Context) {
	contentID, ok := parseIDParam(c, "content_id")
	if !ok {
		return
	}
	result, err := h.service.Get(c.Request.Context(), contentID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	if result.Status == StatusDeleted {
		response.Error(c, http.StatusNotFound, "not_found", "post not found")
		return
	}
	response.OK(c, result)
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrEmptyContent):
		response.Error(c, http.StatusBadRequest, "bad_request", err.Error())
	case errors.Is(err, ErrPostNotFound):
		response.Error(c, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrForbidden):
		response.Error(c, http.StatusForbidden, "forbidden", err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
	}
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	value, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || value <= 0 {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid "+name)
		return 0, false
	}
	return value, true
}
