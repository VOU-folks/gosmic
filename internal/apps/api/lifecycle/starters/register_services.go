package starters

import (
	"context"

	"gosmic/internal/di"
)

func RegisterServices(ctx context.Context, di di.Container) {
	configService := CreateConfigService()

	di.MustSet("configService", configService)
}
