-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals
ADD COLUMN issuance_week_id int NOT NULL DEFAULT 1;

ALTER TABLE referrals ALTER COLUMN issuance_week_id DROP DEFAULT;

ALTER TABLE referrals ADD CONSTRAINT referrals_issuance_week_id_fkey FOREIGN KEY (issuance_week_id) REFERENCES issuance_weeks(id);




-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals 
  DROP CONSTRAINT IF EXISTS referrals_issuance_week_id_fkey;

ALTER TABLE referrals 
  DROP COLUMN issuance_week_id;
-- +goose StatementEnd