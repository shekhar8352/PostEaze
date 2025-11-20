package businessv1

import (
	"context"
	"errors"

	"github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func containsAny(arr pq.StringArray, values ...string) bool {
	for _, item := range arr {
		for _, value := range values {
			if item == value {
				return true
			}
		}
	}
	return false
}

func GetUserById(ctx context.Context, userID string) (*entities.User, error) {
	user, err := repositories.GetUserByID(ctx, userID)
	if err != nil {
		utils.Logger.Error(ctx, "Error fetching user by ID: %v", err)
		return nil, err
	}
	return user, nil
}

func UpdateUser(ctx context.Context, body modelsv1.UpdateUserParams) (*entities.User, error) {
	user, err := repositories.GetUserByID(ctx, body.ID)
	if err != nil {
		utils.Logger.Error(ctx, "Error fetching user by ID: %v", err)
		return nil, err
	}

	if utils.ContainsAny(user.Platforms, "google", "email", "microsoft") {
		return nil, errors.New("user email cannot be updated")
	}

	err = repositories.UpdateUser(ctx, user.ID, user.Email)
	if err != nil {
		utils.Logger.Error(ctx, "Error updating user: %v", err)
		return nil, err
	}
	return user, nil
}
