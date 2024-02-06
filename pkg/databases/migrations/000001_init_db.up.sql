BEGIN;

-- set timezone
SET TIME ZONE 'Asia/Bangkok';

-- auto update
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- create enum
CREATE TYPE "social" AS ENUM (
  'google',
  'github'
);

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR UNIQUE NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "avatar_url" VARCHAR NOT NULL DEFAULT '',
  "is_admin" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE "tokens" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE "oauths" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "social" social NOT NULL,
  "social_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  UNIQUE ("social", "social_id"),
  UNIQUE ("user_id", "social")
);

CREATE TABLE "arts" (
  "id" SERIAL PRIMARY KEY,
  "cover_id" INT UNIQUE NOT NULL,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "creator_id" INT NOT NULL,
  "price" FLOAT NOT NULL DEFAULT 0,
  "download_count" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  UNIQUE ("name", "creator_id")
);

CREATE TABLE "users_starred_arts" (
  "user_id" INT NOT NULL,
  "art_id" INT NOT NULL,
  PRIMARY KEY ("user_id", "art_id")
);

CREATE TABLE "users_bought_arts" (
  "user_id" INT NOT NULL,
  "art_id" INT NOT NULL,
  PRIMARY KEY ("user_id", "art_id")
);

CREATE TABLE "users_created_arts" (
  "user_id" INT NOT NULL,
  "art_id" INT NOT NULL,
  PRIMARY KEY ("user_id", "art_id")
);

CREATE TABLE "tags" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE "arts_tags" (
  "art_id" INT NOT NULL,
  "tag_id" INT NOT NULL,
  PRIMARY KEY ("art_id", "tag_id")
);

CREATE TABLE "files" (
  "id" SERIAL PRIMARY KEY,
  "art_id" INT NOT NULL,
  "filename" VARCHAR NOT NULL,
  "filetype" VARCHAR NOT NULL,
  "url" VARCHAR NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE "codes" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR UNIQUE NOT NULL,
  "value" FLOAT NOT NULL,
  "exp_time" TIMESTAMP NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE "users_used_codes" (
  "user_id" INT NOT NULL,
  "code_id" INT NOT NULL,
  PRIMARY KEY ("user_id", "code_id")
);

ALTER TABLE "tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "oauths" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "arts" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id");
ALTER TABLE "users_starred_arts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "users_starred_arts" ADD FOREIGN KEY ("art_id") REFERENCES "arts" ("id");
ALTER TABLE "users_bought_arts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "users_bought_arts" ADD FOREIGN KEY ("art_id") REFERENCES "arts" ("id");
ALTER TABLE "users_created_arts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "users_created_arts" ADD FOREIGN KEY ("art_id") REFERENCES "arts" ("id");
ALTER TABLE "arts_tags" ADD FOREIGN KEY ("art_id") REFERENCES "arts" ("id");
ALTER TABLE "arts_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id");
ALTER TABLE "files" ADD FOREIGN KEY ("art_id") REFERENCES "arts" ("id");
ALTER TABLE "users_used_codes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "users_used_codes" ADD FOREIGN KEY ("code_id") REFERENCES "codes" ("id");

CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_tokens_table BEFORE UPDATE ON "tokens" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_oauths_table BEFORE UPDATE ON "oauths" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_arts_table BEFORE UPDATE ON "arts" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_files_table BEFORE UPDATE ON "files" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_codes_table BEFORE UPDATE ON "codes" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();

COMMIT;
