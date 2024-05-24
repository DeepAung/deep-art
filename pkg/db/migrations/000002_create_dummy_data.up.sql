INSERT INTO "users" ("username", "email", "password")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '$2a$10$PAsFv3cmdUQPfFvFkefBEOtPAVYnvL9wkyUw5VLDskdBPKayQjagi');

INSERT INTO "arts" ("cover_id", "name", "description", "creator_id", "price")
  VALUES (1, 'the first art', 'just the first art bro.', 1, 0);

INSERT INTO "files" ("art_id", "filename", "filetype", "url")
  VALUES (1, 'kadoru.jpg', 'jpg', './static/storage/kadoru.jpg');

INSERT INTO "files" ("art_id", "filename", "filetype", "url")
  VALUES (1, 'kadoru2.jpg', 'jpg', './static/storage/kadoru2.jpg');

INSERT INTO "files" ("art_id", "filename", "filetype", "url")
  VALUES (1, 'kadoru3.jpg', 'jpg', './static/storage/kadoru3.jpg');

INSERT INTO "tags" ("name")
  VALUES ('tag1'), ('tag2'), ('tag3');

INSERT INTO "arts_tags" ("art_id", "tag_id")
  VALUES (1, 1), (1, 2), (1, 3);
