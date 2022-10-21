-- +goose Up
-- +goose StatementBegin
-- +goose StatementEnd

SET search_path TO rewards_api, public;

CREATE TABLE meta_transaction_requests (
    id char(27) NOT NULL,
    "hash" BYTEA,
    "status" TEXT,
    successful BOOLEAN,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp
);


-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE meta_transaction_requests;
-- +goose StatementEnd
