package api

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/flags"
)

// Init is used to initialise the router.
func Init() error {
	// create a server
	s := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", flags.Port())),
		server.WithKeepAlive(true),
	)

	v1 := s.Group(constants.V1Route)
	{
		addV1UserAuthRoutes(v1)
	}

	// this is a blocking call unless application receives shut down signal
	// or some error occurs
	return s.Run()
}

func addV1UserAuthRoutes(v1 *route.RouterGroup) {
	authv1 := v1.Group(constants.AuthRoute)

	authv1.POST(constants.SignUpRoute)
	authv1.POST(constants.LogInRoute)
	authv1.POST(constants.RefreshRoute)
	authv1.POST(constants.LogOutRoute)
}
