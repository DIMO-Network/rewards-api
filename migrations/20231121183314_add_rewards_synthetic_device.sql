-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards ADD COLUMN synthetic_device_id int;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards DROP COLUMN synthetic_device_id;
-- +goose StatementEnd
