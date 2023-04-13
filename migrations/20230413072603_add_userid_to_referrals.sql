-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals
ADD COLUMN referrer_user_id text;

ALTER TABLE referrals
ADD COLUMN referree_user_id text;



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals 
  DROP COLUMN referrer_user_id;

ALTER TABLE referrals 
  DROP COLUMN referree_user_id;
-- +goose StatementEnd