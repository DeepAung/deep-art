BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "users";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_tokens_table ON "tokens";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauths_table ON "oauths";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_arts_table ON "arts";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_files_table ON "files";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_codes_table ON "codes";

DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "tokens" CASCADE;
DROP TABLE IF EXISTS "oauths" CASCADE;
DROP TABLE IF EXISTS "arts" CASCADE;
DROP TABLE IF EXISTS "users_starred_arts" CASCADE;
DROP TABLE IF EXISTS "users_bought_arts" CASCADE;
DROP TABLE IF EXISTS "users_created_arts" CASCADE;
DROP TABLE IF EXISTS "tags" CASCADE;
DROP TABLE IF EXISTS "arts_tags" CASCADE;
DROP TABLE IF EXISTS "files" CASCADE;
DROP TABLE IF EXISTS "codes" CASCADE;
DROP TABLE IF EXISTS "users_used_codes" CASCADE;

DROP TYPE IF EXISTS "social";

COMMIT;
