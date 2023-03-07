-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TYPE referrals_status AS ENUM ('ReferralComplete', 'ReferralInvalid', 'TxFailed', 'Started');

CREATE TABLE referrals(
    chain_id bigint NOT NULL,
    job_status referrals_status NOT NULL,
    referred bytea NOT NULL,
    CONSTRAINT referred_address_check CHECK (length(referred) = 20),
    referrer bytea NOT NULL,
    CONSTRAINT referrer_address_check CHECK (length(referrer) = 20),
    transaction_hash bytea NOT NULL,
    CONSTRAINT referrals_transaction_hash_check CHECK (length(transaction_hash) = 32),
    log_index integer NOT NULL,
    block_timestamp timestamp with time zone NOT NULL,
    updated_at           timestamptz not null default current_timestamp,

    CONSTRAINT referrals_pkey PRIMARY KEY (chain_id, transaction_hash, log_index)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE referrals;
-- +goose StatementEnd