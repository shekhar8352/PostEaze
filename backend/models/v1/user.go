package modelsv1

import (
	"time"
)

type SignupParams struct {
	Name     string   `json:"name" binding:"required,min=2"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	UserType UserType `json:"user_type" binding:"required,oneof=individual team"`
	TeamName string   `json:"team_name" binding:"required_if=UserType team"`
}

type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenParams struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleEditor  Role = "editor"
	RoleCreator Role = "creator"
)

type UserType string

const (
	UserTypeIndividual UserType = "individual"
	UserTypeTeam       UserType = "team"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	UserType  UserType  `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Memberships []TeamMember `json:"memberships,omitempty"`
}

type Team struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	OwnerID   string       `json:"owner_id"`
	Owner     User         `json:"owner"`
	Members   []TeamMember `json:"members"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type TeamMember struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`

	// Only keep Role here if users can have different roles across teams
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SuccessResponse is the standard success payload
type SuccessResponse struct {
    Status string      `json:"status"`
    Msg    string      `json:"msg"`
    Data   interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard error payload
type ErrorResponse struct {
    Status string `json:"status"`
    Msg    string `json:"msg"`
}

