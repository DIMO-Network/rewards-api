-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards ADD COLUMN override boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards DROP COLUMN override;
-- +goose StatementEnd
