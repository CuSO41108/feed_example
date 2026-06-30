package follow

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
	rg.POST("/follows/:followee_id", h.Follow)
	rg.DELETE("/follows/:followee_id", h.Unfollow)
}

func (h *Handler) Follow(c *gin.Context) {
	followeeID, ok := parseIDParam(c, "followee_id")
	if !ok {
		return
	}
	err := h.service.Follow(c.Request.Context(), middleware.CurrentUserID(c), followeeID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, gin.H{"followee_id": strconv.FormatInt(followeeID, 10), "following": true})
}

func (h *Handler) Unfollow(c *gin.Context) {
	followeeID, ok := parseIDParam(c, "followee_id")
	if !ok {
		return
	}
	err := h.service.Unfollow(c.Request.Context(), middleware.CurrentUserID(c), followeeID)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.OK(c, gin.H{"followee_id": strconv.FormatInt(followeeID, 10), "following": false})
}

func (h *Handler) writeError(c *gin.Context, err error) {
	if errors.Is(err, ErrCannotFollowSelf) {
		response.Error(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	response.Error(c, http.StatusInternalServerError, "internal_error", err.Error())
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	value, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || value <= 0 {
		response.Error(c, http.StatusBadRequest, "bad_request", "invalid "+name)
		return 0, false
	}
	return value, true
}
