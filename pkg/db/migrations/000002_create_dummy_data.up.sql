INSERT INTO "users" ("username", "email", "password")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi');


INSERT INTO "arts" ("cover_id", "name", "description", "creator_id", "price")
  VALUES (1, 'the first art', 'just the first art bro.', 1, 0),
  (4, 'second art broo', '', 1, 0);


INSERT INTO "files" ("art_id", "filename", "filetype", "url")
  VALUES (1, 'kadoru.jpg', 'jpg', './static/storage/kadoru.jpg'),
  (1, 'kadoru2.jpg', 'jpg', './static/storage/kadoru2.jpg'),
  (1, 'kadoru3.jpg', 'jpg', './static/storage/kadoru3.jpg'),
  (1, 'kadoru4.jpg', 'jpg', './static/storage/kadoru4.jpg');


INSERT INTO "users_starred_arts" ("user_id", "art_id")
  VALUES (1, 1);


INSERT INTO "downloaded_arts" ("art_id")
  VALUES (1), (1), (1), (2);


INSERT INTO "tags" ("name")
  VALUES ('tag1'), ('tag2'), ('tag3');


INSERT INTO "arts_tags" ("art_id", "tag_id")
  VALUES (1, 1), (1, 2), (1, 3);

INSERT INTO "codes" ("name", "value", "exp_time")
  VALUES ('GETME100', 100, '2025-12-14T15:04:05Z'),
  ('GETME1000', 1000, '2025-12-14T15:04:05Z'),
  ('GETME10000', 10000, '2025-12-14T15:04:05Z');
