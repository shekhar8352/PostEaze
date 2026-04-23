package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreateTeamMember adds a new member to a team.
func CreateTeamMember(ctx context.Context, tx database.Database, body modelsv1.TeamMember) (*entities.TeamMember, error) {
	data := entities.TeamMember{
		TeamID:      body.TeamID,
		UserID:      body.UserID,
		Role:        body.Role,
		Status:      body.Status,
		InvitedBy:   body.InvitedBy,
		Permissions: body.Permissions,
		IsPrimary:   body.IsPrimary,
	}
	err := tx.QueryRaw(ctx, &data, entities.CreateTeamMember)
	return &data, err
}

// GetAllTeamMembers retrieves all members for all teams (admin-level view).
func GetAllTeamMembers(ctx context.Context) ([]*entities.TeamMember, error) {
	rows, err := database.Get().QueryMultiRaw(ctx, &entities.TeamMember{}, entities.GetAllTeamMembers)
	if err != nil {
		return nil, err
	}

	members := make([]*entities.TeamMember, 0, len(rows))
	for _, row := range rows {
		if member, ok := row.(*entities.TeamMember); ok {
			members = append(members, member)
		}
	}

	return members, nil
}

// GetTeamMembersByTeamID retrieves all members of a specific team.
func GetTeamMembersByTeamID(ctx context.Context, teamID string) ([]*entities.TeamMember, error) {
	data := entities.TeamMember{
		TeamID: teamID,
	}

	rows, err := database.Get().QueryMultiRaw(ctx, &data, entities.GetMembersByTeamID)
	if err != nil {
		return nil, err
	}

	members := make([]*entities.TeamMember, 0, len(rows))
	for _, row := range rows {
		if member, ok := row.(*entities.TeamMember); ok {
			members = append(members, member)
		}
	}

	return members, nil
}

// GetTeamMemberByID retrieves a specific team member by ID.
func GetTeamMemberByID(ctx context.Context, memberID string) (*entities.TeamMember, error) {
	data := entities.TeamMember{
		ID: memberID,
	}
	err := database.Get().QueryRaw(ctx, &data, entities.GetTeamMemberByID)
	return &data, err
}

// GetTeamsByUserID retrieves all teams a user belongs to.
func GetTeamsByUserID(ctx context.Context, userID string) ([]*entities.TeamMember, error) {
	data := entities.TeamMember{
		UserID: userID,
	}

	rows, err := database.Get().QueryMultiRaw(ctx, &data, entities.GetTeamsByUserID)
	if err != nil {
		return nil, err
	}

	memberships := make([]*entities.TeamMember, 0, len(rows))
	for _, row := range rows {
		if membership, ok := row.(*entities.TeamMember); ok {
			memberships = append(memberships, membership)
		}
	}

	return memberships, nil
}

// UpdateTeamMemberRole updates the role of a specific team member.
func UpdateTeamMemberRole(ctx context.Context, tx database.Database, body modelsv1.TeamMember, memberID string) (*entities.TeamMember, error) {
	data := entities.TeamMember{
		ID:   memberID,
		Role: body.Role,
	}
	err := tx.QueryRaw(ctx, &data, entities.UpdateTeamMemberRole)
	return &data, err
}

// UpdateTeamMemberStatus updates the status of a specific team member.
func UpdateTeamMemberStatus(ctx context.Context, tx database.Database, body modelsv1.TeamMember, memberID string) (*entities.TeamMember, error) {
	data := entities.TeamMember{
		ID:     memberID,
		Status: body.Status,
	}
	err := tx.QueryRaw(ctx, &data, entities.UpdateTeamMemberStatus)
	return &data, err
}

// UpdateTeamMemberPermissions updates the permissions of a specific team member.
func UpdateTeamMemberPermissions(ctx context.Context, tx database.Database, body modelsv1.TeamMember, memberID string) (*entities.TeamMember, error) {
	data := entities.TeamMember{
		ID:          memberID,
		Permissions: body.Permissions,
	}
	err := tx.QueryRaw(ctx, &data, entities.UpdateTeamMemberPermissions)
	return &data, err
}

// DeleteTeamMember removes a member from a team (hard delete or status=removed).
func DeleteTeamMember(ctx context.Context, tx database.Database, memberID string) error {
	data := entities.TeamMember{
		ID: memberID,
	}
	return tx.QueryRaw(ctx, &data, entities.DeleteTeamMember)
}
