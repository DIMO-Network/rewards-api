-- +goose Up
-- +goose StatementBegin
ALTER TABLE issuance_weeks ADD COLUMN tokens_per_1k_points numeric(78);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE issuance_weeks DROP COLUMN tokens_per_1k_points;
-- +goose StatementEnd
