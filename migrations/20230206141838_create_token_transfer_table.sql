-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TABLE token_transfers(
    chain_id bigint NOT NULL,
    address_from bytea NOT NULL
    CONSTRAINT token_transfers_address_from_check CHECK (length(address_from) = 20),
    address_to bytea NOT NULL
    CONSTRAINT token_transfers_address_to_check CHECK (length(address_to) = 20),
    amount numeric(78) NOT NULL,
    transaction_hash bytea NOT NULL
    CONSTRAINT token_transfers_transaction_hash_check CHECK (length(transaction_hash) = 32),
    log_index integer NOT NULL,
    block_timestamp timestamp with time zone NOT NULL,
    updated_at           timestamptz not null default current_timestamp,

    CONSTRAINT token_transfers_pkey PRIMARY KEY (chain_id, transaction_hash, log_index)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE token_transfers;
-- +goose StatementEnd