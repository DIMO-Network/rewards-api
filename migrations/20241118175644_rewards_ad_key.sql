-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

CREATE INDEX rewards_aftermarket_token_id_idx ON rewards (aftermarket_token_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

DROP INDEX IF EXISTS rewards_aftermarket_token_id_idx;
-- +goose StatementEnd
