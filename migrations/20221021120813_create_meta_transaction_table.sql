-- +goose Up
-- +goose StatementBegin
-- +goose StatementEnd

SET search_path TO rewards_api, public;

CREATE TABLE meta_transaction_requests (
    id char(27) PRIMARY KEY,
    "hash" BYTEA,
    "status" TEXT,
    successful BOOLEAN,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp
);

ALTER TABLE rewards ADD CONSTRAINT meta_transaction_id_fkey FOREIGN KEY(transfer_meta_transaction_request_id) REFERENCES meta_transaction_requests(id);

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE meta_transaction_requests;
-- +goose StatementEnd
