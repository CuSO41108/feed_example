package follow

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrCannotFollowSelf = errors.New("cannot follow yourself")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Follow(ctx context.Context, followerID int64, followeeID int64) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status int
	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM follow_relations
		WHERE follower_id = ? AND followee_id = ?
		FOR UPDATE`, followerID, followeeID).Scan(&status)
	if err == sql.ErrNoRows {
		now := time.Now().UTC()
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO follow_relations (follower_id, followee_id, status, created_at, updated_at)
			VALUES (?, ?, 1, ?, ?)`, followerID, followeeID, now, now); err != nil {
			return err
		}
		return s.applyCountDelta(ctx, tx, followerID, followeeID, 1)
	}
	if err != nil {
		return err
	}
	if status == 1 {
		return tx.Commit()
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE follow_relations
		SET status = 1, updated_at = ?
		WHERE follower_id = ? AND followee_id = ?`, time.Now().UTC(), followerID, followeeID); err != nil {
		return err
	}
	return s.applyCountDelta(ctx, tx, followerID, followeeID, 1)
}

func (s *Service) Unfollow(ctx context.Context, followerID int64, followeeID int64) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status int
	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM follow_relations
		WHERE follower_id = ? AND followee_id = ?
		FOR UPDATE`, followerID, followeeID).Scan(&status)
	if err == sql.ErrNoRows || status == 0 {
		return tx.Commit()
	}
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE follow_relations
		SET status = 0, updated_at = ?
		WHERE follower_id = ? AND followee_id = ?`, time.Now().UTC(), followerID, followeeID); err != nil {
		return err
	}
	return s.applyCountDelta(ctx, tx, followerID, followeeID, -1)
}

func (s *Service) applyCountDelta(ctx context.Context, tx *sql.Tx, followerID int64, followeeID int64, delta int64) error {
	if delta > 0 {
		if _, err := tx.ExecContext(ctx, `
			UPDATE users
			SET following_count = following_count + 1, updated_at = ?
			WHERE user_id = ?`, time.Now().UTC(), followerID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE users
			SET follower_count = follower_count + 1, updated_at = ?
			WHERE user_id = ?`, time.Now().UTC(), followeeID); err != nil {
			return err
		}
		return tx.Commit()
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET following_count = GREATEST(following_count - 1, 0), updated_at = ?
		WHERE user_id = ?`, time.Now().UTC(), followerID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET follower_count = GREATEST(follower_count - 1, 0), updated_at = ?
		WHERE user_id = ?`, time.Now().UTC(), followeeID); err != nil {
		return err
	}
	return tx.Commit()
}
