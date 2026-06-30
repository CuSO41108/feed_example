package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"friend_zone/internal/middleware"
	"friend_zone/internal/module/auth"
	"friend_zone/internal/module/feed"
	"friend_zone/internal/module/follow"
	"friend_zone/internal/module/post"
	"friend_zone/internal/pkg/response"
)

type Handlers struct {
	Auth   *auth.Handler
	Follow *follow.Handler
	Post   *post.Handler
	Feed   *feed.Handler
}

func NewRouter(jwtSecret string, handlers Handlers) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	router.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})
	router.StaticFile("/", "web/index.html")
	router.Static("/assets", "web/assets")
	router.GET("/swagger/index.html", swaggerIndex)
	router.StaticFile("/swagger/doc.yaml", "docs/openapi.yaml")

	api := router.Group("/api/v1")
	handlers.Auth.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.Auth(jwtSecret))
	handlers.Follow.RegisterRoutes(protected)
	handlers.Post.RegisterRoutes(protected)
	handlers.Feed.RegisterRoutes(protected)

	return router
}

func swaggerIndex(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Friend Zone API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({ url: '/swagger/doc.yaml', dom_id: '#swagger-ui' });
    };
  </script>
</body>
</html>`))
}
