-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TABLE attestations (
    id BYTEA,
    transaction_id char(27) NOT NULL CONSTRAINT attestations_transaction_id_fkey REFERENCES meta_transaction_requests (id),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    "root" BYTEA NOT NULL,
    -- do we want to store chain id?
    CONSTRAINT attestations_pkey PRIMARY KEY (transaction_id, "root")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE attestations;
-- +goose StatementEnd
