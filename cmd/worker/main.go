package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"friend_zone/internal/config"
	"friend_zone/internal/infra/kafka"
	"friend_zone/internal/infra/mysql"
	redisinfra "friend_zone/internal/infra/redis"
	"friend_zone/internal/infra/snowflake"
	"friend_zone/internal/module/fanout"
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

	idgen, err := snowflake.New(2)
	if err != nil {
		log.Fatalf("new snowflake generator: %v", err)
	}
	producer := kafka.NewProducer(cfg.Kafka)
	defer producer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("fanout worker started")
	fanout.NewProcessor(db, rdb, producer, idgen, cfg).Run(ctx)
}
