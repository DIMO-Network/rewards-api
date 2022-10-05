-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TYPE issuance_weeks_job_status AS ENUM ('Started', 'Failed', 'PointsAllocated', 'BeginTokenDistribution', 'Finished');

CREATE TABLE issuance_weeks (
    id int, -- This number is the number of complete weeks that have passed since the
            -- beginning of issuance. Right now the beginning is 2022-01-31 05:00 UTC.
            -- We start counting at 0.
    job_status issuance_weeks_job_status NOT NULL,

    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    points_distributed bigint,
    weekly_token_allocation numeric(28, 0)
);

ALTER TABLE issuance_weeks ADD CONSTRAINT issuance_weeks_id_pkey PRIMARY KEY (id);

CREATE TABLE rewards (
    issuance_week_id int NOT NULL,
    user_device_id char(27) NOT NULL,

    user_id text NOT NULL,

    connection_streak int NOT NULL,
    disconnection_streak int NOT NULL,
    streak_points int NOT NULL,

    integration_ids text[] NOT NULL DEFAULT '{}',
    integration_points int NOT NULL,

    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp
);

ALTER TABLE rewards ADD CONSTRAINT rewards_issuance_week_id_fkey FOREIGN KEY (issuance_week_id) REFERENCES issuance_weeks(id);
ALTER TABLE rewards ADD CONSTRAINT rewards_issuance_week_id_user_device_id_pkey PRIMARY KEY (issuance_week_id, user_device_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE device_week_rewards;
DROP TABLE issuance_week;
DROP TYPE issuance_week_job_status;
-- +goose StatementEnd
