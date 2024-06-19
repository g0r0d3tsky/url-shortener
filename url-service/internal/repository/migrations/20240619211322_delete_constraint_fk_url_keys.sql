-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE url_keys
    DROP CONSTRAINT url_keys_url_id_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE url_keys
    ADD CONSTRAINT url_keys_url_id_fkey FOREIGN KEY (url_id) REFERENCES url_data (id);
-- +goose StatementEnd
