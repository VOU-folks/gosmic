package api

import (
	"context"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"net/http"
	"osm-api/internal/di"
)

func RegisterRoutes(ctx context.Context, di di.Container, engine *echo.Echo) {
	appInfo := AppInfo{
		Name:    ctx.Value("AppName").(string),
		Version: ctx.Value("AppVersion").(string),
	}

	engine.GET("/", func(c echo.Context) error {
		return c.JSONPretty(
			http.StatusOK,
			appInfo,
			"  ",
		)
	})
	engine.GET("/_metrics", echoprometheus.NewHandler())
	engine.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}
