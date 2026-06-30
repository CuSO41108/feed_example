package post

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"friend_zone/internal/config"
	"friend_zone/internal/infra/kafka"
	"friend_zone/internal/infra/snowflake"
	"friend_zone/internal/module/fanout"
)

const (
	StatusNormal  = 1
	StatusDeleted = 2
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
	ErrEmptyContent = errors.New("content_text is required")
)

type Service struct {
	db     *sql.DB
	idgen  *snowflake.Generator
	config config.FeedConfig
}

type CreateRequest struct {
	ContentText string `json:"content_text" binding:"required,max=2000"`
}

type Post struct {
	ContentID   int64     `json:"content_id"`
	AuthorID    int64     `json:"author_id"`
	ContentText string    `json:"content_text"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publish_time"`
}

func NewService(db *sql.DB, idgen *snowflake.Generator, cfg config.FeedConfig) *Service {
	return &Service{db: db, idgen: idgen, config: cfg}
}

func (s *Service) Create(ctx context.Context, authorID int64, req CreateRequest) (Post, error) {
	if req.ContentText == "" {
		return Post{}, ErrEmptyContent
	}

	contentID := s.idgen.NextID()
	eventID := s.idgen.NextID()
	now := time.Now().UTC()
	post := Post{
		ContentID:   contentID,
		AuthorID:    authorID,
		ContentText: req.ContentText,
		Status:      StatusNormal,
		PublishTime: now,
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Post{}, err
	}
	defer tx.Rollback()

	var followerCount int64
	if err := tx.QueryRowContext(ctx, `SELECT follower_count FROM users WHERE user_id = ? AND status = 1`, authorID).Scan(&followerCount); err != nil {
		return Post{}, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO posts (content_id, author_id, content_text, status, publish_time, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		contentID, authorID, req.ContentText, StatusNormal, now, now, now); err != nil {
		return Post{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO author_outbox (author_id, content_id, publish_time)
		VALUES (?, ?, ?)`, authorID, contentID, now); err != nil {
		return Post{}, err
	}

	message := fanout.PostPublishedMessage{
		EventID:       eventID,
		ContentID:     contentID,
		AuthorID:      authorID,
		PublishTime:   now,
		FollowerCount: followerCount,
		BigVThreshold: s.config.BigVThreshold,
	}
	if err := insertOutboxEvent(ctx, tx, eventID, kafka.TopicPostPublished, message, now); err != nil {
		return Post{}, err
	}

	if err := tx.Commit(); err != nil {
		return Post{}, err
	}
	return post, nil
}

func (s *Service) Delete(ctx context.Context, operatorID int64, contentID int64) error {
	eventID := s.idgen.NextID()
	now := time.Now().UTC()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var authorID int64
	var status int
	err = tx.QueryRowContext(ctx, `
		SELECT author_id, status
		FROM posts
		WHERE content_id = ?
		FOR UPDATE`, contentID).Scan(&authorID, &status)
	if err == sql.ErrNoRows {
		return ErrPostNotFound
	}
	if err != nil {
		return err
	}
	if authorID != operatorID {
		return ErrForbidden
	}
	if status == StatusDeleted {
		return tx.Commit()
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE posts
		SET status = ?, updated_at = ?
		WHERE content_id = ?`, StatusDeleted, now, contentID); err != nil {
		return err
	}
	message := fanout.PostDeletedMessage{
		EventID:     eventID,
		ContentID:   contentID,
		AuthorID:    authorID,
		DeletedByID: operatorID,
	}
	if err := insertOutboxEvent(ctx, tx, eventID, kafka.TopicPostDeleted, message, now); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) Get(ctx context.Context, contentID int64) (Post, error) {
	var p Post
	err := s.db.QueryRowContext(ctx, `
		SELECT content_id, author_id, content_text, status, publish_time
		FROM posts
		WHERE content_id = ?`, contentID).Scan(&p.ContentID, &p.AuthorID, &p.ContentText, &p.Status, &p.PublishTime)
	if err == sql.ErrNoRows {
		return Post{}, ErrPostNotFound
	}
	if err != nil {
		return Post{}, err
	}
	return p, nil
}

func insertOutboxEvent(ctx context.Context, tx *sql.Tx, eventID int64, topic string, payload any, now time.Time) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO event_outbox (event_id, topic, payload, status, retry_count, created_at, updated_at)
		VALUES (?, ?, ?, 0, 0, ?, ?)`,
		eventID, topic, string(data), now, now)
	return err
}
