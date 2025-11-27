-- +goose Down
-- +goose StatementBegin

/******************************************************************************************
 * DOWN MIGRATION
 * Drops channel_tokens first (due to FK dependencies), then channels.
 ******************************************************************************************/

-- Drop token indexes
DROP INDEX IF EXISTS idx_channel_tokens_expires_at;
DROP INDEX IF EXISTS idx_channel_tokens_channel_id;

-- Drop tokens table
DROP TABLE IF EXISTS channel_tokens;


-- Drop channel indexes
DROP INDEX IF EXISTS idx_channels_provider;
DROP INDEX IF EXISTS idx_channels_team_id;
DROP INDEX IF EXISTS idx_channels_owner_user_id;

-- Drop channels table
DROP TABLE IF EXISTS channels;

-- +goose StatementEnd
