package modelsv1

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleEditor  Role = "editor"
	RoleCreator Role = "creator"
)

type Team struct {
	Name      string       `json:"name"`
	Description string       `json:"description"`
	Visibility string       `json:"visibility"`
	AvatarURL  string       `json:"avatar_url"`
	OwnerID   string       `json:"owner_id"`
	Status    string       `json:"status"`
}

type TeamMember struct {
	ID           string                 `json:"id"`
	TeamID       string                 `json:"team_id"`
	UserID       string                 `json:"user_id"`
	Role         string                 `json:"role"`
	Status       string                 `json:"status"`
	InvitedBy    *string                `json:"invited_by,omitempty"`
	Permissions  map[string]interface{} `json:"permissions"`
	IsPrimary    bool                   `json:"is_primary"`
}