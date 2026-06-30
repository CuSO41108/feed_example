package user

import (
	"context"
	"database/sql"
	"time"
)

type ActivityService struct {
	db           *sql.DB
	activeWindow time.Duration
}

func NewActivityService(db *sql.DB, activeWindow time.Duration) *ActivityService {
	return &ActivityService{db: db, activeWindow: activeWindow}
}

func (s *ActivityService) ActiveUntil(now time.Time) time.Time {
	return now.Add(s.activeWindow)
}

func (s *ActivityService) MarkLogin(ctx context.Context, userID int64) error {
	now := time.Now().UTC()
	return s.upsert(ctx, userID, "last_login_at", now)
}

func (s *ActivityService) MarkFeedRefresh(ctx context.Context, userID int64) error {
	now := time.Now().UTC()
	return s.upsert(ctx, userID, "last_feed_refresh_at", now)
}

func (s *ActivityService) upsert(ctx context.Context, userID int64, column string, now time.Time) error {
	activeUntil := s.ActiveUntil(now)
	query := "INSERT INTO user_activity (user_id, " + column + ", active_until, updated_at) VALUES (?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " + column + "=VALUES(" + column + "), active_until=VALUES(active_until), updated_at=VALUES(updated_at)"
	_, err := s.db.ExecContext(ctx, query, userID, now, activeUntil, now)
	return err
}

func (s *ActivityService) IsActive(ctx context.Context, userID int64) (bool, error) {
	var activeUntil time.Time
	err := s.db.QueryRowContext(ctx, "SELECT active_until FROM user_activity WHERE user_id = ?", userID).Scan(&activeUntil)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !activeUntil.Before(time.Now().UTC()), nil
}
