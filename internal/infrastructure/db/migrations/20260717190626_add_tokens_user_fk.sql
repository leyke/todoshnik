-- +goose Up

ALTER TABLE tokens
ADD CONSTRAINT fk_tokens_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;


-- +goose Down

ALTER TABLE tokens
DROP CONSTRAINT fk_tokens_user;