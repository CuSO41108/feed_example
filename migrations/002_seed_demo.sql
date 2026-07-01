SET NAMES utf8mb4;
SET time_zone = '+00:00';

SET @seed_now = UTC_TIMESTAMP(3);
SET @seed_password_hash = '$2a$10$0tBLxmQE0wLHwtkFG3XIFOrOm9MOabCFvUIUrRbhJCH2lBV8thHFS';

INSERT INTO users (user_id, username, nickname, avatar_key, password_hash, follower_count, following_count, status, created_at, updated_at)
VALUES
  (100000000000000000, 'demo_reader', '演示读者', 'avatar-20', @seed_password_hash, 0, 20, 1, @seed_now, @seed_now),
  (100000000000000001, 'demo_author_01', '林知夏', 'avatar-01', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000002, 'demo_author_02', '周望舒', 'avatar-02', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000003, 'demo_author_03', '陈禾', 'avatar-03', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000004, 'demo_author_04', '宋一澜', 'avatar-04', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000005, 'demo_author_05', '许青岚', 'avatar-05', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000006, 'demo_author_06', '何砚', 'avatar-06', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000007, 'demo_author_07', '顾南乔', 'avatar-07', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000008, 'demo_author_08', '陆闻舟', 'avatar-08', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000009, 'demo_author_09', '唐星野', 'avatar-09', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000010, 'demo_author_10', '沈微澜', 'avatar-10', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000011, 'demo_author_11', '江予白', 'avatar-11', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000012, 'demo_author_12', '叶清和', 'avatar-12', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000013, 'demo_author_13', '梁时雨', 'avatar-13', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000014, 'demo_author_14', '苏眠', 'avatar-14', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000015, 'demo_author_15', '赵予安', 'avatar-15', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000016, 'demo_author_16', '秦书意', 'avatar-16', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000017, 'demo_author_17', '孟栖迟', 'avatar-17', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000018, 'demo_author_18', '夏远山', 'avatar-18', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000019, 'demo_author_19', '乔北辰', 'avatar-19', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now),
  (100000000000000020, 'demo_author_20', '白鹿鸣', 'avatar-20', @seed_password_hash, 1, 0, 1, @seed_now, @seed_now)
ON DUPLICATE KEY UPDATE
  nickname = VALUES(nickname),
  avatar_key = VALUES(avatar_key),
  password_hash = VALUES(password_hash),
  follower_count = VALUES(follower_count),
  following_count = VALUES(following_count),
  status = VALUES(status),
  updated_at = VALUES(updated_at);

INSERT INTO user_activity (user_id, last_login_at, last_feed_refresh_at, active_until, updated_at)
VALUES
  (100000000000000000, @seed_now, @seed_now, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000001, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000002, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000003, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000004, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000005, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000006, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000007, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000008, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000009, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000010, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000011, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000012, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000013, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000014, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000015, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000016, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000017, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000018, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000019, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now),
  (100000000000000020, @seed_now, NULL, DATE_ADD(@seed_now, INTERVAL 7 DAY), @seed_now)
ON DUPLICATE KEY UPDATE
  last_login_at = VALUES(last_login_at),
  active_until = VALUES(active_until),
  updated_at = VALUES(updated_at);

INSERT INTO follow_relations (follower_id, followee_id, status, created_at, updated_at)
VALUES
  (100000000000000000, 100000000000000001, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000002, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000003, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000004, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000005, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000006, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000007, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000008, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000009, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000010, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000011, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000012, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000013, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000014, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000015, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000016, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000017, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000018, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000019, 1, @seed_now, @seed_now),
  (100000000000000000, 100000000000000020, 1, @seed_now, @seed_now)
ON DUPLICATE KEY UPDATE
  status = VALUES(status),
  updated_at = VALUES(updated_at);

INSERT INTO posts (content_id, author_id, content_text, status, publish_time, created_at, updated_at)
VALUES
  (200000000000000001, 100000000000000001, '今天把关注流的收件箱压测了一轮，分页游标终于稳定了。', 1, DATE_SUB(@seed_now, INTERVAL 5 MINUTE), @seed_now, @seed_now),
  (200000000000000002, 100000000000000002, '午后咖啡配系统设计图，脑子里全是 push 和 pull 的边界。', 1, DATE_SUB(@seed_now, INTERVAL 12 MINUTE), @seed_now, @seed_now),
  (200000000000000003, 100000000000000003, '把 Redis ZSET 的最近 1000 条缓存梳理了一遍，刷新速度很舒服。', 1, DATE_SUB(@seed_now, INTERVAL 18 MINUTE), @seed_now, @seed_now),
  (200000000000000004, 100000000000000004, '今天的随机动态：先写清楚数据模型，再写代码，少走很多弯路。', 1, DATE_SUB(@seed_now, INTERVAL 26 MINUTE), @seed_now, @seed_now),
  (200000000000000005, 100000000000000005, 'Kafka fanout chunk 跑起来后，整个发布链路看起来顺眼多了。', 1, DATE_SUB(@seed_now, INTERVAL 35 MINUTE), @seed_now, @seed_now),
  (200000000000000006, 100000000000000006, '给大 V 做推拉结合时，最关键的是别让普通用户刷新时等太久。', 1, DATE_SUB(@seed_now, INTERVAL 43 MINUTE), @seed_now, @seed_now),
  (200000000000000007, 100000000000000007, '今天的页面背景比纯白舒服，信息流也更像朋友圈了。', 1, DATE_SUB(@seed_now, INTERVAL 51 MINUTE), @seed_now, @seed_now),
  (200000000000000008, 100000000000000008, '调了一下头像入口，关注动作藏在用户头像里，界面清爽很多。', 1, DATE_SUB(@seed_now, INTERVAL 67 MINUTE), @seed_now, @seed_now),
  (200000000000000009, 100000000000000009, '发动态、入 outbox、拆 fanout，串起来之后主链路就完整了。', 1, DATE_SUB(@seed_now, INTERVAL 83 MINUTE), @seed_now, @seed_now),
  (200000000000000010, 100000000000000010, '下拉刷新拿最新，上滑加载拿更早，排序一定要按时间和 ID 双字段。', 1, DATE_SUB(@seed_now, INTERVAL 104 MINUTE), @seed_now, @seed_now),
  (200000000000000011, 100000000000000011, '今天给 demo 数据补了昵称和头像，终于不是一排用户 ID 了。', 1, DATE_SUB(@seed_now, INTERVAL 132 MINUTE), @seed_now, @seed_now),
  (200000000000000012, 100000000000000012, '如果缓存没命中就回 MySQL，第一版这样已经足够稳。', 1, DATE_SUB(@seed_now, INTERVAL 155 MINUTE), @seed_now, @seed_now),
  (200000000000000013, 100000000000000013, '把正文详情延迟到读 Feed 时批量查，收件箱就轻很多。', 1, DATE_SUB(@seed_now, INTERVAL 181 MINUTE), @seed_now, @seed_now),
  (200000000000000014, 100000000000000014, '今天记录：逻辑删除只改 status，读 Feed 时过滤掉 deleted。', 1, DATE_SUB(@seed_now, INTERVAL 219 MINUTE), @seed_now, @seed_now),
  (200000000000000015, 100000000000000015, '本地 docker compose 一键起来，开发体验就会顺很多。', 1, DATE_SUB(@seed_now, INTERVAL 248 MINUTE), @seed_now, @seed_now),
  (200000000000000016, 100000000000000016, '做关注流时，幂等比想象中重要，每个 content 只能进同一个收件箱一次。', 1, DATE_SUB(@seed_now, INTERVAL 296 MINUTE), @seed_now, @seed_now),
  (200000000000000017, 100000000000000017, '今天把 API 文档又补了一点，前后端对字段更有默契。', 1, DATE_SUB(@seed_now, INTERVAL 337 MINUTE), @seed_now, @seed_now),
  (200000000000000018, 100000000000000018, '随机内容也要像真实用户写的，不然 Feed 一眼就是假数据。', 1, DATE_SUB(@seed_now, INTERVAL 389 MINUTE), @seed_now, @seed_now),
  (200000000000000019, 100000000000000019, '发布以后先写发件箱，再异步 fanout 到粉丝收件箱，这条链路很清晰。', 1, DATE_SUB(@seed_now, INTERVAL 451 MINUTE), @seed_now, @seed_now),
  (200000000000000020, 100000000000000020, '晚上收工前再刷一遍 Timeline，最新内容排在最前面才安心。', 1, DATE_SUB(@seed_now, INTERVAL 524 MINUTE), @seed_now, @seed_now)
ON DUPLICATE KEY UPDATE
  content_text = VALUES(content_text),
  status = VALUES(status),
  publish_time = VALUES(publish_time),
  updated_at = VALUES(updated_at);

INSERT INTO author_outbox (author_id, content_id, publish_time)
SELECT author_id, content_id, publish_time
FROM posts
WHERE content_id BETWEEN 200000000000000001 AND 200000000000000020
ON DUPLICATE KEY UPDATE
  publish_time = VALUES(publish_time);

INSERT INTO user_feed_inbox (user_id, content_id, author_id, publish_time)
SELECT 100000000000000000, content_id, author_id, publish_time
FROM posts
WHERE content_id BETWEEN 200000000000000001 AND 200000000000000020
ON DUPLICATE KEY UPDATE
  publish_time = VALUES(publish_time);
