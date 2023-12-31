CREATE TABLE lists (
    list_id BIGSERIAL PRIMARY KEY,
    list_title VARCHAR NOT NULL,
    user_id BIGINT REFERENCES users ON DELETE CASCADE,
    UNIQUE(list_title, user_id)
);