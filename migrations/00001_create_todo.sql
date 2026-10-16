-- +goose Up
CREATE TABLE todos(
    todo_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    todo_name VARCHAR(20) NOT NULL UNIQUE,
    todo_description VARCHAR(255) NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS todos;
