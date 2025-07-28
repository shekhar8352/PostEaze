package main

import (
	"context"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/api"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils/configs"
	sql "github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/shekhar8352/PostEaze/utils/flags"
	httpclient "github.com/shekhar8352/PostEaze/utils/http"
)

func main() {
	// this main function initializes all the services required to run the application e.g database , router , configs , redis etc
	ctx := context.Background()
	initEnv()
	initConfigs(ctx)
	initDatabase(ctx)
	initRouter(ctx)
	initHttp(ctx)
}
func initEnv() {
	env.InitEnv()
}
func initConfigs(ctx context.Context) {
	var err error
	configNames := []string{constants.APIConfig, constants.ApplicationConfig, constants.DatabaseConfig}
	if flags.Mode() == constants.DevMode {
		err = configs.InitDev(flags.BaseConfigPath(), configNames...)
	} else if flags.Mode() == constants.ReleaseMode {
		err = configs.InitRelease(flags.Env(), flags.AWSRegion(), configNames...)
	}
	if err != nil {
		log.Fatal(ctx, "error in initialising configs", err)
	}
}

func initDatabase(ctx context.Context) {
	// to be done
	driverName, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseDriverNameConfigKey)
	if err != nil {
		log.Fatal(ctx, " failed to get database driver name ", err)
		return
	}
	urlString, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseURLConfigKey)
	if err != nil {
		log.Fatal(ctx, " failed to get database url ", err)
		return
	}
	url := env.ApplyEnvironmentToString(urlString)
	err = sql.Init(ctx, sql.Config{
		DriverName: driverName,
		URL:        url,
		MaxOpenConnections: int(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseMaxOpenConnectionsConfigKey, 1)),
		MaxIdleConnections: int(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseMaxIdleConnectionsConfigKey, 0)),
		ConnectionMaxLifetime: time.Duration(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseMaxConnectionLifetimeInSecondsConfigKey, 30)) * time.Second,
		ConnectionMaxIdleTime: time.Duration(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseMaxConnectionIdleTimeInSecondsConfigKey, 10)) * time.Second,
	})
	if err != nil {
		log.Fatal(ctx, " failed to initialize database ", err)
	}
}

func initHttp(ctx context.Context) {
	httpclient.InitHttp(
		httpclient.NewRequestConfig(constants.APIGetCatsFactConfigKey,
			configs.Get().GetMapD(constants.APIConfig, constants.APIGetCatsFactConfigKey, nil)),
	)

}

func initRouter(ctx context.Context) {
	err := api.Init()
	if err != nil {
		log.Fatal(ctx, " error in initialising router ", err)
	}
}
