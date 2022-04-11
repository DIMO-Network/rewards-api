-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE issuance_weeks ADD COLUMN starts_at timestamptz;
ALTER TABLE issuance_weeks ADD COLUMN ends_at timestamptz;
UPDATE issuance_weeks SET starts_at = timestamptz '2022-01-31 05:00:00+00' + interval '7 day' * id;
UPDATE issuance_weeks SET ends_at = timestamptz '2022-01-31 05:00:00+00' + interval '7 day' * (id + 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE issuance_weeks DROP COLUMN starts_at;
ALTER TABLE issuance_weeks DROP COLUMN ends_at;
-- +goose StatementEnd
