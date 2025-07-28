package repositories

import (
	"context"
	"time"

	"github.com/shekhar8352/PostEaze/entities"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func CreateUser(ctx context.Context, tx database.Database, user modelsv1.User) (*entities.User, error) {
	data := entities.User{
		Name:     user.Name,
		Email:    user.Email,
		UserType: string(user.UserType),
		Password: user.Password,
	}
	err := tx.QueryRaw(ctx, &data, entities.CreateUser)
	return &data, err
}

func InsertRefreshTokenOfUser(ctx context.Context, userID, refreshToken string, expireAt time.Time) error {

	data := entities.User{
		ID:           userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expireAt,
	}
	err := database.Get().QueryRaw(ctx, &data, entities.InsertRefreshToken)
	return err
}

func GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	data := entities.User{
		Email: email,
	}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByEmail)
	if err != nil {
		return nil, err
	}
	return &data, nil

}

func GetUserbyToken(ctx context.Context, token string) (*entities.User, error) {
	data := entities.User{
		RefreshToken: token,
	}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByToken)
	if err != nil {
		return nil, err
	}

	// now fetch user details through user ID

	err = database.Get().QueryRaw(ctx, &data, entities.GetUserByToken)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func RevokeTokenForUser(ctx context.Context, userID string) error {
	data := entities.User{
		ID: userID,
	}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByToken)
	return err
}
