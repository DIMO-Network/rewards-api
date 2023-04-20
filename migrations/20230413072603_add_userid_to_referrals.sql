-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals ADD COLUMN referrer_user_id text NOT NULL DEFAULT '';
ALTER TABLE referrals ALTER COLUMN referrer_user_id DROP DEFAULT;


ALTER TABLE referrals ADD COLUMN referree_user_id text NOT NULL DEFAULT '';
ALTER TABLE referrals ALTER COLUMN referree_user_id DROP DEFAULT;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals 
  DROP COLUMN referrer_user_id;

ALTER TABLE referrals 
  DROP COLUMN referree_user_id;
-- +goose StatementEnd