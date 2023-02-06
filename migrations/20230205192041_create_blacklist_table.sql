-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TABLE blacklist (
    user_ethereum_address CHAR(42),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    note TEXT NOT NULL,

    CONSTRAINT blacklist_user_ethereum_address_pkey PRIMARY KEY (user_ethereum_address)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE blacklist;
-- +goose StatementEnd
