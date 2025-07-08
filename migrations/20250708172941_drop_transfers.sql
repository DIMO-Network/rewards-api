-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE known_wallets;
DROP TYPE wallet_type;
DROP table token_transfers;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

-- There is no going back.
-- +goose StatementEnd
