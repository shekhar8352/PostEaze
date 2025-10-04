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
}

type TeamMember struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
	// Only keep Role here if users can have different roles across teams
	Role      Role      `json:"role"`
}