-- +goose Up
CREATE TABLE users(
  id UUID PRIMARY KEY,

  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  isAdmin BOOLEAN DEFAULT false  NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE  users;