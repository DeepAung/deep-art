CREATE TABLE "users" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "username" VARCHAR UNIQUE NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "avatar_url" VARCHAR NOT NULL DEFAULT '',
  "is_admin" BOOLEAN NOT NULL DEFAULT false,
  "coin" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "tokens" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "user_id" INT NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE "oauths" (
  "user_id" INT NOT NULL,
  "provider" VARCHAR NOT NULL,
  "provider_user_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("provider", "provider_user_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- follower is follwing followee
-- follower is users
-- followee is creators
CREATE TABLE "follow" (
  "user_id_follower" INT NOT NULL,
  "user_id_followee" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("user_id_follower", "user_id_followee"),
  FOREIGN KEY ("user_id_follower") REFERENCES "users" ("id"),
  FOREIGN KEY ("user_id_followee") REFERENCES "users" ("id")
);

CREATE TABLE "arts" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "cover_url" VARCHAR UNIQUE NOT NULL,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "creator_id" INT NOT NULL,
  "price" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE ("name", "creator_id"),
  FOREIGN KEY ("creator_id") REFERENCES "users" ("id")
);

CREATE TABLE "downloaded_arts" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "art_id" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("art_id") REFERENCES "arts" ("id")
);

CREATE TABLE "users_starred_arts" (
  "user_id" INT NOT NULL,
  "art_id" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("user_id", "art_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("art_id") REFERENCES "arts" ("id")
);

CREATE TABLE "users_bought_arts" (
  "user_id" INT NOT NULL,
  "art_id" INT NOT NULL,
  PRIMARY KEY ("user_id", "art_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("art_id") REFERENCES "arts" ("id")
);

CREATE TABLE "tags" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE "arts_tags" (
  "art_id" INT NOT NULL,
  "tag_id" INT NOT NULL,
  PRIMARY KEY ("art_id", "tag_id"),
  FOREIGN KEY ("art_id") REFERENCES "arts" ("id"),
  FOREIGN KEY ("tag_id") REFERENCES "tags" ("id")
);

CREATE TABLE "files" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "art_id" INT NOT NULL,
  "filename" VARCHAR NOT NULL,
  "filetype" VARCHAR NOT NULL,
  "url" VARCHAR NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("art_id") REFERENCES "arts" ("id")
);

CREATE TABLE "codes" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" VARCHAR UNIQUE NOT NULL,
  "value" INT NOT NULL,
  "exp_time" TIMESTAMP NOT NULL
);

CREATE TABLE "users_used_codes" (
  "user_id" INT NOT NULL,
  "code_id" INT NOT NULL,
  PRIMARY KEY ("user_id", "code_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("code_id") REFERENCES "codes" ("id")
);

CREATE TRIGGER [update_timestamp_users] AFTER UPDATE ON "users" FOR EACH ROW WHEN NEW."updated_at" < OLD."updated_at"
BEGIN UPDATE "users" SET "updated_at"=CURRENT_TIMESTAMP WHERE id=OLD.id; END;

CREATE TRIGGER [update_timestamp_tokens] AFTER UPDATE ON "tokens" FOR EACH ROW WHEN NEW."updated_at" < OLD."updated_at"
BEGIN UPDATE "tokens" SET "updated_at"=CURRENT_TIMESTAMP WHERE id=OLD.id; END;

CREATE TRIGGER [update_timestamp_oauths] AFTER UPDATE ON "oauths" FOR EACH ROW WHEN NEW."updated_at" < OLD."updated_at"
BEGIN UPDATE "oauths" SET "updated_at"=CURRENT_TIMESTAMP WHERE id=OLD.id; END;

CREATE TRIGGER [update_timestamp_arts] AFTER UPDATE ON "arts" FOR EACH ROW WHEN NEW."updated_at" < OLD."updated_at"
BEGIN UPDATE "arts" SET "updated_at"=CURRENT_TIMESTAMP WHERE id=OLD.id; END;

CREATE TRIGGER [update_timestamp_files] AFTER UPDATE ON "files" FOR EACH ROW WHEN NEW."updated_at" < OLD."updated_at"
BEGIN UPDATE "files" SET "updated_at"=CURRENT_TIMESTAMP WHERE id=OLD.id; END;
