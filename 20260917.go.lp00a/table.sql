CREATE TABLE IF NOT EXISTS todos (
  id          INTEGER   PRIMARY KEY AUTOINCREMENT,
  task        TEXT      NOT NULL,
  is_done     BOOLEAN   NOT NULL,
  created_at  TIMESTAMP NOT NULL,
  updated_at  TIMESTAMP NOT NULL
);
