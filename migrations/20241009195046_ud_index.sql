-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE INDEX rewards_user_device_id_idx ON rewards (user_device_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP INDEX rewards_user_device_id_idx;
-- +goose StatementEnd
