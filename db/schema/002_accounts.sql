-- +goose Up
CREATE TYPE currency AS ENUM ('EUR', 'USD', 'RSD');

CREATE TABLE accounts(
  id UUID PRIMARY KEY,

  name TEXT NOT NULL,
  number UUID NOT NULL,
  balance DECIMAL DEFAULT 0,
  currency currency,

  user_id UUID NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

ALTER TABLE accounts 
ADD CONSTRAINT fk_accounts_users 
FOREIGN KEY (user_id)
 REFERENCES users(id);

-- +goose Down
DROP TABLE accounts;
DROP TYPE currency;