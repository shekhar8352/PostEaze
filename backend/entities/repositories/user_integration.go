package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func UpsertUserIntegration(ctx context.Context, ui *entities.UserIntegration) error {
	db := database.GetDB()
	q := `
		INSERT INTO user_integrations (
			user_id, provider, provider_account_email, access_token, refresh_token,
			scopes, expires_at, revoked
		) VALUES ($1, $2, $3, $4, $5, $6, $7, FALSE)
		ON CONFLICT (user_id, provider) DO UPDATE SET
			provider_account_email = EXCLUDED.provider_account_email,
			access_token = EXCLUDED.access_token,
			refresh_token = COALESCE(EXCLUDED.refresh_token, user_integrations.refresh_token),
			scopes = EXCLUDED.scopes,
			expires_at = EXCLUDED.expires_at,
			revoked = FALSE,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, q,
		ui.UserID, ui.Provider, ui.ProviderAccountEmail, ui.AccessToken, ui.RefreshToken,
		ui.Scopes, ui.ExpiresAt,
	).Scan(&ui.ID, &ui.CreatedAt, &ui.UpdatedAt)
}

func GetUserIntegration(ctx context.Context, userID uuid.UUID, provider string) (*entities.UserIntegration, error) {
	db := database.GetDB()
	q := `
		SELECT id, user_id, provider, provider_account_email, access_token, refresh_token,
		       scopes, expires_at, revoked, created_at, updated_at
		FROM user_integrations
		WHERE user_id = $1 AND provider = $2 AND revoked = FALSE
	`
	var ui entities.UserIntegration
	err := db.QueryRowContext(ctx, q, userID, provider).Scan(
		&ui.ID, &ui.UserID, &ui.Provider, &ui.ProviderAccountEmail, &ui.AccessToken, &ui.RefreshToken,
		&ui.Scopes, &ui.ExpiresAt, &ui.Revoked, &ui.CreatedAt, &ui.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ui, nil
}

func UpdateUserIntegrationTokens(ctx context.Context, id int64, accessToken []byte, expiresAt *time.Time) error {
	db := database.GetDB()
	q := `
		UPDATE user_integrations
		SET access_token = $2, expires_at = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := db.ExecContext(ctx, q, id, accessToken, expiresAt)
	return err
}

func RevokeUserIntegration(ctx context.Context, userID uuid.UUID, provider string) error {
	db := database.GetDB()
	q := `
		UPDATE user_integrations
		SET revoked = TRUE, updated_at = NOW()
		WHERE user_id = $1 AND provider = $2
	`
	res, err := db.ExecContext(ctx, q, userID, provider)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
