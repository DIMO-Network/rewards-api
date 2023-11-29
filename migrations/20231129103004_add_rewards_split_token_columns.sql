-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    ADD COLUMN aftermarket_device_points int,
    ADD COLUMN synthetic_device_points int,
    ADD COLUMN aftermarket_device_tokens numeric(26),
    ADD COLUMN synthetic_device_tokens numeric(26),
    ADD COLUMN streak_tokens numeric(26);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    DROP COLUMN aftermarket_device_points,
    DROP COLUMN synthetic_device_points,
    DROP COLUMN aftermarket_device_tokens,
    DROP COLUMN synthetic_device_tokens,
    DROP COLUMN streak_tokens;
-- +goose StatementEnd
