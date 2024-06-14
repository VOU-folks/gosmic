package starters

import (
	"context"

	"osm-api/internal/di"
)

func RegisterServices(ctx context.Context, di di.Container) {
	configService := CreateConfigService()

	di.MustSet("configService", configService)
}
