package businessv1

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/database"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func CreateTeam(ctx context.Context, body modelsv1.Team) (*entities.Team, error) {
	tx, err := database.GetTx(ctx, nil)
	if err != nil {
		utils.Logger.Error(ctx, "Error creating team: %v", err)
		return nil, err
	}

	team, err := repositories.CreateTeam(ctx, tx, body.Name, body.OwnerID)
	if err != nil {
		utils.Logger.Error(ctx, "Error creating team: %v", err)
		return nil, err
	}

	err = database.CommitTx(tx)
	if err != nil {
		utils.Logger.Error(ctx, "Error committing transaction: %v", err)
		return nil, err
	}

	return team, nil
}

func GetAllTeams(ctx context.Context) ([]*entities.Team, error) {
	teams, err := repositories.GetAllTeams(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Error fetching teams: %v", err)
		return nil, err
	}
	return teams, nil
}

func GetTeamByID(ctx context.Context, teamID string) (*entities.Team, error) {
	team, err := repositories.GetTeamByID(ctx, teamID)
	if err != nil {
		utils.Logger.Error(ctx, "Error fetching team by ID: %v", err)
		return nil, err
	}
	return team, nil
}
