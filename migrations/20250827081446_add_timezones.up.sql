CREATE TABLE time_zone_users (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,          -- например "Europe/Moscow"
    UNIQUE(user_id)
);