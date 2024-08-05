-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id bigserial primary key,
  username varchar(255) not null, 
  password varchar(255) not null, 
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  deleted_at timestamp default current_timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
