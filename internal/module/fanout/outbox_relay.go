package fanout

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"friend_zone/internal/config"
	"friend_zone/internal/infra/kafka"
)

type OutboxRelay struct {
	db       *sql.DB
	producer *kafka.Producer
	interval time.Duration
}

type pendingEvent struct {
	EventID int64
	Topic   string
	Payload []byte
}

func NewOutboxRelay(db *sql.DB, producer *kafka.Producer, cfg config.FeedConfig) *OutboxRelay {
	return &OutboxRelay{db: db, producer: producer, interval: cfg.OutboxRelayInterval}
}

func (r *OutboxRelay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		if err := r.drainOnce(ctx); err != nil {
			log.Printf("outbox relay error: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *OutboxRelay) drainOnce(ctx context.Context) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT event_id, topic, CAST(payload AS CHAR) AS payload
		FROM event_outbox
		WHERE status = 0 AND retry_count < 20
		ORDER BY created_at ASC
		LIMIT 100`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var event pendingEvent
		if err := rows.Scan(&event.EventID, &event.Topic, &event.Payload); err != nil {
			return err
		}
		if err := r.producer.PublishBytes(ctx, event.Topic, strconv.FormatInt(event.EventID, 10), event.Payload); err != nil {
			_, _ = r.db.ExecContext(ctx, `
				UPDATE event_outbox
				SET retry_count = retry_count + 1, updated_at = ?
				WHERE event_id = ?`, time.Now().UTC(), event.EventID)
			continue
		}
		_, _ = r.db.ExecContext(ctx, `
			UPDATE event_outbox
			SET status = 1, updated_at = ?
			WHERE event_id = ?`, time.Now().UTC(), event.EventID)
	}
	return rows.Err()
}
