-- +goose Up
-- +goose StatementBegin
REVOKE CREATE ON schema public FROM public;
CREATE SCHEMA IF NOT EXISTS rewards_api;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA rewards_api CASCADE;
GRANT CREATE, USAGE ON schema public TO public;
-- +goose StatementEnd
