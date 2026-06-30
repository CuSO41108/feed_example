package feed

import (
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
	rg.GET("/feed/timeline", h.Timeline)
}

func (h *Handler) Timeline(c *gin.Context) {
	limit := 0
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "bad_request", "invalid limit")
			return
		}
		limit = parsed
	}
	query := Query{
		Direction: c.DefaultQuery("direction", DirectionLatest),
		Cursor:    c.Query("cursor"),
		Limit:     limit,
	}
	if query.Direction != DirectionLatest && query.Direction != DirectionNewer && query.Direction != DirectionOlder {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid direction")
		return
	}

	result, err := h.service.Timeline(c.Request.Context(), middleware.CurrentUserID(c), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	response.OK(c, result)
}
