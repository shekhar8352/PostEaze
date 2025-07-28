package models

import (
	"time"

	"github.com/google/uuid"
)

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
	ID        uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string        `json:"name"`
	Email     string        `gorm:"uniqueIndex" json:"email"`
	Password  string        `json:"-"`
	UserType  UserType      `gorm:"type:varchar(20)" json:"user_type"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`

	TeamID *uuid.UUID `gorm:"type:uuid" json:"team_id,omitempty"`
	Team   *Team      `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"team,omitempty"`

	Memberships []TeamMember `gorm:"foreignKey:UserID" json:"memberships,omitempty"`
}

type Team struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string       `json:"name"`
	OwnerID   uuid.UUID    `gorm:"type:uuid" json:"owner_id"`
	Owner     User         `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"owner"`
	Members   []TeamMember `gorm:"foreignKey:TeamID" json:"members"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type TeamMember struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeamID    uuid.UUID `gorm:"type:uuid" json:"team_id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`

	// Only keep Role here if users can have different roles across teams
	Role Role `gorm:"type:varchar(20)" json:"role"`

	Team      Team `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"team,omitempty"`
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

