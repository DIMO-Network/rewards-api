-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TABLE vins (
    vin char(17) CONSTRAINT vins_pkey PRIMARY KEY,
    first_week_earning int NOT NULL
);

ALTER TABLE rewards
    ADD COLUMN new_vin boolean NOT NULL DEFAULT FALSE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE rewards
    DROP COLUMN new_vin;

DROP TABLE vins;
-- +goose StatementEnd
