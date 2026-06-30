package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"friend_zone/internal/pkg/response"
)

const CurrentUserIDKey = "current_user_id"

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "missing bearer token")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid token claims")
			c.Abort()
			return
		}
		sub, err := claims.GetSubject()
		if err != nil || sub == "" {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "missing subject")
			c.Abort()
			return
		}
		userID, ok := parseInt64(sub)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "invalid subject")
			c.Abort()
			return
		}

		c.Set(CurrentUserIDKey, userID)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) int64 {
	value, exists := c.Get(CurrentUserIDKey)
	if !exists {
		return 0
	}
	userID, _ := value.(int64)
	return userID
}

func parseInt64(raw string) (int64, bool) {
	var out int64
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		out = out*10 + int64(ch-'0')
	}
	return out, out > 0
}
