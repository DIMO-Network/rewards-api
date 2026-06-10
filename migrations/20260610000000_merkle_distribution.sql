-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TABLE merkle_roots (
    issuance_week_id int PRIMARY KEY REFERENCES issuance_weeks(id),
    pool_id int NOT NULL,
    root bytea NOT NULL,
    total_allocation numeric(38) NOT NULL,
    proofs_uri text NOT NULL,
    meta_transaction_request_id text REFERENCES meta_transaction_requests(id),
    set_successful boolean NOT NULL DEFAULT false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE merkle_roots;
-- +goose StatementEnd
