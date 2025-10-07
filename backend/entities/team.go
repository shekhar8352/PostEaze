package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/database"
)

const (
	CreateTeam = iota
	GetAllTeams
	GetTeamByID
	GetTeamByOwnerID
	UpdateTeam
)

type Team struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	OwnerID           string          `json:"owner_id"`
	Description       *string         `json:"description"`
	AvatarURL         *string         `json:"avatar_url"`
	Visibility        string          `json:"visibility"`
	Status            string          `json:"status"`
	OwnerRoleOverride *string         `json:"owner_role_override"`
	Settings          json.RawMessage `json:"settings"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
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

// --- Validation helpers ---

func (o *Team) Validate() error {
	if o.Visibility == "" {
		o.Visibility = "private"
	}

	if o.Status == "" {
		o.Status = "active"
	}
	validStatus := map[string]bool{"active": true, "archived": true, "deleted": true}
	if !validStatus[o.Status] {
		return errors.New("invalid team status")
	}

	if len(o.Settings) == 0 {
		o.Settings = json.RawMessage(`{}`)
	}
	return nil
}

// --- Query Builders ---

func (o *Team) GetQuery(code int) string {
	switch code {
	case CreateTeam:
		return `
		INSERT INTO teams (name, owner_id, visibility, description, avatar_url, settings, status)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, '{}'::jsonb), $7)
		RETURNING id, created_at, updated_at;`
	case GetTeamByID:
		return `
		SELECT id, name, visibility, description, avatar_url, status, owner_id, settings, created_at, updated_at
		FROM teams WHERE id = $1;`
	case UpdateTeam:
		return `
			UPDATE teams 
			SET name=$1, description=$2, avatar_url=$3, visibility=$4, updated_at=NOW()
			WHERE id=$5 
			RETURNING updated_at, owner_id, settings, created_at, status;`
	}
	return constants.Empty
}

func (o *Team) GetQueryValues(code int) []any {
	switch code {
	case CreateTeam:
		_ = o.Validate()
		return []any{o.Name, o.OwnerID, o.Visibility, o.Description, o.AvatarURL, o.Settings, o.Status}
	case GetTeamByID:
		return []any{o.ID}
	case UpdateTeam:
		fmt.Println(o)
		return []any{o.Name, o.Description, o.AvatarURL, o.Visibility, o.ID}
	}
	return nil
}

func (o *Team) GetMultiQuery(code int) string {
	switch code {
	case GetAllTeams:
		return `
		SELECT id, name, visibility, description, avatar_url, status, owner_id, settings, created_at, updated_at
		FROM teams;`
	case GetTeamByOwnerID:
		return `
		SELECT id, name, visibility, description, avatar_url, status, owner_id, settings, created_at, updated_at
		FROM teams WHERE owner_id = $1;`
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

// --- Scanner bindings ---

func (o *Team) GetNextRaw() database.RawEntity {
	return new(Team)
}

func (o *Team) BindRawRow(code int, row database.Scanner) error {
	switch code {
	case CreateTeam:
		return row.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
	case GetAllTeams:
		return row.Scan(&o.ID, &o.Name, &o.Visibility, &o.Description, &o.AvatarURL,
			&o.Status, &o.OwnerID, &o.Settings, &o.CreatedAt, &o.UpdatedAt)
	case GetTeamByID:
		return row.Scan(&o.ID, &o.Name, &o.Visibility, &o.Description, &o.AvatarURL,
			&o.Status, &o.OwnerID, &o.Settings, &o.CreatedAt, &o.UpdatedAt)
	case GetTeamByOwnerID:
		return row.Scan(&o.ID, &o.Name, &o.Visibility, &o.Description, &o.AvatarURL,
			&o.Status, &o.OwnerID, &o.Settings, &o.CreatedAt, &o.UpdatedAt)
	case UpdateTeam:
		return row.Scan(&o.UpdatedAt, &o.OwnerID, &o.Settings, &o.CreatedAt, &o.Status)
	}
	return nil
}

func (o *Team) GetExec(code int) string {
	return constants.Empty
}

func (o *Team) GetExecValues(code int, _ string) []any {
	return nil
}
