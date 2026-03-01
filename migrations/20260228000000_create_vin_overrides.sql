-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TABLE vin_overrides (
    token_id int NOT NULL,
    vin char(17) NOT NULL,
    note text,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    CONSTRAINT vin_overrides_pkey PRIMARY KEY (token_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE vin_overrides;
-- +goose StatementEnd
