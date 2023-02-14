-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE token_transfers
    ADD chain_id VARCHAR;

UPDATE token_transfers 
-- should we set the default value to polygon since we've only been polling that so far?
SET chain_id = '---'
WHERE chain_id IS NULL;

ALTER TABLE token_transfers
    ALTER COLUMN chain_id SET NOT NULL;

ALTER TABLE token_transfers
    DROP CONSTRAINT token_transfers_pkey;

ALTER TABLE token_transfers
    ADD PRIMARY KEY (transaction_hash, log_index, chain_id);

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE token_transfers
    DROP COLUMN chain_id;



-- +goose StatementEnd
