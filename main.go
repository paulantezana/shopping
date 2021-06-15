package main

import (
	"github.com/paulantezana/shopping/endpoint"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/paulantezana/shopping/migration"
	"github.com/paulantezana/shopping/provider"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Migration database
	migration.Migrate()

	// Configuration cor
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// Assets
	static := e.Group("/static")
	static.Static("", "static")

	// Root router success
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Sting API services
	endpoint.PublicApi(e)
	endpoint.ProtectedApi(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = provider.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}
