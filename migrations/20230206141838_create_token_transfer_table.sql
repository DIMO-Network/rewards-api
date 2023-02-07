-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TABLE token_transfers(
    contract_address bytea NOT NULL
        CONSTRAINT token_transfers_contract_address_check CHECK (length(contract_address) = 20),
    user_address_from bytea NOT NULL
        CONSTRAINT token_transfers_user_address_to_check CHECK (length(user_address_from) = 20),
    user_address_to bytea NOT NULL
        CONSTRAINT token_transfers_user_address_from_check CHECK (length(user_address_to) = 20),
    amount numeric(78, 0) NOT NULL,
    tx_type varchar,
    created_at           timestamptz not null default current_timestamp,
    updated_at           timestamptz not null default current_timestamp,

    CONSTRAINT token_transfers_pkey PRIMARY KEY (contract_address, user_address_from, user_address_to, amount, created_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE token_transfers;
-- +goose StatementEnd
