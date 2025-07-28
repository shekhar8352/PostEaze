package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func SaveTeam(ctx context.Context, tx database.Database, teamName string, ownerID string) (string, error) {

	data := entities.Team{
		Name:    teamName,
		OwnerID: ownerID,
	}
	err := tx.QueryRaw(ctx, &data, entities.CreateTeam)
	return data.ID, err
}

func AddListOfUsersToTeam(ctx context.Context, tx database.Database, teamID string, members []string, role string) error {
	data := entities.Team{
		ID: teamID,
	}
	for i := range members {
		data.Members = append(data.Members, entities.TeamMember{UserID: members[i], Role: role})
	}
	err := tx.QueryRaw(ctx, &data, entities.AddUsersToTeam)
	return err

}
