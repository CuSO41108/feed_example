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

- Health check: `http://localhost:8080/healthz`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Kafka UI: `http://localhost:8081`

## Scope

This first version focuses on the Timeline Feed main path:

- JWT auth
- Follow and unfollow
- Text-only post publishing
- Logical post deletion
- Author outbox
- User feed inbox
- Pull-to-refresh and infinite scroll with `(publish_time, content_id)` cursor
- Big V push-pull hybrid fanout
- Kafka fanout worker
- Redis ZSET cache for active users' latest 1000 inbox entries

