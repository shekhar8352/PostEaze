package businessv1

import (
	"context"
	"errors"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func Signup(ctx context.Context, params modelsv1.SignupParams) (map[string]interface{}, error) {
	// Note: This function is kept for backward compatibility
	// New users should use Firebase authentication
	
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &modelsv1.User{
		Name:      params.Name,
		Email:     params.Email,
		Platforms: []string{"email"}, // Traditional email signup
	}
	
	// Create legacy user entity for backward compatibility
	legacyUser := &entities.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		UserType: string(params.UserType),
	}
	
	// start transaction
	tx, err := database.GetTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	if params.UserType == modelsv1.UserTypeTeam {
		// create user
		userCreated, err := repositories.CreateUser(ctx, tx, *legacyUser)
		if err != nil {
			database.RollbackTx(tx)
			return nil, err
		}
		user.ID = userCreated.ID
		user.CreatedAt = userCreated.CreatedAt
		user.UpdatedAt = userCreated.UpdatedAt

		// create team
		teamID, err := repositories.SaveTeam(ctx, tx, params.TeamName, userCreated.ID)
		if err != nil {
			database.RollbackTx(tx)
			return nil, err
		}

		// add user to its team
		err = repositories.AddListOfUsersToTeam(ctx, tx, teamID, []string{userCreated.ID}, string(modelsv1.RoleAdmin))
		if err != nil {
			database.RollbackTx(tx)
			return nil, err
		}
	} else {
		userCreated, err := repositories.CreateUser(ctx, tx, *legacyUser)
		if err != nil {
			database.RollbackTx(tx)
			return nil, err
		}
		user.ID = userCreated.ID
		user.CreatedAt = userCreated.CreatedAt
		user.UpdatedAt = userCreated.UpdatedAt
	}
	
	err = database.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, string(params.UserType))
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	err = repositories.InsertRefreshTokenOfUser(ctx, user.ID, refreshToken, utils.GetRefreshTokenExpiry())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func Login(ctx context.Context, params modelsv1.LoginParams) (map[string]interface{}, error) {
	// Note: This function is kept for backward compatibility with password-based authentication
	utils.Logger.Info(ctx, "Attempting to login user with email: %s", params.Email)
	user, err := repositories.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(params.Password, user.Password) {
		utils.Logger.Error(ctx, "Error validating password for user with email: %s", params.Email)
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.UserType)
	if err != nil {
		utils.Logger.Error(ctx, "Error generating access token for user with email: %s", params.Email)
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.Logger.Error(ctx, "Error generating refresh token for user with email: %s", params.Email)
		return nil, err
	}

	err = repositories.InsertRefreshTokenOfUser(ctx, user.ID, refreshToken, utils.GetRefreshTokenExpiry())
	if err != nil {
		utils.Logger.Error(ctx, "Error inserting refresh token for user with ID : %s", user.ID)
	}

	userDetail := &modelsv1.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Platforms: []string{"email"}, // Legacy users use email platform
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.Logger.Info(ctx, "Logged in user successfully: %s", user)
	return map[string]interface{}{
		"user":          userDetail,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func RefreshToken(ctx context.Context, token string) (map[string]string, error) {
	user, err := repositories.GetUserbyToken(ctx, token)
	if err != nil {
		return nil, err
	}
	
	// For Firebase users, use "individual" as default user type
	userType := "individual"
	if user.UserType != "" {
		userType = user.UserType
	}
	
	newAccess, err := utils.GenerateAccessToken(user.ID, userType)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token": newAccess,
	}, nil
}

func Logout(ctx context.Context, refreshToken string) error {
	user, err := repositories.GetUserbyToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	repositories.RevokeTokenForUser(ctx, user.ID)

	return err
}

func AuthenticateWithFirebase(ctx context.Context, params modelsv1.FirebaseAuthParams) (map[string]interface{}, error) {
	// Validate Firebase token
	firebaseService := utils.GetFirebaseService()
	if firebaseService == nil {
		return nil, errors.New("Firebase service not initialized")
	}

	firebaseUser, err := firebaseService.ValidateToken(ctx, params.FirebaseToken)
	if err != nil {
		utils.Logger.Error(ctx, "Firebase token validation failed: %v", err)
		return nil, errors.New("invalid Firebase token")
	}

	// Validate email requirement based on platform
	if params.Platform != "facebook" && firebaseUser.Email == "" {
		utils.Logger.Error(ctx, "Email is required for this platform: %s", params.Platform)
		return nil, errors.New("email is required for this platform")
	}

	// 🔥 Check if user exists by firebase_id instead of local ID
	existingUser, err := repositories.GetUserByFirebaseID(ctx, firebaseUser.UID)
	if err != nil {
		if errors.Is(err, database.ErrNoRecords) {
			// User doesn't exist → create new user
			utils.Logger.Info(ctx, "User does not exist, creating new user with FirebaseID: %s", firebaseUser.UID)
			return createNewFirebaseUser(ctx, params, firebaseUser)
		}
		// Other DB errors
		utils.Logger.Error(ctx, "Database error while checking user existence: %v", err)
		return nil, err
	}

	// User exists → update platforms if needed
	utils.Logger.Info(ctx, "User exists with FirebaseID: %s, ID: %s, Name: %s, Email: %s, Platforms: %v",
		existingUser.FirebaseID, existingUser.ID, existingUser.Name, existingUser.Email, existingUser.Platforms)

	return authenticateExistingUser(ctx, existingUser, params.Platform)
}

func createNewFirebaseUser(ctx context.Context, params modelsv1.FirebaseAuthParams, firebaseUser *utils.FirebaseUser) (map[string]interface{}, error) {
	// Start transaction
	tx, err := database.GetTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user := modelsv1.User{
		FirebaseID: firebaseUser.UID,
		Name:       firebaseUser.Name,
		Email:      firebaseUser.Email,
		Platforms:  []string{params.Platform},
	}

	userCreated, err := repositories.CreateUserWithFirebase(ctx, tx, user)
	if err != nil {
		database.RollbackTx(tx)
		return nil, err
	}
	user.ID = userCreated.ID
	user.CreatedAt = userCreated.CreatedAt
	user.UpdatedAt = userCreated.UpdatedAt

	err = database.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(userCreated.ID, "individual")
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(userCreated.ID)
	if err != nil {
		return nil, err
	}

	err = repositories.InsertRefreshTokenOfUser(ctx, userCreated.ID, refreshToken, utils.GetRefreshTokenExpiry())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func authenticateExistingUser(ctx context.Context, existingUser *entities.User, platform string) (map[string]interface{}, error) {
	// Ensure platform is in user's list
	platformExists := false
	for _, p := range existingUser.Platforms {
		if p == platform {
			platformExists = true
			break
		}
	}

	if !platformExists {
		updatedPlatforms := append(existingUser.Platforms, platform)
		err := repositories.UpdateUserPlatforms(ctx, existingUser.ID, updatedPlatforms)
		if err != nil {
			return nil, err
		}
		existingUser.Platforms = updatedPlatforms
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(existingUser.ID, "individual")
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		return nil, err
	}

	err = repositories.InsertRefreshTokenOfUser(ctx, existingUser.ID, refreshToken, utils.GetRefreshTokenExpiry())
	if err != nil {
		return nil, err
	}

	userDetail := &modelsv1.User{
		ID:         existingUser.ID,
		FirebaseID: existingUser.FirebaseID,
		Name:       existingUser.Name,
		Email:      existingUser.Email,
		Platforms:  existingUser.Platforms,
		CreatedAt:  existingUser.CreatedAt,
		UpdatedAt:  existingUser.UpdatedAt,
	}

	return map[string]interface{}{
		"user":          userDetail,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}