package services

import (
	"context"
	"errors"
	"posteaze-backend/models"
	"posteaze-backend/pkg/config"
	"posteaze-backend/repository"
	"posteaze-backend/utils"

	"gorm.io/gorm"
)

type SignupParams struct {
	Name     string          `json:"name" binding:"required,min=2"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	UserType models.UserType `json:"user_type" binding:"required,oneof=individual team"`
	TeamName string          `json:"team_name" binding:"required_if=UserType team"`
}

type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthService struct {
	UserRepo repository.UserRepository
	TokenRepo repository.RefreshTokenRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		UserRepo: repository.NewUserRepository(config.GetAppContext().DB),
		TokenRepo: repository.NewRefreshTokenRepository(config.GetAppContext().DB),
	}
}

func (s *AuthService) Signup(ctx context.Context, params SignupParams) (map[string]interface{}, error) {
	db := config.GetAppContext().DB

	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: hashedPassword,
		UserType: params.UserType,
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if params.UserType == models.UserTypeTeam {
			team := models.Team{Name: params.TeamName}
			if err := tx.Create(&team).Error; err != nil {
				return err
			}
			user.TeamID = &team.ID
			if err := tx.Create(user).Error; err != nil {
				return err
			}
			team.OwnerID = user.ID
			return tx.Save(&team).Error
		}
		return tx.Create(user).Error
	})
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
	err = s.TokenRepo.Create(ctx, &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: utils.GetRefreshTokenExpiry(), // you can define this as helper
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}


func (s *AuthService) Login(ctx context.Context, params LoginParams) (map[string]interface{}, error) {
	utils.Logger.Info("Attempting to login user with email: ", params.Email)

	user, err := s.UserRepo.FindByEmail(ctx, params.Email)
	if err != nil {
		utils.Logger.Error("Error finding user by email: ", err)
		return nil, errors.New("invalid credentials")
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

	utils.Logger.Info("Logged in user successfully: ", user)
	return map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
