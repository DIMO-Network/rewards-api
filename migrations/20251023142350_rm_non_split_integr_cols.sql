-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards DROP COLUMN integration_points;
ALTER TABLE rewards DROP COLUMN tokens;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

-- No real going back here.
-- +goose StatementEnd
