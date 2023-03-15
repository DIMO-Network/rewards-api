-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TYPE referrals_transfer_failure_reason AS ENUM (
    'ReferralInvalid',
    'TxReverted'
);

CREATE TABLE referrals (
    referee bytea NOT NULL CONSTRAINT referrals_referee_check CHECK (length(referee) = 20),
    referrer bytea NOT NULL CONSTRAINT referrals_referrer_check CHECK (length(referrer) = 20),
    transfer_successful boolean,
    transfer_failure_reason referrals_transfer_failure_reason,
    request_id char(27) NOT NULL CONSTRAINT referrals_request_id_fkey REFERENCES meta_transaction_requests (id),
    CONSTRAINT referrals_pkey PRIMARY KEY (referee, referrer)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE referrals;
<<<<<<< HEAD
=======

DROP TYPE referrals_transfer_failure_reason;
>>>>>>> main
-- +goose StatementEnd