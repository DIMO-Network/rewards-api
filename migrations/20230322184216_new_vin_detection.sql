-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TABLE vins (
    vin char(17) CONSTRAINT vins_pkey PRIMARY KEY,
    first_earning_week int NOT NULL,
    first_earning_token_id numeric(78) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

DROP TABLE vins;
-- +goose StatementEnd
