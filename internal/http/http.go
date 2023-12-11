package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type HttpServer struct {
	App  *echo.Echo
	port string
}

func Bootstrap(port string) *HttpServer {
	// create the fiber app
	app := echo.New()

	// add middleware
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	// add health check
	app.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy!")
	})

	// add docs
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	return &HttpServer{
		App:  app,
		port: port,
	}
}

func (h *HttpServer) Start() {
	h.App.Start("0.0.0.0:" + h.port)
}

func (h *HttpServer) Cleanup() {
	h.App.Shutdown(context.Background())
}
