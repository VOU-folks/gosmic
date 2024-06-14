package config

import "osm-api/internal/services"

func GetConfigService(configFilePath string) (*services.ConfigService, error) {
	var err error
	var configService *services.ConfigService

	defaultConfigFilePath := "config.yaml"
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}
	configService = services.NewConfigService(configFilePath)
	err = configService.Init()
	if err != nil {
		return nil, err
	}

	return configService, nil
}
