package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/shopping/api"
	"github.com/paulantezana/shopping/config"
	"github.com/paulantezana/shopping/migration"
	"os"
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

	// Root router
	//e.GET("/", func(c echo.Context) error {
	//	return c.NoContent(http.StatusOK)
	//})

	// Assets
	static := e.Group("/static")
	static.Static("", "static")

	// Sting API services
	api.PublicApi(e)
	api.ProtectedApi(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = config.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}
