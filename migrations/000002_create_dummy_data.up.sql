INSERT INTO "users" ("username", "email", "password", "is_admin")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi', 0),
  ('admin', 'admin@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi', 1),
  ('user1', 'user1@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi', 0),
  ('user2', 'user2@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi', 0),
  ('user3', 'user3@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi', 0);

INSERT INTO "follow" ("user_id_follower", "user_id_followee")
  VALUES (2, 1), (3, 1), (4, 1);


INSERT INTO "arts" ("cover_url", "name", "description", "creator_id", "price")
  VALUES ('/static/storage/kadoru.jpg', 'the first art', 'just the first art bro.', 1, 0),
  ('/static/storage/kadoru4.jpg', 'second art broo', '', 1, 0);


INSERT INTO "files" ("art_id", "filename", "filetype", "url")
  VALUES (1, 'kadoru.jpg', 'jpg', '/static/storage/kadoru.jpg'),
  (1, 'kadoru2.jpg', 'jpg', '/static/storage/kadoru2.jpg'),
  (1, 'kadoru3.jpg', 'jpg', '/static/storage/kadoru3.jpg'),
  (1, 'kadoru4.jpg', 'jpg', '/static/storage/kadoru4.jpg');

INSERT INTO "users_starred_arts" ("user_id", "art_id")
  VALUES (1, 1);


INSERT INTO "downloaded_arts" ("art_id")
  VALUES (1), (1), (1), (2);


INSERT INTO "tags" ("name")
  VALUES ('tag1'), ('tag2'), ('tag3');


INSERT INTO "arts_tags" ("art_id", "tag_id")
  VALUES (1, 1), (1, 2), (1, 3), (2, 1);

INSERT INTO "codes" ("name", "value", "exp_time")
  VALUES ('GETME100', 100, '2025-12-14T15:04:05Z'),
  ('GETME1000', 1000, '2025-12-14T15:04:05Z'),
  ('GETME10000', 10000, '2025-12-14T15:04:05Z');

------------------------------------------------------------------------------

WITH RECURSIVE cte(x) AS (
  SELECT 1 x
  UNION ALL SELECT x+1 FROM cte WHERE x < 100
)
INSERT INTO "files" ("art_id", "filename", "filetype", "url")
  SELECT x+4, printf("%05d.png", x+4), 'png', printf("/static/storage/genimages/%05d.png", x+4)
  FROM cte;

WITH RECURSIVE cte(x) AS (
  SELECT 1 x
  UNION ALL SELECT x+1 FROM cte WHERE x < 100
)
INSERT INTO "arts" ("cover_url", "name", "description", "creator_id", "price")
  SELECT printf("/static/storage/genimages/%05d.png", x+4),
    printf("dummy art no. %d", x+2),
    'dummy description',
    1,
    x+2
  FROM cte;
