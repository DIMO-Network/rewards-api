-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TYPE rewards_transfer_failure_reason AS ENUM ('DidntQualify', 'TxReverted');

ALTER TABLE rewards
    -- Max weekly is 1.3e6 * 1e18, so 25 digits or numeric(25) would be enough.
    ADD COLUMN user_ethereum_address CHAR(42),
    ADD COLUMN user_device_token_id NUMERIC(78),
    ADD COLUMN transfer_meta_transaction_request_id CHAR(27),
    ADD COLUMN transfer_successful BOOLEAN,
    ADD COLUMN transfer_failure_reason rewards_transfer_failure_reason;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    DROP COLUMN tokens,
    DROP COLUMN user_ethereum_address,
    DROP COLUMN user_device_token_id,
    DROP COLUMN transfer_meta_transaction_request_id,
    DROP COLUMN transfer_successful,
    DROP COLUMN transfer_failure_reason;

DROP TYPE rewards_transfer_failure_reason;
-- +goose StatementEnd
