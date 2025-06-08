package services

import (
	"context"
	"errors"
	"posteaze-backend/models"
	"posteaze-backend/pkg/config"
	"posteaze-backend/utils"

	"gorm.io/gorm"
)

type SignupParams struct {
	Name     string
	Email    string
	Password string
	UserType models.UserType
	TeamName string
}

type LoginParams struct {
	Email    string
	Password string
}

type AuthService struct{}

func (s *AuthService) Signup(ctx context.Context, params SignupParams) (*models.User, error) {
	db := config.GetAppContext().DB

	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
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
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			team.OwnerID = user.ID
			return tx.Save(&team).Error
		}
		return tx.Create(&user).Error
	})

	return &user, err
}

func (s *AuthService) Login(ctx context.Context, params LoginParams) (*models.User, error) {
	db := config.GetAppContext().DB

	var user models.User
	if err := db.WithContext(ctx).Where("email = ?", params.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(params.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
