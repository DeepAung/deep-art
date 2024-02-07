BEGIN;

INSERT INTO "users" (
  "username",
  "email",
  "password",
)
VALUES
  ('deepaung', 'i.deepaung@gmail.com', '$2a$10$bMm02k51y5q.MDIYbB13POV5cTcUwZONd77K2RXg/PdOQiybwLDca'),
  ('test', 'test@gmail.com', '$2a$10$xiZr4jo2Rq.GWgQPeEpUseivlnxIHoM5LYECEUKuvPONaxh30s8Ru'),
  ('example', 'example@gmail.com', '$2a$10$bMm02k51y5q.MDIYbB13POV5cTcUwZONd77K2RXg/PdOQiybwLDca');

COMMIT;
