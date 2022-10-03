-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path TO rewards_api, public;

CREATE TABLE IF NOT EXISTS weekly_point_total (
    id int, -- This number is the number of complete weeks that have passed since the
            -- beginning of issuance. Right now the beginning is 2022-01-31 05:00 UTC.
            -- We start counting at 0.
    points numeric(78,0),
    week_start timestamptz NOT NULL,
    week_end timestamptz NOT NULL
);

ALTER TABLE weekly_point_total ADD CONSTRAINT weekly_points_id_pkey PRIMARY KEY (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

SET search_path TO rewards_api, public;

DROP TABLE weekly_point_total;

-- +goose StatementEnd
