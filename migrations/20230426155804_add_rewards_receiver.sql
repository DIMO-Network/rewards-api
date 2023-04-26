-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards ADD COLUMN rewards_receiver_ethereum_address CHAR(42);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards DROP COLUMN rewards_receiver_ethereum_address;
-- +goose StatementEnd
