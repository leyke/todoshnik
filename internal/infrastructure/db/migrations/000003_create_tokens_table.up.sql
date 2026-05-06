CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INT,
    hash  VARCHAR(255),
    device  VARCHAR(255),
    expires_at  BIGINT
)