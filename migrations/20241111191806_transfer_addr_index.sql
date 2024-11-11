-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE INDEX token_transfers_address_from_idx ON token_transfers (address_from);
CREATE INDEX token_transfers_address_to_idx ON token_transfers (address_from);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS token_transfers_address_from_idx;
DROP INDEX IF EXISTS token_transfers_address_to_idx;
-- +goose StatementEnd
