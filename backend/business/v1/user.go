package businessv1

import (
	"context"
	"errors"

	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func Signup(ctx context.Context, params modelsv1.SignupParams) (map[string]interface{}, error) {
	// db := config.GetAppContext().DB

	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &modelsv1.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: hashedPassword,
		UserType: params.UserType,
	}
	// start transaction
	tx, err := database.GetTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if params.UserType == modelsv1.UserTypeTeam {
		// create user
		userCreated, err := repositories.CreateUser(ctx, tx, *user)
		if err != nil {
			return nil, err
		}
		user.ID = userCreated.ID
		user.CreatedAt = userCreated.CreatedAt
		user.UpdatedAt = userCreated.UpdatedAt

		// create team
		teamID, err := repositories.SaveTeam(ctx, tx, params.TeamName, userCreated.ID)
		if err != nil {
			return nil, err
		}

		// add user to its team
		err = repositories.AddListOfUsersToTeam(ctx, tx, teamID, []string{userCreated.ID}, string(modelsv1.RoleAdmin))
		if err != nil {
			return nil, err
		}
	} else {
		userCreated, err := repositories.CreateUser(ctx, tx, *user)
		if err != nil {
			return nil, err
		}
		user.ID = userCreated.ID
		user.CreatedAt = userCreated.CreatedAt
		user.UpdatedAt = userCreated.UpdatedAt
	}
	err = database.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, string(user.UserType))
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	// err = s.TokenRepo.Create(ctx, &models.RefreshToken{
	// 	UserID:    user.ID,
	// 	Token:     refreshToken,
	// 	ExpiresAt: utils.GetRefreshTokenExpiry(), // you can define this as helper
	// })
	// if err != nil {
	// 	return nil, err
	// }
	err = repositories.InsertRefreshTokenOfUser(ctx, user.ID, refreshToken, utils.GetRefreshTokenExpiry())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func Login(ctx context.Context, params modelsv1.LoginParams) (map[string]interface{}, error) {
	utils.Logger.Info("Attempting to login user with email: ", params.Email)

	// user, err := s.UserRepo.FindByEmail(ctx, params.Email)
	// if err != nil {
	// 	utils.Logger.Error("Error finding user by email: ", err)
	// 	return nil, errors.New("invalid credentials")
	// }
	user, err := repositories.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(params.Password, user.Password) {
		utils.Logger.Error("Error validating password for user with email: ", params.Email)
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, string(user.UserType))
	if err != nil {
		utils.Logger.Error("Error generating access token for user with email: ", params.Email)
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.Logger.Error("Error generating refresh token for user with email: ", params.Email)
		return nil, err
	}

	err = repositories.InsertRefreshTokenOfUser(ctx, user.ID, refreshToken, utils.GetRefreshTokenExpiry())
	if err != nil {
		utils.Logger.Error("Error inserting refresh token for user with ID : ", user.ID)
	}

	userDetail := &modelsv1.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserType:  modelsv1.UserType(user.UserType),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.Logger.Info("Logged in user successfully: ", user)
	return map[string]interface{}{
		"user":          userDetail,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func RefreshToken(ctx context.Context, token string) (map[string]string, error) {
	// refreshToken, err := s.TokenRepo.FindValidToken(ctx, token)
	// if err != nil {
	// 	return nil, err
	// }
	user, err := repositories.GetUserbyToken(ctx, token)
	if err != nil {
		return nil, err
	}
	newAccess, err := utils.GenerateAccessToken(user.ID, string(user.UserType))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token": newAccess,
	}, nil
}

func Logout(ctx context.Context, refreshToken string) error {
	// err := s.TokenRepo.RevokeAllForUser(ctx, refreshToken)
	// if err != nil {
	// 	return nil
	// }
	user, err := repositories.GetUserbyToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	repositories.RevokeTokenForUser(ctx, user.ID)

	return err
}
