package modelsv1

import (
	"time"
)

type FirebaseAuthParams struct {
	FirebaseID       string `json:"firebase_id" binding:"required"`
	FirebaseToken string `json:"firebase_token" binding:"required"`
	Platform      string `json:"platform" binding:"required,oneof=email google facebook microsoft"`
}

type RefreshTokenParams struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserType string

const (
	UserTypeIndividual UserType = "individual"
	UserTypeTeam       UserType = "team"
)

type User struct {
	ID        string    `json:"id"`
	FirebaseID string    `json:"firebase_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Platforms []string  `json:"platforms"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Memberships []TeamMember `json:"memberships,omitempty"`
}

type UpdateUserParams struct {
	ID string `json:"id"`
	Email string `json:"email"`
}