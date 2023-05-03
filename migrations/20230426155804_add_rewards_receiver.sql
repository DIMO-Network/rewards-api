-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards ADD COLUMN rewards_receiver_ethereum_address CHAR(42);
UPDATE rewards SET rewards_receiver_ethereum_address = user_ethereum_address;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards DROP COLUMN rewards_receiver_ethereum_address;
-- +goose StatementEnd
