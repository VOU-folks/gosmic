package api

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"osm-api/internal/apps/api/lifecycle/starters"
	"osm-api/internal/apps/api/lifecycle/stoppers"

	"golang.org/x/crypto/acme/autocert"

	"osm-api/internal/di"
)

type App struct {
	ctx    context.Context
	di     di.Container
	engine *echo.Echo
}

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewApp(ctx context.Context, di di.Container) *App {
	engine := echo.New()
	engine.HideBanner = true

	sslEnabled := ctx.Value("SslEnabled").(bool)
	if sslEnabled {
		engine.AutoTLSManager.HostPolicy = autocert.HostWhitelist(ctx.Value("TlsManagerHostWhitelist").(string))
		engine.AutoTLSManager.Cache = autocert.DirCache(ctx.Value("TlsManagerCacheDir").(string))
		engine.AutoTLSManager.Prompt = autocert.AcceptTOS
		engine.AutoTLSManager.Email = ctx.Value("TlsManagerEmail").(string)
	}

	starters.RegisterServices(ctx, di)

	RegisterMiddlewares(ctx, di, engine)
	RegisterRoutes(ctx, di, engine)

	return &App{
		ctx:    ctx,
		di:     di,
		engine: engine,
	}
}

func (app *App) SetListener(listener net.Listener) {
	app.engine.Listener = listener
}

func (app *App) SetServer(server *http.Server) {
	app.engine.Server = server
}

func (app *App) Start() error {
	listenAt := app.ctx.Value("ListenAt").(string)
	sslEnabled := app.ctx.Value("SslEnabled").(bool)

	err := app.startListener(sslEnabled, listenAt)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (app *App) startListener(sslEnabled bool, listenAt string) error {
	if sslEnabled {
		return app.engine.StartAutoTLS(listenAt)
	}
	return app.engine.Start(listenAt)
}

func (app *App) Stop() error {
	stoppers.StopServices(app.di)

	return app.engine.Shutdown(app.ctx)
}
