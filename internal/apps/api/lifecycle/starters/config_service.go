package starters

import (
	"gosmic/internal/services"
)

func CreateConfigService() *services.ConfigService {
	service := services.NewConfigService("config.yaml")

	err := service.Init()
	if err != nil {
		panic(err)
	}

	return service
}
