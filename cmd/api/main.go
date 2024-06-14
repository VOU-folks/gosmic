package main

import (
	"context"
	"flag"
	"fmt"
	"osm-api/internal/apps/api"
	"osm-api/internal/di"
	"osm-api/internal/handlers"
	"osm-api/internal/services"
	"osm-api/internal/structs"
)

var (
	diContainer di.Container

	app *api.App
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to the config file. Example: api -config /full/path/to/config.yaml")
	flag.Parse()

	config, err := getConfig(*configFile)
	if err != nil {
		panic(err)
	}
	apiConfig := config.Apps.Api

	ctx := context.Background()
	ctx = context.WithValue(ctx, "AppName", apiConfig.Name)
	ctx = context.WithValue(ctx, "AppVersion", apiConfig.Version)
	ctx = context.WithValue(ctx, "TlsManagerEmail", apiConfig.Ssl.Manager.Email)
	ctx = context.WithValue(ctx, "TlsManagerHostWhitelist", apiConfig.Ssl.Manager.HostWhitelist)
	ctx = context.WithValue(ctx, "TlsManagerCacheDir", apiConfig.Ssl.Manager.CacheDir)
	ctx = context.WithValue(ctx, "SslEnabled", apiConfig.Ssl.Enabled)
	ctx = context.WithValue(ctx, "ListenAt", fmt.Sprintf("%v:%d", apiConfig.Listener.Host, apiConfig.Listener.Port))

	diContainer = di.NewDI(ctx)
	app = api.NewApp(ctx, diContainer)

	registerDependencies(diContainer)

	handlers.StartLifecycle(app)
}

func getConfig(configFilePath string) (structs.Config, error) {
	var err error
	var configService *services.ConfigService

	defaultConfigFilePath := "config.yaml"
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}
	configService = services.NewConfigService(configFilePath)
	err = configService.Init()
	if err != nil {
		return structs.Config{}, err
	}

	return configService.GetConfig(), nil
}

func registerDependencies(di di.Container) {
	di.MustSet("app", app)
}
