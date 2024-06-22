-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE url_data
RENAME COLUMN expires_at TO visited_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE url_data
    RENAME COLUMN VISITED_AT TO expires_at;
-- +goose StatementEnd
