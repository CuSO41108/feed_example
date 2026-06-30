package fanout

import "time"

type PostPublishedMessage struct {
	EventID       int64     `json:"event_id"`
	ContentID     int64     `json:"content_id"`
	AuthorID      int64     `json:"author_id"`
	PublishTime   time.Time `json:"publish_time"`
	FollowerCount int64     `json:"follower_count"`
	BigVThreshold int64     `json:"big_v_threshold"`
}

type PostDeletedMessage struct {
	EventID     int64 `json:"event_id"`
	ContentID   int64 `json:"content_id"`
	AuthorID    int64 `json:"author_id"`
	DeletedByID int64 `json:"deleted_by_id"`
}

type FanoutChunkMessage struct {
	TaskID      int64     `json:"task_id"`
	EventID     int64     `json:"event_id"`
	ContentID   int64     `json:"content_id"`
	AuthorID    int64     `json:"author_id"`
	PublishTime time.Time `json:"publish_time"`
	FollowerIDs []int64   `json:"follower_ids"`
	RetryCount  int       `json:"retry_count"`
}
