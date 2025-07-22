package constants

// flags for application
const (
	ModeKey               = "mode"
	ModeUsage             = "run mode of the application, can be dev or release"
	DevMode               = "dev"
	ReleaseMode           = "release"
	DefaultMode           = DevMode
	PortKey               = "port"
	PortUsage             = "port to run the application on"
	DefaultPort           = 8080
	BaseConfigPathKey     = "base-config-path"
	BaseConfigPathUsage   = "path to configs directory"
	DefaultBaseConfigPath = "resources/configs/local"
	EnvKey                = "ENV"
	AWSRegionKey          = "AWS_REGION"
)
