package feed

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"friend_zone/internal/config"
	activity "friend_zone/internal/module/user"
	cursorpkg "friend_zone/internal/pkg/cursor"
)

const (
	DirectionLatest = "latest"
	DirectionNewer  = "newer"
	DirectionOlder  = "older"
)

type Service struct {
	db       *sql.DB
	redis    *goredis.Client
	activity *activity.ActivityService
	config   config.FeedConfig
}

type Query struct {
	Direction string
	Cursor    string
	Limit     int
}

type Item struct {
	ContentID   int64     `json:"content_id"`
	AuthorID    int64     `json:"author_id"`
	ContentText string    `json:"content_text"`
	PublishTime time.Time `json:"publish_time"`
}

type Response struct {
	Items        []Item `json:"items"`
	TopCursor    string `json:"top_cursor,omitempty"`
	BottomCursor string `json:"bottom_cursor,omitempty"`
	HasMore      bool   `json:"has_more"`
}

type indexEntry struct {
	ContentID   int64
	AuthorID    int64
	PublishTime time.Time
}

func NewService(db *sql.DB, redis *goredis.Client, activity *activity.ActivityService, cfg config.FeedConfig) *Service {
	return &Service{db: db, redis: redis, activity: activity, config: cfg}
}

func (s *Service) Timeline(ctx context.Context, userID int64, q Query) (Response, error) {
	if q.Direction == "" {
		q.Direction = DirectionLatest
	}
	limit := s.normalizeLimit(q.Limit)

	var cur cursorpkg.Cursor
	var hasCursor bool
	if q.Cursor != "" {
		decoded, err := cursorpkg.Decode(q.Cursor)
		if err != nil {
			return Response{}, err
		}
		cur = decoded
		hasCursor = true
	}

	if err := s.activity.MarkFeedRefresh(ctx, userID); err != nil {
		return Response{}, err
	}

	candidateLimit := limit * 5
	if candidateLimit < limit+1 {
		candidateLimit = limit + 1
	}

	inboxEntries, err := s.fetchInbox(ctx, userID, q.Direction, cur, hasCursor, candidateLimit)
	if err != nil {
		return Response{}, err
	}
	bigVEntries, err := s.fetchBigVOutbox(ctx, userID, q.Direction, cur, hasCursor, candidateLimit)
	if err != nil {
		return Response{}, err
	}

	entries := mergeEntries(inboxEntries, bigVEntries)
	if len(entries) == 0 {
		return Response{Items: []Item{}}, nil
	}

	details, err := s.fetchPostDetails(ctx, userID, entries)
	if err != nil {
		return Response{}, err
	}

	items := make([]Item, 0, limit+1)
	for _, entry := range entries {
		item, ok := details[entry.ContentID]
		if !ok {
			continue
		}
		items = append(items, item)
		if len(items) == limit+1 {
			break
		}
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	resp := Response{Items: items, HasMore: hasMore}
	if len(items) > 0 {
		top, _ := cursorpkg.Encode(cursorpkg.Cursor{PublishTime: items[0].PublishTime, ContentID: items[0].ContentID})
		bottom, _ := cursorpkg.Encode(cursorpkg.Cursor{PublishTime: items[len(items)-1].PublishTime, ContentID: items[len(items)-1].ContentID})
		resp.TopCursor = top
		resp.BottomCursor = bottom
	}
	return resp, nil
}

func (s *Service) fetchInbox(ctx context.Context, userID int64, direction string, cur cursorpkg.Cursor, hasCursor bool, limit int) ([]indexEntry, error) {
	if s.redis != nil && direction == DirectionLatest && !hasCursor {
		entries, err := s.fetchInboxFromRedis(ctx, userID, limit)
		if err == nil && len(entries) >= limit {
			return entries, nil
		}
	}

	query := `
		SELECT content_id, author_id, publish_time
		FROM user_feed_inbox
		WHERE user_id = ?`
	args := []any{userID}
	query, args = addCursorCondition(query, args, direction, cur, hasCursor)
	query += `
		ORDER BY publish_time DESC, content_id DESC
		LIMIT ?`
	args = append(args, limit)
	return scanEntries(ctx, s.db, query, args...)
}

func (s *Service) fetchInboxFromRedis(ctx context.Context, userID int64, limit int) ([]indexEntry, error) {
	values, err := s.redis.ZRevRangeWithScores(ctx, redisInboxKey(userID), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	entries := make([]indexEntry, 0, len(values))
	for _, value := range values {
		member, ok := value.Member.(string)
		if !ok {
			continue
		}
		contentID, authorID, ok := parseRedisMember(member)
		if !ok {
			continue
		}
		entries = append(entries, indexEntry{
			ContentID:   contentID,
			AuthorID:    authorID,
			PublishTime: time.UnixMilli(int64(value.Score)).UTC(),
		})
	}
	return entries, nil
}

func (s *Service) fetchBigVOutbox(ctx context.Context, userID int64, direction string, cur cursorpkg.Cursor, hasCursor bool, limit int) ([]indexEntry, error) {
	minTime := time.Now().UTC().Add(-s.config.BigVPullWindow)
	query := `
		SELECT ao.content_id, ao.author_id, ao.publish_time
		FROM author_outbox ao
		JOIN users u ON u.user_id = ao.author_id
		JOIN follow_relations fr ON fr.followee_id = ao.author_id
		WHERE fr.follower_id = ?
		  AND fr.status = 1
		  AND u.follower_count >= ?
		  AND ao.publish_time >= ?`
	args := []any{userID, s.config.BigVThreshold, minTime}
	query, args = addCursorConditionWithPrefix(query, args, "ao.", direction, cur, hasCursor)
	query += `
		ORDER BY ao.publish_time DESC, ao.content_id DESC
		LIMIT ?`
	args = append(args, limit)
	return scanEntries(ctx, s.db, query, args...)
}

func (s *Service) fetchPostDetails(ctx context.Context, userID int64, entries []indexEntry) (map[int64]Item, error) {
	ids := make([]int64, 0, len(entries))
	seen := make(map[int64]struct{}, len(entries))
	for _, entry := range entries {
		if _, ok := seen[entry.ContentID]; ok {
			continue
		}
		seen[entry.ContentID] = struct{}{}
		ids = append(ids, entry.ContentID)
	}
	if len(ids) == 0 {
		return map[int64]Item{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, 0, len(ids)+2)
	args = append(args, userID)
	for i, id := range ids {
		placeholders[i] = "?"
		args = append(args, id)
	}
	args = append(args, userID)

	query := fmt.Sprintf(`
		SELECT p.content_id, p.author_id, p.content_text, p.publish_time
		FROM posts p
		JOIN follow_relations fr ON fr.followee_id = p.author_id
		WHERE fr.follower_id = ?
		  AND fr.status = 1
		  AND p.status = 1
		  AND p.content_id IN (%s)
		  AND p.author_id <> ?`, strings.Join(placeholders, ","))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[int64]Item, len(ids))
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ContentID, &item.AuthorID, &item.ContentText, &item.PublishTime); err != nil {
			return nil, err
		}
		out[item.ContentID] = item
	}
	return out, rows.Err()
}

func scanEntries(ctx context.Context, db *sql.DB, query string, args ...any) ([]indexEntry, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []indexEntry{}
	for rows.Next() {
		var entry indexEntry
		if err := rows.Scan(&entry.ContentID, &entry.AuthorID, &entry.PublishTime); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func mergeEntries(groups ...[]indexEntry) []indexEntry {
	seen := map[int64]indexEntry{}
	for _, group := range groups {
		for _, entry := range group {
			old, exists := seen[entry.ContentID]
			if !exists || entry.PublishTime.After(old.PublishTime) {
				seen[entry.ContentID] = entry
			}
		}
	}
	out := make([]indexEntry, 0, len(seen))
	for _, entry := range seen {
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].PublishTime.Equal(out[j].PublishTime) {
			return out[i].ContentID > out[j].ContentID
		}
		return out[i].PublishTime.After(out[j].PublishTime)
	})
	return out
}

func addCursorCondition(query string, args []any, direction string, cur cursorpkg.Cursor, hasCursor bool) (string, []any) {
	return addCursorConditionWithPrefix(query, args, "", direction, cur, hasCursor)
}

func addCursorConditionWithPrefix(query string, args []any, prefix string, direction string, cur cursorpkg.Cursor, hasCursor bool) (string, []any) {
	if !hasCursor || direction == DirectionLatest {
		return query, args
	}
	switch direction {
	case DirectionNewer:
		query += " AND (" + prefix + "publish_time > ? OR (" + prefix + "publish_time = ? AND " + prefix + "content_id > ?))"
	case DirectionOlder:
		query += " AND (" + prefix + "publish_time < ? OR (" + prefix + "publish_time = ? AND " + prefix + "content_id < ?))"
	default:
		return query, args
	}
	args = append(args, cur.PublishTime, cur.PublishTime, cur.ContentID)
	return query, args
}

func (s *Service) normalizeLimit(limit int) int {
	if limit <= 0 {
		return s.config.DefaultPageSize
	}
	if limit > s.config.MaxPageSize {
		return s.config.MaxPageSize
	}
	return limit
}

func redisInboxKey(userID int64) string {
	return "feed:inbox:" + strconv.FormatInt(userID, 10)
}

func parseRedisMember(member string) (int64, int64, bool) {
	parts := strings.Split(member, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	contentID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	authorID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	return contentID, authorID, true
}
