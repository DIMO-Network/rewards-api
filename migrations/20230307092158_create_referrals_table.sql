-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TYPE referrals_job_status AS ENUM ('Complete', 'Started');

CREATE TYPE referrals_transfer_failure_reason AS ENUM ('ReferralInvalid', 'TxReverted');

CREATE TABLE referrals(
    id char(27) NOT NULL,
    job_status referrals_job_status NOT NULL,
    
    referred bytea NOT NULL,
    CONSTRAINT referred_address_check CHECK (length(referred) = 20),
    referrer bytea NOT NULL,
    CONSTRAINT referrer_address_check CHECK (length(referrer) = 20),
    
    transfer_successful BOOLEAN,
    transfer_failure_reason referrals_transfer_failure_reason,
    
    CONSTRAINT referrals_pkey PRIMARY KEY (referred, referrer)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE referrals;
-- +goose StatementEnd