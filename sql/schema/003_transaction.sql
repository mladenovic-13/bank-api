-- +goose Up
CREATE TYPE transaction_type AS ENUM ('TRANSFER', 'PAYMENT', 'WITHDRAWAL', 'DEPOSIT');


CREATE TABLE transactions(
  id UUID PRIMARY KEY,
  sender_number UUID NOT NULL REFERENCES accounts(number),
  receiver_number UUID NOT NULL REFERENCES accounts(number),
  amount DECIMAL(10, 2) NOT NULL,
  currency currency NOT NULL,
  transaction_type transaction_type NOT NULL,
  created_at TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE transactions;
DROP TYPE transaction_type;