package entities

import (
	"encoding/json"
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/database"
)

const (
	CreateTeamMember = iota + 100 // offset from team constants
	GetTeamMemberByID
	GetMembersByTeamID
	GetTeamsByUserID
	UpdateTeamMemberRole
	UpdateTeamMemberStatus
	UpdateTeamMemberPermissions
	UpdateTeamMemberLastActive
	DeleteTeamMember
	GetAllTeamMembers
)

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

// --- Validation ---

func (m *TeamMember) Validate() error {
	if m.Status == "" {
		m.Status = "active"
	}
	if m.Role == "" {
		m.Role = "member"
	}
	if m.Permissions == nil {
		m.Permissions = map[string]interface{}{}
	}
	return nil
}

// --- Query Builders ---

func (m *TeamMember) GetQuery(code int) string {
	switch code {
	case CreateTeamMember:
		return `
			INSERT INTO team_members 
				(team_id, user_id, role, status, invited_by, permissions, is_primary, joined_at)
			VALUES 
				($1, $2, $3, $4, $5, COALESCE($6, '{}'::jsonb), $7, NOW())
			RETURNING id, created_at, updated_at;`
	case GetTeamMemberByID:
		return `
			SELECT id, team_id, user_id, role, status, joined_at, invited_by, permissions, 
			       is_primary, last_active_at, created_at, updated_at
			FROM team_members
			WHERE id = $1;`
	case UpdateTeamMemberRole:
		return `
			UPDATE team_members
			SET role = $1, updated_at = NOW()
			WHERE id = $2
			RETURNING id, role, updated_at;`
	case UpdateTeamMemberStatus:
		return `
			UPDATE team_members
			SET status = $1, updated_at = NOW()
			WHERE id = $2
			RETURNING id, status, updated_at;`
	case UpdateTeamMemberPermissions:
		return `
			UPDATE team_members
			SET permissions = COALESCE($1, '{}'::jsonb), updated_at = NOW()
			WHERE id = $2
			RETURNING id, permissions, updated_at;`
	case UpdateTeamMemberLastActive:
		return `
			UPDATE team_members
			SET last_active_at = NOW(), updated_at = NOW()
			WHERE id = $1
			RETURNING id, last_active_at, updated_at;`
	case DeleteTeamMember:
		return `
			DELETE FROM team_members
			WHERE id = $1
			RETURNING id;`
	}
	return constants.Empty
}

func (m *TeamMember) GetQueryValues(code int) []any {
	switch code {
	case CreateTeamMember:
		_ = m.Validate()
		permissionsJSON, _ := json.Marshal(m.Permissions)
		return []any{m.TeamID, m.UserID, m.Role, m.Status, m.InvitedBy, permissionsJSON, m.IsPrimary}
	case GetTeamMemberByID, DeleteTeamMember, UpdateTeamMemberLastActive:
		return []any{m.ID}
	case UpdateTeamMemberRole:
		return []any{m.Role, m.ID}
	case UpdateTeamMemberStatus:
		return []any{m.Status, m.ID}
	case UpdateTeamMemberPermissions:
		permissionsJSON, _ := json.Marshal(m.Permissions)
		return []any{permissionsJSON, m.ID}
	}
	return nil
}

// --- Multi Query Builders ---

func (m *TeamMember) GetMultiQuery(code int) string {
	switch code {
	case GetMembersByTeamID:
		return `
			SELECT id, team_id, user_id, role, status, joined_at, invited_by, permissions, 
			       is_primary, last_active_at, created_at, updated_at
			FROM team_members
			WHERE team_id = $1
			ORDER BY joined_at ASC;`
	case GetTeamsByUserID:
		return `
			SELECT id, team_id, user_id, role, status, joined_at, invited_by, permissions, 
			       is_primary, last_active_at, created_at, updated_at
			FROM team_members
			WHERE user_id = $1
			ORDER BY joined_at ASC;`
	case GetAllTeamMembers:
		return `
			SELECT id, team_id, user_id, role, status, joined_at, invited_by, permissions, 
			       is_primary, last_active_at, created_at, updated_at
			FROM team_members
			ORDER BY created_at DESC;`
	}
	return constants.Empty
}

func (m *TeamMember) GetMultiQueryValues(code int) []any {
	switch code {
	case GetMembersByTeamID:
		return []any{m.TeamID}
	case GetTeamsByUserID:
		return []any{m.UserID}
	case GetAllTeamMembers:
		return []any{}
	}
	return nil
}

// --- Scanner Bindings ---

func (m *TeamMember) GetNextRaw() database.RawEntity {
	return new(TeamMember)
}

func (m *TeamMember) BindRawRow(code int, row database.Scanner) error {
	switch code {
	case CreateTeamMember:
		return row.Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	case GetTeamMemberByID, GetMembersByTeamID, GetTeamsByUserID, GetAllTeamMembers:
		return row.Scan(
			&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.Status,
			&m.JoinedAt, &m.InvitedBy, &m.Permissions, &m.IsPrimary,
			&m.LastActiveAt, &m.CreatedAt, &m.UpdatedAt,
		)
	case UpdateTeamMemberRole:
		return row.Scan(&m.ID, &m.Role, &m.UpdatedAt)
	case UpdateTeamMemberStatus:
		return row.Scan(&m.ID, &m.Status, &m.UpdatedAt)
	case UpdateTeamMemberPermissions:
		return row.Scan(&m.ID, &m.Permissions, &m.UpdatedAt)
	case UpdateTeamMemberLastActive:
		return row.Scan(&m.ID, &m.LastActiveAt, &m.UpdatedAt)
	case DeleteTeamMember:
		return row.Scan(&m.ID)
	}
	return nil
}

func (m *TeamMember) GetExec(code int) string {
	return constants.Empty
}

func (m *TeamMember) GetExecValues(code int, _ string) []any {
	return nil
}
