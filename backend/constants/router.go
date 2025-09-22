package constants

const (
	ApiRoute = "/api"
	V1Route  = "/v1"

	AuthRoute    = "/auth"
	Authenticate = "/authenticate"
	RefreshRoute = "/refresh"
	LogOutRoute  = "/logout"
	LogRoute     = "/log"

	LogById   = "/byId/:log_id"
	LogByDate = "/byDate/:date"

	UserRoute   = "/user"
	GetUserById = "/:user_id"
	UpdateUser  = "/:user_id"

	TeamRoute   = "/team"
	GetAllTeams = "/all"
	GetTeamById = "/:team_id"
	CreateTeam  = "/create"
	UpdateTeam  = "/:team_id"
)
