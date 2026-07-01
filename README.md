# Friend Zone

Timeline Feed first-version backend skeleton.

## Stack

- Go + Gin
- MySQL
- Redis
- Kafka
- Docker Compose

## Run

```bash
docker compose up --build
```

API:

- Web app: `http://localhost:8080/`
- Health check: `http://localhost:8080/healthz`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Kafka UI: `http://localhost:8081`

Demo login:

- Username: `demo_reader`
- Password: `demo123456`
- The seed data includes 20 followed authors with nicknames, avatar keys, and random text posts.

## Architecture Design

This project implements a simplified Timeline Feed architecture inspired by Chapter 11 of
`亿级流量系统架构设计与实战`. The current version focuses on keeping the core design direction
correct while avoiding unnecessary production-level complexity for a small-data demo.

The main design idea is push-pull hybrid fanout:

- Normal authors use push mode: after publishing a post, the system asynchronously writes feed index entries into followers' inboxes.
- Big V authors use hybrid mode: active followers receive pushed inbox entries, while inactive followers pull Big V posts from the author outbox when reading the feed.
- The post publishing service does not traverse followers directly. It writes the post and an outbox event, then the Timeline worker consumes Kafka events and performs fanout.
- The user inbox stores only feed index data: `user_id`, `content_id`, `author_id`, and `publish_time`. Full post and author details are assembled on the read path.
- Timeline pagination uses cursor pagination with `(publish_time, content_id)` instead of `offset`, so posts with the same publish time still have stable ordering.
- Inbox writes are idempotent through `UNIQUE(user_id, content_id)` and `INSERT IGNORE`.
- Redis ZSET is used as a hot inbox cache for active users, while MySQL remains the source of truth.

## Completion Check

Implemented in this version:

- JWT auth, follow, and unfollow.
- Text-only post publishing and logical deletion.
- Author outbox and user feed inbox tables.
- Post event outbox and Kafka relay.
- Timeline fanout worker with chunk messages, retry topic, and DLQ topic.
- Pull-to-refresh and infinite scroll using `(publish_time, content_id)` cursors.
- Push-pull hybrid feed reads that merge inbox entries with Big V author outbox entries.
- Redis ZSET cache for active users' latest inbox entries.
- Static frontend for login, posting, follow operations, and Timeline Feed reads.

## Future Optimization

The current version assumes author scale is relatively stable: Big V authors remain Big V authors,
and normal authors remain normal authors. Under this assumption, these production-level details are
intentionally deferred:

- Record `fanout_mode` or `is_big_v_at_publish` on `author_outbox` or `posts`, so historical posts keep the fanout decision made at publish time.
- Increase Kafka partitions and worker parallelism when fanout latency becomes visible.
- Add a fanout task progress table keyed by content ID and task number, recording the latest successfully pushed follower for resumable delivery.
- Replace the current single-SQL Big V outbox query with a stricter heap-based K-way merge if author outboxes are later split across services or shards.
- Add a global feed retention window or maximum scroll depth for database inbox history, so Timeline Feed cannot be scrolled indefinitely.
