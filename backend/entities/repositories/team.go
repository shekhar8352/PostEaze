package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func CreateTeam(ctx context.Context, tx database.Database, teamName string, ownerID string) (*entities.Team, error) { 
	data := entities.Team{
		Name:    teamName,
		OwnerID: ownerID,
	}
	err := tx.QueryRaw(ctx, &data, entities.CreateTeam)
	return &data, err
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

func GetAllTeams(ctx context.Context) ([]*entities.Team, error) {
    rows, err := database.Get().QueryMultiRaw(ctx, &entities.Team{}, entities.GetAllTeams)
    if err != nil {
        return nil, err
    }

    teams := make([]*entities.Team, 0, len(rows))
    for _, row := range rows {
        if team, ok := row.(*entities.Team); ok {
            teams = append(teams, team)
        }
    }

    return teams, nil
}

func GetTeamByID(ctx context.Context, teamID string) (*entities.Team, error) {
	data := entities.Team{
		ID: teamID,
	}
	err := database.Get().QueryRaw(ctx, &data, entities.GetTeamByID)
	return &data, err
}

func GetTeamByOwnerID(ctx context.Context, ownerID string) ([]*entities.Team, error) {
	data := entities.Team{
		OwnerID: ownerID,
	}

    rows, err := database.Get().QueryMultiRaw(ctx, &data, entities.GetTeamByOwnerID)
    if err != nil {
        return nil, err
    }

    teams := make([]*entities.Team, 0, len(rows))
    for _, row := range rows {
        if team, ok := row.(*entities.Team); ok {
            teams = append(teams, team)
        }
    }

    return teams, nil
}