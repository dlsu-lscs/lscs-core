-- +goose Up
-- +goose StatementBegin

ALTER TABLE members ADD COLUMN image_url VARCHAR(512) DEFAULT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE members DROP COLUMN IF EXISTS image_url;

-- +goose StatementEnd
