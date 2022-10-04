-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path TO rewards_api, public;

CREATE TABLE IF NOT EXISTS token_allocation (
    issuance_week_id int NOT NULL,
    user_device_id char(27) NOT NULL,
    tokens numeric(78,0),
    week_start timestamptz NOT NULL,
    week_end timestamptz NOT NULL
);

ALTER TABLE token_allocation ADD CONSTRAINT tokens_id_pkey PRIMARY KEY (issuance_week_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

SET search_path TO rewards_api, public;

DROP TABLE weekly_point_total;

-- +goose StatementEnd
