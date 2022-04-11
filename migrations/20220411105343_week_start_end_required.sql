-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE issuance_weeks ALTER COLUMN starts_at SET NOT NULL;
ALTER TABLE issuance_weeks ALTER COLUMN ends_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE issuance_weeks ALTER COLUMN starts_at DROP NOT NULL;
ALTER TABLE issuance_weeks ALTER COLUMN ends_at DROP NOT NULL;
-- +goose StatementEnd
