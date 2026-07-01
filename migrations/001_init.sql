SET NAMES utf8mb4;
SET time_zone = '+00:00';

CREATE TABLE IF NOT EXISTS users (
  user_id BIGINT NOT NULL,
  username VARCHAR(64) NOT NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  avatar_key VARCHAR(32) NOT NULL DEFAULT 'avatar-01',
  password_hash VARCHAR(255) NOT NULL,
  follower_count BIGINT NOT NULL DEFAULT 0,
  following_count BIGINT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (user_id),
  UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_activity (
  user_id BIGINT NOT NULL,
  last_login_at DATETIME(3) NULL,
  last_feed_refresh_at DATETIME(3) NULL,
  active_until DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (user_id),
  KEY idx_active_until (active_until)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS follow_relations (
  follower_id BIGINT NOT NULL,
  followee_id BIGINT NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (follower_id, followee_id),
  KEY idx_followee_status_follower (followee_id, status, follower_id),
  KEY idx_follower_status_followee (follower_id, status, followee_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS posts (
  content_id BIGINT NOT NULL,
  author_id BIGINT NOT NULL,
  content_text TEXT NOT NULL,
  status TINYINT NOT NULL DEFAULT 1,
  publish_time DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (content_id),
  KEY idx_author_time_id (author_id, publish_time, content_id),
  KEY idx_time_id (publish_time, content_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS author_outbox (
  author_id BIGINT NOT NULL,
  content_id BIGINT NOT NULL,
  publish_time DATETIME(3) NOT NULL,
  PRIMARY KEY (author_id, publish_time, content_id),
  UNIQUE KEY uk_content_id (content_id),
  KEY idx_time_id (publish_time, content_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_feed_inbox (
  user_id BIGINT NOT NULL,
  content_id BIGINT NOT NULL,
  author_id BIGINT NOT NULL,
  publish_time DATETIME(3) NOT NULL,
  PRIMARY KEY (user_id, publish_time, content_id),
  UNIQUE KEY uk_user_content (user_id, content_id),
  KEY idx_user_author_time (user_id, author_id, publish_time, content_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS event_outbox (
  event_id BIGINT NOT NULL,
  topic VARCHAR(128) NOT NULL,
  payload JSON NOT NULL,
  status TINYINT NOT NULL DEFAULT 0,
  retry_count INT NOT NULL DEFAULT 0,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (event_id),
  KEY idx_status_created (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
