package entities

import (
	"encoding/json"
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/database"
)

const (
	CreateTeam = iota
	GetAllTeams
	GetTeamByID
	GetTeamByOwnerID
)

type Team struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	OwnerID           string                 `json:"owner_id"`
	Description       *string                `json:"description"`
	AvatarURL         *string                `json:"avatar_url"`
	Visibility        string                 `json:"visibility"`
	Status            string                 `json:"status"`
	OwnerRoleOverride *string                `json:"owner_role_override"`
	Settings          json.RawMessage `json:"settings"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}


type TeamMember struct {
	ID           string                 `json:"id"`
	TeamID       string                 `json:"team_id"`
	UserID       string                 `json:"user_id"`
	Role         string                 `json:"role"`
	Status       string                 `json:"status"`
	JoinedAt     time.Time              `json:"joined_at"`
	InvitedBy    *string                `json:"invited_by,omitempty"`
	Permissions  map[string]interface{} `json:"permissions"`
	IsPrimary    bool                   `json:"is_primary"`
	LastActiveAt *time.Time             `json:"last_active_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (o *Team) GetQuery(code int) string {
	switch code {
	case CreateTeam:
		return `INSERT INTO teams (name , owner_id, visibility, description, avatar_url ) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at;`
	case GetTeamByID:
		return `SELECT id, name, visibility, description, avatar_url, status, owner_id, settings, created_at, updated_at FROM teams WHERE id = $1;`
	}
	return constants.Empty
}

func (o *Team) GetQueryValues(code int) []any {
	switch code {
	case CreateTeam:
		return []any{o.Name, o.OwnerID, o.Visibility, o.Description, o.AvatarURL}
	case GetTeamByID:
		return []any{o.ID}
	}
	return nil
}

func (o *Team) GetMultiQuery(code int) string {
	switch code {
	case GetAllTeams:
		return `SELECT id, name, visibility, description, avatar_url, status, owner_id, created_at, updated_at FROM teams;`
	case GetTeamByOwnerID:
		return `SELECT id, name, visibility, description, avatar_url, status, owner_id, created_at, updated_at FROM teams WHERE owner_id = $1;`
	}
	return constants.Empty
}

func (o *Team) GetMultiQueryValues(code int) []any {
	switch code {
	case GetAllTeams:
		return []any{}
	case GetTeamByOwnerID:
		return []any{o.OwnerID}
	}
	return nil
}

func (o *Team) GetNextRaw() database.RawEntity {
	return new(Team)
}

func (o *Team) BindRawRow(code int, row database.Scanner) error {
	switch code {
	case CreateTeam:
		return row.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
	case GetAllTeams:
		return row.Scan(&o.ID, &o.Name, &o.Visibility, &o.Description, &o.AvatarURL, &o.Status, &o.OwnerID, &o.CreatedAt, &o.UpdatedAt)
	case GetTeamByID:
    return row.Scan(
        &o.ID,
        &o.Name,
        &o.Visibility,
        &o.Description,
        &o.AvatarURL,
        &o.Status,
        &o.OwnerID,
        &o.Settings,
        &o.CreatedAt,
        &o.UpdatedAt,
    )
	case GetTeamByOwnerID:
		return row.Scan(&o.ID, &o.Name, &o.Visibility, &o.Description, &o.AvatarURL, &o.Status, &o.Settings, &o.CreatedAt, &o.UpdatedAt)
	}
	return nil
}

func (o *Team) GetExec(code int) string {
	switch code {
	default:
		return constants.Empty
	}
}

func (o *Team) GetExecValues(code int, _ string) []any {
	return nil
}
