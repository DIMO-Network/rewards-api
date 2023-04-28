-- +goose Up
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals RENAME COLUMN referree_user_id TO referee_user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path = rewards_api, public;

ALTER TABLE referrals RENAME COLUMN referee_user_id TO referree_user_id;
-- +goose StatementEnd
