-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards ALTER COLUMN user_device_token_id TYPE int USING user_device_token_id::int;
ALTER TABLE rewards ALTER COLUMN user_device_token_id SET NOT NULL;

ALTER TABLE rewards DROP CONSTRAINT rewards_issuance_week_id_user_device_id_pkey;
ALTER TABLE rewards ADD PRIMARY KEY (issuance_week_id, user_device_token_id);

ALTER TABLE rewards DROP COLUMN user_device_id;
ALTER TABLE rewards DROP COLUMN user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

-- Really no going back here.
-- +goose StatementEnd
