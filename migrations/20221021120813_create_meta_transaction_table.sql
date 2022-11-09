-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TYPE meta_transaction_request_status AS ENUM ('Unsubmitted', 'Submitted', 'Mined', 'Confirmed');

CREATE TABLE meta_transaction_requests (
    id char(27) PRIMARY KEY,
    "hash" TEXT,
    "status" meta_transaction_request_status NOT NULL,
    successful BOOLEAN,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp
);

ALTER TABLE rewards ADD CONSTRAINT rewards_transfer_meta_transaction_request_id_fkey FOREIGN KEY (transfer_meta_transaction_request_id) REFERENCES meta_transaction_requests(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE meta_transaction_requests CASCADE;

DROP TYPE meta_transaction_request_status;
-- +goose StatementEnd
