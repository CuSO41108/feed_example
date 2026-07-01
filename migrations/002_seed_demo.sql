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
SELECT readers.user_id, authors.user_id, 1, @seed_now, @seed_now
FROM users AS readers
JOIN users AS authors ON authors.username LIKE 'demo_author\_%'
WHERE readers.status = 1
  AND authors.status = 1
  AND (readers.username = 'demo_reader' OR readers.username LIKE '%reader%')
  AND readers.user_id <> authors.user_id
ON DUPLICATE KEY UPDATE
  status = VALUES(status),
  updated_at = VALUES(updated_at);

UPDATE users AS u
LEFT JOIN (
  SELECT followee_id AS user_id, COUNT(*) AS follower_count
  FROM follow_relations
  WHERE status = 1
  GROUP BY followee_id
) AS followers ON followers.user_id = u.user_id
LEFT JOIN (
  SELECT follower_id AS user_id, COUNT(*) AS following_count
  FROM follow_relations
  WHERE status = 1
  GROUP BY follower_id
) AS following ON following.user_id = u.user_id
SET u.follower_count = COALESCE(followers.follower_count, 0),
    u.following_count = COALESCE(following.following_count, 0),
    u.updated_at = @seed_now
WHERE u.username = 'demo_reader'
   OR u.username LIKE '%reader%'
   OR u.username LIKE 'demo_author\_%';

INSERT INTO posts (content_id, author_id, content_text, status, publish_time, created_at, updated_at)
VALUES
  (200000000000000001, 100000000000000001, '早上路过楼下花店，老板把向日葵摆到门口，整条街都亮了一点。', 1, DATE_SUB(@seed_now, INTERVAL 5 MINUTE), @seed_now, @seed_now),
  (200000000000000002, 100000000000000002, '今天的早餐是豆浆、油条和一个刚出锅的茶叶蛋，热乎乎的很满足。', 1, DATE_SUB(@seed_now, INTERVAL 12 MINUTE), @seed_now, @seed_now),
  (200000000000000003, 100000000000000003, '傍晚去河边散步，风把桂花味吹了一路，走着走着心情就慢下来了。', 1, DATE_SUB(@seed_now, INTERVAL 18 MINUTE), @seed_now, @seed_now),
  (200000000000000004, 100000000000000004, '地铁上看到一个小朋友认真给玩具熊系安全带，今天的小确幸达成。', 1, DATE_SUB(@seed_now, INTERVAL 26 MINUTE), @seed_now, @seed_now),
  (200000000000000005, 100000000000000005, '午休去买咖啡，店员在杯套上画了一个笑脸，下午开工都轻松一点。', 1, DATE_SUB(@seed_now, INTERVAL 35 MINUTE), @seed_now, @seed_now),
  (200000000000000006, 100000000000000006, '下班顺手买了半个西瓜，冰箱里一放，今晚的快乐就有着落了。', 1, DATE_SUB(@seed_now, INTERVAL 43 MINUTE), @seed_now, @seed_now),
  (200000000000000007, 100000000000000007, '今天云很好看，像一大块被太阳晒软的棉花糖。', 1, DATE_SUB(@seed_now, INTERVAL 51 MINUTE), @seed_now, @seed_now),
  (200000000000000008, 100000000000000008, '周末把阳台整理了一下，薄荷长得特别精神，靠近就有清凉味。', 1, DATE_SUB(@seed_now, INTERVAL 67 MINUTE), @seed_now, @seed_now),
  (200000000000000009, 100000000000000009, '晚饭做了番茄牛腩，汤汁拌米饭太香了，明天还想继续吃。', 1, DATE_SUB(@seed_now, INTERVAL 83 MINUTE), @seed_now, @seed_now),
  (200000000000000010, 100000000000000010, '雨停以后空气特别干净，路灯下面还能看到一点点水汽。', 1, DATE_SUB(@seed_now, INTERVAL 104 MINUTE), @seed_now, @seed_now),
  (200000000000000011, 100000000000000011, '给书桌换了一个小台灯，暖黄色一开，房间突然像周五晚上。', 1, DATE_SUB(@seed_now, INTERVAL 132 MINUTE), @seed_now, @seed_now),
  (200000000000000012, 100000000000000012, '今天没有赶时间，慢慢吃完一碗面，连汤都觉得刚刚好。', 1, DATE_SUB(@seed_now, INTERVAL 155 MINUTE), @seed_now, @seed_now),
  (200000000000000013, 100000000000000013, '路边新开了一家面包店，黄油香从门口飘出来，没忍住买了两个。', 1, DATE_SUB(@seed_now, INTERVAL 181 MINUTE), @seed_now, @seed_now),
  (200000000000000014, 100000000000000014, '晚上洗完衣服晒到阳台，风吹起来的时候有种很安心的生活感。', 1, DATE_SUB(@seed_now, INTERVAL 219 MINUTE), @seed_now, @seed_now),
  (200000000000000015, 100000000000000015, '今天给自己买了一束小雏菊，放在餐桌上，吃饭都变得认真了。', 1, DATE_SUB(@seed_now, INTERVAL 248 MINUTE), @seed_now, @seed_now),
  (200000000000000016, 100000000000000016, '跑步回来喝了一大杯冰水，耳机里刚好播到喜欢的歌。', 1, DATE_SUB(@seed_now, INTERVAL 296 MINUTE), @seed_now, @seed_now),
  (200000000000000017, 100000000000000017, '楼下小猫今天终于肯靠近一点了，蹲在台阶上看了我三秒。', 1, DATE_SUB(@seed_now, INTERVAL 337 MINUTE), @seed_now, @seed_now),
  (200000000000000018, 100000000000000018, '把冰箱里剩下的菜都做成了炒饭，意外很好吃，粒粒分明。', 1, DATE_SUB(@seed_now, INTERVAL 389 MINUTE), @seed_now, @seed_now),
  (200000000000000019, 100000000000000019, '傍晚骑车经过学校操场，听到有人在练吉他，夏天好像提前到了。', 1, DATE_SUB(@seed_now, INTERVAL 451 MINUTE), @seed_now, @seed_now),
  (200000000000000020, 100000000000000020, '睡前把明天要穿的衣服搭好了，感觉已经提前赢了一点明天。', 1, DATE_SUB(@seed_now, INTERVAL 524 MINUTE), @seed_now, @seed_now)
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
SELECT readers.user_id, posts.content_id, posts.author_id, posts.publish_time
FROM users AS readers
JOIN posts ON posts.content_id BETWEEN 200000000000000001 AND 200000000000000020
WHERE readers.status = 1
  AND (readers.username = 'demo_reader' OR readers.username LIKE '%reader%')
ON DUPLICATE KEY UPDATE
  publish_time = VALUES(publish_time);
