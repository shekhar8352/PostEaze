package api

import (
	"github.com/gin-gonic/gin"
	apiv1 "github.com/shekhar8352/PostEaze/api/v1"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/middleware"
)

// Init is used to initialise the router.
func Init() error {
	// create a server
	// s := server.Default(
	// 	server.WithHostPorts(fmt.Sprintf(":%d", flags.Port())),
	// 	server.WithKeepAlive(true),
	// )

	s := gin.Default()
	s.Use(middleware.GinLoggingMiddleware())

	api := s.Group(constants.ApiRoute)

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	
	v1 := api.Group(constants.V1Route)
	{
		addV1UserAuthRoutes(v1)
		addV1LogRoutes(v1)
	}

	// this is a blocking call unless application receives shut down signal
	// or some error occurs
	return s.Run()
}

func addV1UserAuthRoutes(v1 *gin.RouterGroup) {
	authv1 := v1.Group(constants.AuthRoute)

	authv1.POST(constants.SignUpRoute, apiv1.SignupHandler)
	authv1.POST(constants.LogInRoute, apiv1.LoginHandler)
	authv1.POST(constants.RefreshRoute, middleware.AuthMiddleware(), apiv1.RefreshTokenHandler)
	authv1.POST(constants.LogOutRoute, middleware.AuthMiddleware(), apiv1.LogoutHandler)
}

func addV1LogRoutes(v1 *gin.RouterGroup) {
	logv1 := v1.Group(constants.LogRoute)
	logv1.GET(constants.LogByDate, apiv1.GetLogsByDate)
	logv1.GET(constants.LogById, apiv1.GetLogByIDHandler)
}
