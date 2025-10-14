-- +goose Up
-- +goose StatementBegin
SET search_path TO rewards_api, public;

ALTER TABLE referrals DROP COLUMN referrer_user_id;
ALTER TABLE referrals DROP COLUMN referee_user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO rewards_api, public;
-- Can't easily go back on this one, unless we could construct user ids from addresses in SQL.
-- +goose StatementEnd
