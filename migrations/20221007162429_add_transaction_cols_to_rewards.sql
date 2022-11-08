-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    ADD COLUMN tokens NUMERIC(26),
    ADD COLUMN user_ethereum_address CHAR(42),
    ADD COLUMN user_device_token_id NUMERIC(78, 0),
    ADD COLUMN transfer_meta_transaction_request_id char(27),
    ADD COLUMN transfer_successful BOOLEAN,
    ADD COLUMN transfer_fail_reason CHAR(20);

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
