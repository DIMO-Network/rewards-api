-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TABLE token_transfers(
    address_from bytea NOT NULL,
    address_to bytea NOT NULL,
    amount numeric(78) NOT NULL,
    transaction_hash bytea NOT NULL,
    log_index integer NOT NULL,
    block_timestamp timestamp with time zone NOT NULL,
    updated_at           timestamptz not null default current_timestamp,

    CONSTRAINT token_transfers_pkey PRIMARY KEY (transaction_hash, log_index)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE token_transfers;
-- +goose StatementEnd