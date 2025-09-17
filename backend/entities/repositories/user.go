package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/shekhar8352/PostEaze/entities"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func CreateUser(ctx context.Context, tx database.Database, user entities.User) (*entities.User, error) {
	err := tx.QueryRaw(ctx, &user, entities.CreateUser)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func InsertRefreshTokenOfUser(ctx context.Context, userID, refreshToken string, expireAt time.Time) error {
	data := entities.User{
		ID:           userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expireAt,
	}
	return database.Get().QueryRaw(ctx, &data, entities.InsertRefreshToken)
}

func GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	data := entities.User{Email: email}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByEmail)
	if err != nil {
		if errors.Is(err, database.ErrNoRecords) {
			return nil, database.ErrNoRecords
		}
		return nil, err
	}
	return &data, nil
}

func GetUserbyToken(ctx context.Context, token string) (*entities.User, error) {
	data := entities.User{RefreshToken: token}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByToken)
	if err != nil {
		if errors.Is(err, database.ErrNoRecords) {
			return nil, database.ErrNoRecords
		}
		return nil, err
	}

	// fetch full user details by ID
	user, err := GetUserByID(ctx, data.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func RevokeTokenForUser(ctx context.Context, userID string) error {
	data := entities.User{ID: userID}
	return database.Get().QueryRaw(ctx, &data, entities.RevokeTokens)
}

func CreateUserWithFirebase(ctx context.Context, tx database.Database, user modelsv1.User) (*entities.User, error) {
	data := entities.User{
		ID:         user.ID,
		FirebaseID: user.FirebaseID,
		Name:       user.Name,
		Email:      user.Email,
		Platforms:  user.Platforms,
	}
	err := tx.QueryRaw(ctx, &data, entities.CreateUserWithFirebase)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	data := entities.User{ID: userID}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByID)
	if err != nil {
		if errors.Is(err, database.ErrNoRecords) {
			return nil, database.ErrNoRecords
		}
		return nil, err
	}
	return &data, nil
}

func GetUserByFirebaseID(ctx context.Context, firebaseID string) (*entities.User, error) {
	data := entities.User{FirebaseID: firebaseID}
	err := database.Get().QueryRaw(ctx, &data, entities.GetUserByFirebaseID)
	if err != nil {
		if errors.Is(err, database.ErrNoRecords) {
			return nil, database.ErrNoRecords
		}
		return nil, err
	}
	return &data, nil
}

func UpdateUserPlatforms(ctx context.Context, userID string, platforms []string) error {
	data := entities.User{
		ID:        userID,
		Platforms: platforms,
	}
	return database.Get().QueryRaw(ctx, &data, entities.UpdateUserPlatforms)
}
