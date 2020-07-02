package endpoint

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/shopping/controller"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// PublicApi function public urls
func PublicApi(e *echo.Echo) {
	pb := e.Group("/api/v1/public")

	pb.POST("/user/login", controller.Login)
	//pb.POST("/user/forgot/search", controller.ForgotSearch)
	//pb.POST("/user/forgot/validate", controller.ForgotValidate)
	//pb.POST("/user/forgot/change", controller.ForgotChange)
}

// ProtectedApi function protected urls
func ProtectedApi(e *echo.Echo) {
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(provider.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// Global settings
	ar.GET("/user/by/token", controller.GetUserByToken)
	ar.POST("/user/by/id", controller.GetUserByID)
	ar.POST("/user/paginate", controller.PaginateUser)
	ar.POST("/user/create", controller.CreateUser)
	ar.POST("/user/update", controller.UpdateUser)
	ar.POST("/user/delete", controller.DeleteUser)
}
