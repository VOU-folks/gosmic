package api

import (
	"context"
	"fmt"
	"io"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gosmic/internal/di"
)

func RegisterMiddlewares(ctx context.Context, di di.Container, engine *echo.Echo) {
	attachRequestIdMiddleware(engine)

	engine.Use(middleware.Recover())
	engine.Use(echoprometheus.NewMiddleware("api"))

	//attachHeaderDumper(engine)
	//attachBodyDumper(engine)

	engine.Use(
		middleware.LoggerWithConfig(
			middleware.LoggerConfig{
				Format: `${time_rfc3339} | ${remote_ip} | ${method} ${uri} | ${status} | ${latency_human} | ${bytes_in} / ${bytes_out} (in/out) | ${id}` + "\n",
			},
		),
	)
}

func attachRequestIdMiddleware(engine *echo.Echo) {
	config := middleware.RequestIDConfig{
		TargetHeader: "Request-Id",
		RequestIDHandler: func(c echo.Context, requestId string) {
			c.Set("RequestId", requestId)
			c.Request().Header.Set("X-Request-Id", requestId)
		},
	}
	engine.Pre(middleware.RequestIDWithConfig(config))
}

func attachHeaderDumper(engine *echo.Echo) {
	engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Logging headers for request:", c.Request())
			for k, v := range c.Request().Header {
				fmt.Println(k, v)
			}
			return next(c)
		}
	})
}

func attachBodyDumper(engine *echo.Echo) {
	engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Logging body for request:", c.Request().Header.Get("X-Request-Id"))
			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				fmt.Println("Failed to read request body:", err)
			}
			fmt.Println(string(body))
			return next(c)
		}
	})
}
