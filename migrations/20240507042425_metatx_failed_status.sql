-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TYPE meta_transaction_request_status ADD VALUE 'Failed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;

-- No great way to do this.
-- +goose StatementEnd
