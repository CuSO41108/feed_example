package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"friend_zone/internal/config"
	"friend_zone/internal/infra/kafka"
	"friend_zone/internal/infra/mysql"
	redisinfra "friend_zone/internal/infra/redis"
	"friend_zone/internal/infra/snowflake"
	"friend_zone/internal/module/auth"
	"friend_zone/internal/module/fanout"
	"friend_zone/internal/module/feed"
	"friend_zone/internal/module/follow"
	"friend_zone/internal/module/post"
	activity "friend_zone/internal/module/user"
	"friend_zone/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := mysql.Open(cfg.MySQL)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer db.Close()

	rdb := redisinfra.New(cfg.Redis)
	defer rdb.Close()

	idgen, err := snowflake.New(1)
	if err != nil {
		log.Fatalf("new snowflake generator: %v", err)
	}
	producer := kafka.NewProducer(cfg.Kafka)
	defer producer.Close()

	activitySvc := activity.NewActivityService(db, cfg.Feed.ActiveWindow)
	authSvc := auth.NewService(db, idgen, activitySvc, cfg.JWT.Secret, cfg.JWT.TTL)
	followSvc := follow.NewService(db)
	postSvc := post.NewService(db, idgen, cfg.Feed)
	feedSvc := feed.NewService(db, rdb, activitySvc, cfg.Feed)

	router := server.NewRouter(cfg.JWT.Secret, server.Handlers{
		Auth:   auth.NewHandler(authSvc),
		Follow: follow.NewHandler(followSvc),
		Post:   post.NewHandler(postSvc),
		Feed:   feed.NewHandler(feedSvc),
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go fanout.NewOutboxRelay(db, producer, cfg.Feed).Run(ctx)

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("api listening on %s", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("api server failed: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("api shutdown error: %v", err)
	}
}
