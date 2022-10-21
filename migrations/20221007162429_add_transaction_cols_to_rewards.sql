-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    ADD COLUMN tokens numeric(26),
    ADD COLUMN user_ethereum_address CHAR(42),
    ADD COLUMN user_device_token_id numeric(78, 0),
    ADD COLUMN transfer_meta_transaction_request_id VARCHAR,
    ADD COLUMN transfer_successful BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    DROP COLUMN tokens;
    DROP user_ethereum_address
    DROP user_device_token_id
    DROP transfer_meta_transaction_request_id
    DROP transfer_successful;
-- +goose StatementEnd
