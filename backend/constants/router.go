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

	TeamRoute        = "/team"
	GetAllTeams      = "/all"
	GetTeamById      = "/:id"
	GetTeamByOwnerID = "/owner/:id"
	CreateTeam       = "/create"
	UpdateTeam       = "/update"
	UpdateTeamStatus = "/update-status"

	MetaRoute    = "/meta"
	MetaCallback = "/callback"

	ChannelRoute           = "/channels"
	InstagramRoute         = "/instagram"
	CreateInstagramChannel = "/create"
	WebhookRoute           = "/webhooks"
	InstagramWebhook       = "/instagram"

	ScheduledPostsRoute = "/scheduled-posts"
)
