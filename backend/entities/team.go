package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/database"
)

const (
	CreateTeam = iota
	AddUsersToTeam
)

type Team struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	OwnerID   string       `json:"owner_id"`
	Members   []TeamMember `json:"members"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type TeamMember struct {
	UserID string `json:"user_id"`
	// Only keep Role here if users can have different roles across teams
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (o *Team) GetQuery(code int) string {
	switch code {
	case CreateTeam:
		return `INSERT INTO teams (name , owner_id ) VALUES ( $1 , $2) RETURNING id;`
	case AddUsersToTeam:
		baseQuery := `INSERT INTO team_members (team_id, user_id, role) VALUES `
		valueStrings := make([]string, 0, len(o.Members))
		argCounter := 1

		for range o.Members {
			// For each member, add 3 placeholders: ($1, $2, $3), ($4, $5, $6), ...
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", argCounter, argCounter+1, argCounter+2))
			argCounter += 3
		}

		placeholderString := strings.Join(valueStrings, ", ")
		return baseQuery + placeholderString + ";"

	}
	return constants.Empty
}

func (o *Team) GetQueryValues(code int) []any {
	switch code {
	case CreateTeam:
		return []interface{}{o.Name, o.OwnerID}
	case AddUsersToTeam:
		args := make([]interface{}, 0, len(o.Members)*3)
		for _, member := range o.Members {
			args = append(args, o.ID, member.UserID, member.Role)
		}
		return args
	}
	return nil
}

func (o *Team) GetMultiQuery(code int) string {
	switch code {
	}
	return constants.Empty
}

func (o *Team) GetMultiQueryValues(code int) []any {
	switch code {
	}
	return nil
}

func (o *Team) GetNextRaw() database.RawEntity {
	return new(Team)
}

func (o *Team) BindRawRow(code int, row database.Scanner) error {
	switch code {
	case CreateTeam:
		row.Scan(&o.ID)
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
