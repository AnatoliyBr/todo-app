CREATE TABLE tasks (
    task_id BIGSERIAL PRIMARY KEY,
    task_title VARCHAR NOT NULL,
    details VARCHAR NOT NULL,
    deadline TIMESTAMP NOT NULL,
    done BOOLEAN DEFAULT FALSE,
    list_id BIGINT REFERENCES lists ON DELETE CASCADE,
    UNIQUE(task_title, list_id)
);