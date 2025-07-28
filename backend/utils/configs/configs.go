package configs

import (
	"github.com/shekhar8352/PostEaze/constants"
	configs "github.com/sinhashubham95/go-config-client"
)

var client configs.Client

// InitDev is used to initialise configs for dev mode.
func InitDev(directory string, configNames ...string) (err error) {
	client, err = configs.New(configs.Options{
		Provider: configs.FileBased,
		Params: map[string]any{
			constants.ConfigDirectoryKey: directory,
			constants.ConfigNamesKey:     configNames,
			constants.ConfigTypeKey:      constants.ConfigYAML,
		},
	})
	return
}

// InitRelease is used to initialise configs for release mode.
func InitRelease(env, region string, configNames ...string) (err error) {
	client, err = configs.New(configs.Options{
		Provider: configs.AWSAppConfig,
		Params: map[string]any{
			constants.ConfigIDKey:              constants.ApplicationName,
			constants.ConfigRegionKey:          region,
			constants.ConfigEnvKey:             env,
			constants.ConfigAppKey:             constants.ApplicationName,
			constants.ConfigNamesKey:           configNames,
			constants.ConfigTypeKey:            constants.ConfigYAML,
			constants.ConfigCredentialsModeKey: configs.AppConfigSharedCredentialMode,
		},
	})
	return
}

// Get is used to get the client.
func Get() configs.Client {
	return client
}
