-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TYPE issuance_week_job_status AS ENUM ('Started', 'Failed', 'Finished');

CREATE TABLE issuance_week (
    id int PRIMARY KEY, -- This number is the number of complete weeks that have passed since the
                        -- beginning of issuance. Right now this is 0900 UTC on March 14, 2022. We
                        -- start counting at 0.
    job_status issuance_week_job_status NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp
);

CREATE TABLE device_week_rewards (
    user_device_id char(27) NOT NULL,
    issuance_week_id int NOT NULL,
    device_definition_id char(27) NOT NULL,
    vin char(17) NOT NULL,
    user_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,

    miles_driven double precision NOT NULL,
    connected BOOLEAN NOT NULL,
    weeks_connected_streak int NOT NULL,
    weeks_disconnected_streak int NOT NULL,
    streak_points int NOT NULL,
    connection_method_points int NOT NULL,

    engine_type_points int NOT NULL,
    rarity_points int NOT NULL,

    PRIMARY KEY(issuance_week_id, user_device_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;
DROP TABLE device_week_rewards;
-- +goose StatementEnd
