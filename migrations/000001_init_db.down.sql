DROP TRIGGER IF EXISTS update_timestamp_users;
DROP TRIGGER IF EXISTS update_timestamp_tokens;
DROP TRIGGER IF EXISTS update_timestamp_oauths;
DROP TRIGGER IF EXISTS update_timestamp_arts;
DROP TRIGGER IF EXISTS update_timestamp_files;

DROP TABLE IF EXISTS "users"; -- CASCADE;
DROP TABLE IF EXISTS "tokens"; -- CASCADE;
DROP TABLE IF EXISTS "oauths"; -- CASCADE;
DROP TABLE IF EXISTS "follow"; -- CASCADE;
DROP TABLE IF EXISTS "arts"; -- CASCADE;
DROP TABLE IF EXISTS "downloaded_arts"; -- CASCADE;
DROP TABLE IF EXISTS "users_starred_arts"; -- CASCADE;
DROP TABLE IF EXISTS "users_bought_arts"; -- CASCADE;
DROP TABLE IF EXISTS "tags"; -- CASCADE;
DROP TABLE IF EXISTS "arts_tags"; -- CASCADE;
DROP TABLE IF EXISTS "files"; -- CASCADE;
DROP TABLE IF EXISTS "codes"; -- CASCADE;
DROP TABLE IF EXISTS "users_used_codes"; -- CASCADE;
