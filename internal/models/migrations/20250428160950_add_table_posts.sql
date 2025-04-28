-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS POSTS (
    ID INTEGER PRIMARY KEY,
    CONTENT TEXT,
    CREATED_AT TEXT NOT NULL DEFAULT (DATETIME('subsec'))
) STRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE POSTS;
-- +goose StatementEnd
