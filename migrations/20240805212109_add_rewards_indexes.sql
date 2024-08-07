-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE INDEX rewards_user_ethereum_address_idx ON rewards (user_ethereum_address);
CREATE INDEX rewards_user_device_token_id_idx ON rewards (user_device_token_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP INDEX rewards_user_device_token_id_idx;
DROP INDEX rewards_user_ethereum_address_idx;
-- +goose StatementEnd
