-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api;
CREATE TABLE device_week_rewards (
    user_device_id char(27) NOT NULL,
    user_id text NOT NULL,
    week int NOT NULL,
    miles_driven double NOT NULL,
    connected BOOLEAN NOT NULL,
    weeks_connected_streak int NOT NULL,
    weeks_disconnected_streak int NOT NULL,
    streak_points int NOT NULL,
    connection_method_points int NOT NULL,
    engine_type_points int NOT NULL,
    rarity_points int NOT NULL,
    created_at timestamptz NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api;
DROP TABLE device_week_rewards;
-- +goose StatementEnd
