package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App   AppConfig
	HTTP  HTTPConfig
	MySQL MySQLConfig
	Redis RedisConfig
	Kafka KafkaConfig
	JWT   JWTConfig
	Feed  FeedConfig
}

type AppConfig struct {
	Env string
}

type HTTPConfig struct {
	Addr string
}

type MySQLConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type FeedConfig struct {
	BigVThreshold       int64
	ActiveWindow        time.Duration
	BigVPullWindow      time.Duration
	RedisInboxLimit     int64
	FanoutChunkSize     int
	DefaultPageSize     int
	MaxPageSize         int
	OutboxRelayInterval time.Duration
}

func Load() Config {
	return Config{
		App: AppConfig{
			Env: getEnv("APP_ENV", "dev"),
		},
		HTTP: HTTPConfig{
			Addr: getEnv("HTTP_ADDR", ":8080"),
		},
		MySQL: MySQLConfig{
			DSN: getEnv("MYSQL_DSN", "friend_zone:friend_zone@tcp(127.0.0.1:3306)/friend_zone?parseTime=true&loc=UTC&charset=utf8mb4,utf8"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers: splitCSV(getEnv("KAFKA_BROKERS", "127.0.0.1:9092")),
			GroupID: getEnv("KAFKA_GROUP_ID", "friend-zone-worker"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "dev-secret-change-me"),
			TTL:    time.Duration(getEnvInt("JWT_TTL_HOURS", 24)) * time.Hour,
		},
		Feed: FeedConfig{
			BigVThreshold:       int64(getEnvInt("BIG_V_THRESHOLD", 10000)),
			ActiveWindow:        time.Duration(getEnvInt("ACTIVE_WINDOW_HOURS", 24*7)) * time.Hour,
			BigVPullWindow:      time.Duration(getEnvInt("BIG_V_PULL_WINDOW_HOURS", 24*30)) * time.Hour,
			RedisInboxLimit:     int64(getEnvInt("REDIS_INBOX_LIMIT", 1000)),
			FanoutChunkSize:     getEnvInt("FANOUT_CHUNK_SIZE", 500),
			DefaultPageSize:     getEnvInt("DEFAULT_PAGE_SIZE", 20),
			MaxPageSize:         getEnvInt("MAX_PAGE_SIZE", 100),
			OutboxRelayInterval: time.Duration(getEnvInt("OUTBOX_RELAY_INTERVAL_SECONDS", 2)) * time.Second,
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
