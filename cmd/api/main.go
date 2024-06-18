package main

import (
	"context"
	"flag"
	"fmt"
	"gosmic/internal/apps/api"
	. "gosmic/internal/config"
	"gosmic/internal/di"
	"gosmic/internal/handlers"
)

var (
	diContainer di.Container

	app *api.App
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to the config file. Example: api -config /full/path/to/config.yaml")
	flag.Parse()

	config, err := GetConfigService(*configFile)
	if err != nil {
		panic(err)
	}
	apiConfig := config.GetApiConfig()

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

func registerDependencies(di di.Container) {
	di.MustSet("app", app)
}
