-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TABLE known_wallets (
    chain_id bigint,
    address bytea NOT NULL CONSTRAINT known_wallets_address_from_check CHECK (length(address) = 20),
    CONSTRAINT known_wallets_pkey PRIMARY KEY (chain_id, address)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE known_wallets;
-- +goose StatementEnd
