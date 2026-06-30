package fanout

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"friend_zone/internal/config"
	"friend_zone/internal/infra/kafka"
	"friend_zone/internal/infra/snowflake"
)

type Processor struct {
	db       *sql.DB
	redis    *goredis.Client
	producer *kafka.Producer
	idgen    *snowflake.Generator
	config   config.FeedConfig
	kafkaCfg config.KafkaConfig
}

func NewProcessor(db *sql.DB, redis *goredis.Client, producer *kafka.Producer, idgen *snowflake.Generator, appCfg config.Config) *Processor {
	return &Processor{
		db:       db,
		redis:    redis,
		producer: producer,
		idgen:    idgen,
		config:   appCfg.Feed,
		kafkaCfg: appCfg.Kafka,
	}
}

func (p *Processor) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, topic := range []string{kafka.TopicPostPublished, kafka.TopicPostDeleted, kafka.TopicFeedFanoutChunk, kafka.TopicFeedFanoutRetry} {
		topic := topic
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.consume(ctx, topic)
		}()
	}
	wg.Wait()
}

func (p *Processor) consume(ctx context.Context, topic string) {
	reader := kafka.NewReader(p.kafkaCfg, topic)
	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("fetch kafka message from %s: %v", topic, err)
			time.Sleep(time.Second)
			continue
		}

		if err := p.handleMessage(ctx, topic, msg.Value); err != nil {
			log.Printf("handle kafka message topic=%s offset=%d: %v", topic, msg.Offset, err)
			if topic == kafka.TopicFeedFanoutChunk || topic == kafka.TopicFeedFanoutRetry {
				p.retryChunk(ctx, msg.Value)
			}
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("commit kafka message topic=%s offset=%d: %v", topic, msg.Offset, err)
		}
	}
}

func (p *Processor) handleMessage(ctx context.Context, topic string, payload []byte) error {
	switch topic {
	case kafka.TopicPostPublished:
		var msg PostPublishedMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}
		return p.handlePostPublished(ctx, msg)
	case kafka.TopicPostDeleted:
		var msg PostDeletedMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}
		return p.handlePostDeleted(ctx, msg)
	case kafka.TopicFeedFanoutChunk, kafka.TopicFeedFanoutRetry:
		var msg FanoutChunkMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}
		return p.handleFanoutChunk(ctx, msg)
	default:
		return nil
	}
}

func (p *Processor) handlePostPublished(ctx context.Context, msg PostPublishedMessage) error {
	activeOnly := msg.FollowerCount >= p.config.BigVThreshold
	var lastFollowerID int64

	for {
		followerIDs, err := p.fetchFollowers(ctx, msg.AuthorID, lastFollowerID, activeOnly, p.config.FanoutChunkSize)
		if err != nil {
			return err
		}
		if len(followerIDs) == 0 {
			return nil
		}

		chunk := FanoutChunkMessage{
			TaskID:      p.idgen.NextID(),
			EventID:     msg.EventID,
			ContentID:   msg.ContentID,
			AuthorID:    msg.AuthorID,
			PublishTime: msg.PublishTime,
			FollowerIDs: followerIDs,
		}
		if err := p.producer.PublishJSON(ctx, kafka.TopicFeedFanoutChunk, strconv.FormatInt(chunk.TaskID, 10), chunk); err != nil {
			return err
		}
		lastFollowerID = followerIDs[len(followerIDs)-1]
	}
}

func (p *Processor) handleFanoutChunk(ctx context.Context, msg FanoutChunkMessage) error {
	if len(msg.FollowerIDs) == 0 {
		return nil
	}
	if err := p.insertInboxBatch(ctx, msg); err != nil {
		return err
	}
	activeIDs, err := p.filterActiveUsers(ctx, msg.FollowerIDs)
	if err != nil {
		return err
	}
	return p.writeRedisInbox(ctx, activeIDs, msg)
}

func (p *Processor) handlePostDeleted(ctx context.Context, msg PostDeletedMessage) error {
	var lastFollowerID int64
	for {
		followerIDs, err := p.fetchFollowers(ctx, msg.AuthorID, lastFollowerID, true, p.config.FanoutChunkSize)
		if err != nil {
			return err
		}
		if len(followerIDs) == 0 {
			return nil
		}
		if err := p.removeRedisInbox(ctx, followerIDs, msg); err != nil {
			return err
		}
		lastFollowerID = followerIDs[len(followerIDs)-1]
	}
}

func (p *Processor) fetchFollowers(ctx context.Context, authorID int64, lastFollowerID int64, activeOnly bool, limit int) ([]int64, error) {
	query := `
		SELECT f.follower_id
		FROM follow_relations f`
	args := []any{authorID, lastFollowerID}
	if activeOnly {
		query += `
		JOIN user_activity ua ON ua.user_id = f.follower_id AND ua.active_until >= ?`
		args = []any{time.Now().UTC(), authorID, lastFollowerID}
	}
	query += `
		WHERE f.followee_id = ?
		  AND f.status = 1
		  AND f.follower_id > ?
		ORDER BY f.follower_id ASC
		LIMIT ?`
	args = append(args, limit)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (p *Processor) insertInboxBatch(ctx context.Context, msg FanoutChunkMessage) error {
	placeholders := make([]string, 0, len(msg.FollowerIDs))
	args := make([]any, 0, len(msg.FollowerIDs)*4)
	for _, userID := range msg.FollowerIDs {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		args = append(args, userID, msg.ContentID, msg.AuthorID, msg.PublishTime)
	}
	query := "INSERT IGNORE INTO user_feed_inbox (user_id, content_id, author_id, publish_time) VALUES " + strings.Join(placeholders, ",")
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *Processor) filterActiveUsers(ctx context.Context, userIDs []int64) ([]int64, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(userIDs))
	args := make([]any, 0, len(userIDs)+1)
	args = append(args, time.Now().UTC())
	for i, id := range userIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	query := fmt.Sprintf(`
		SELECT user_id
		FROM user_activity
		WHERE active_until >= ?
		  AND user_id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (p *Processor) writeRedisInbox(ctx context.Context, userIDs []int64, msg FanoutChunkMessage) error {
	if len(userIDs) == 0 || p.redis == nil {
		return nil
	}
	pipe := p.redis.Pipeline()
	member := redisMember(msg.ContentID, msg.AuthorID)
	for _, userID := range userIDs {
		key := redisInboxKey(userID)
		pipe.ZAdd(ctx, key, goredis.Z{Score: float64(msg.PublishTime.UnixMilli()), Member: member})
		pipe.ZRemRangeByRank(ctx, key, 0, -p.config.RedisInboxLimit-1)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (p *Processor) removeRedisInbox(ctx context.Context, userIDs []int64, msg PostDeletedMessage) error {
	if len(userIDs) == 0 || p.redis == nil {
		return nil
	}
	pipe := p.redis.Pipeline()
	member := redisMember(msg.ContentID, msg.AuthorID)
	for _, userID := range userIDs {
		pipe.ZRem(ctx, redisInboxKey(userID), member)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (p *Processor) retryChunk(ctx context.Context, payload []byte) {
	var msg FanoutChunkMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return
	}
	msg.RetryCount++
	topic := kafka.TopicFeedFanoutRetry
	if msg.RetryCount > 3 {
		topic = kafka.TopicFeedFanoutDLQ
	}
	if err := p.producer.PublishJSON(ctx, topic, strconv.FormatInt(msg.TaskID, 10), msg); err != nil {
		log.Printf("publish retry/dlq failed task=%d: %v", msg.TaskID, err)
	}
}

func redisInboxKey(userID int64) string {
	return "feed:inbox:" + strconv.FormatInt(userID, 10)
}

func redisMember(contentID int64, authorID int64) string {
	return fmt.Sprintf("%020d:%d", contentID, authorID)
}
