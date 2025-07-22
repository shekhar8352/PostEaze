package constants

// config init constants
const (
	ConfigDirectoryKey       = "configsDirectory"
	ConfigNamesKey           = "configNames"
	ConfigTypeKey            = "configType"
	ConfigYAML               = "yaml"
	ConfigIDKey              = "id"
	ConfigRegionKey          = "region"
	ConfigCredentialsModeKey = "credentialsMode"
	ConfigEnvKey             = "env"
	ConfigAppKey             = "app"
)

// config names
const (
	ApplicationConfig = "application"
	LoggerConfig      = "logger"
	DatabaseConfig    = "database"
	APIConfig         = "api"
)

// config constants
const (
	LoggerLevelConfigKey                            = "level"
	LoggerParamsConfigKey                           = "params"
	DatabaseDriverNameConfigKey                     = "driverName"
	DatabaseURLConfigKey                            = "url"
	DatabaseMaxOpenConnectionsConfigKey             = "maxOpenConnections"
	DatabaseMaxIdleConnectionsConfigKey             = "maxIdleConnections"
	DatabaseMaxConnectionLifetimeInSecondsConfigKey = "maxConnectionLifetimeInSeconds"
	DatabaseMaxConnectionIdleTimeInSecondsConfigKey = "maxConnectionIdleTimeInSeconds"
)
