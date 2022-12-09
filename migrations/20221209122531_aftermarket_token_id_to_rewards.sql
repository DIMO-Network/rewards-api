-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards
    ADD COLUMN aftermarket_token_id  NUMERIC(25);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE rewards 
    DROP COLUMN aftermarket_token_id;

-- +goose StatementEnd
