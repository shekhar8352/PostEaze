package flags

import (
	"flag"
	"os"

	"github.com/shekhar8352/PostEaze/constants"
)

var (
	mode           = flag.String(constants.ModeKey, constants.DefaultMode, constants.ModeUsage)
	port           = flag.Int(constants.PortKey, constants.DefaultPort, constants.PortUsage)
	baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.DefaultBaseConfigPath,
		constants.BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

// Mode is the application running mode.
func Mode() string {
	return *mode
}

// Port is used to get the port.
func Port() int {
	return *port
}

// BaseConfigPath is used to get the config path for file-based configs.
func BaseConfigPath() string {
	return *baseConfigPath
}

// Env is the environment.
func Env() string {
	return os.Getenv(constants.EnvKey)
}

// AWSRegion is the aws region.
func AWSRegion() string {
	return os.Getenv(constants.AWSRegionKey)
}
