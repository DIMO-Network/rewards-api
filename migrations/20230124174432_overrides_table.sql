-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE TABLE overrides (
    issuance_week_id int,
    user_device_id char(27),
    connection_streak int NOT NULL,
    note text NOT NULL,

    CONSTRAINT overrides_issuance_week_id_user_device_id_pkey PRIMARY KEY (issuance_week_id, user_device_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP TABLE overrides;
-- +goose StatementEnd
