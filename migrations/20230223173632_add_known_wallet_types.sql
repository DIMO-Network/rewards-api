-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

CREATE TYPE wallet_type AS ENUM ('Baseline', 'Referrals', 'Marketplace');

ALTER TABLE known_wallets ADD COLUMN "type" wallet_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE known_wallets DROP COLUMN "type";

DROP TYPE wallet_type;
-- +goose StatementEnd
