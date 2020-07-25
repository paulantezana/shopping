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

	// Company Local
	ar.POST("/company/local/by/id", controller.GetCompanyLocalByID)
	ar.POST("/company/local/paginate", controller.PaginateCompanyLocal)
	ar.POST("/company/local/create", controller.CreateCompanyLocal)
	ar.POST("/company/local/update", controller.UpdateCompanyLocal)
	ar.POST("/company/local/delete", controller.DeleteCompanyLocal)

	// Utils all
	ar.GET("/util/additionalLegendType/all", controller.GetAllUtilAdditionalLegendType)
	ar.GET("/util/catAffectationIgvType/all", controller.GetAllUtilCatAffectationIgvType)
	ar.GET("/util/creditDebitType/all", controller.GetAllUtilCreditDebitType)
	ar.GET("/util/currencyType/all", controller.GetAllUtilCurrencyType)
	ar.GET("/util/documentType/all", controller.GetAllUtilDocumentType)
	ar.GET("/util/geographicalLocation/all", controller.GetAllUtilGeographicalLocation)
	ar.GET("/util/identityDocumentType/all", controller.GetAllUtilIdentityDocumentType)
	ar.GET("/util/operationType/all", controller.GetAllUtilOperationType)
	ar.GET("/util/perceptionType/all", controller.GetAllUtilPerceptionType)
	ar.GET("/util/productType/all", controller.GetAllUtilProductType)
	ar.GET("/util/subjectDetractionType/all", controller.GetAllUtilSubjectDetractionType)
	ar.GET("/util/systemIscType/all", controller.GetAllUtilSystemIscType)
	ar.GET("/util/transferReasonType/all", controller.GetAllUtilTransferReasonType)
	ar.GET("/util/transportModeType/all", controller.GetAllUtilTransportModeType)
	ar.GET("/util/tributeType/all", controller.GetAllUtilTributeType)
	ar.GET("/util/unitMeasureType/all", controller.GetAllUtilUnitMeasureType)

	// Utils Search
	ar.POST("/util/geographicalLocation/search", controller.GetSearchUtilGeographicalLocation)
	ar.POST("/util/productType/search", controller.GetAllUtilGeographicalLocation)
	ar.POST("/util/productType/search", controller.GetAllUtilGeographicalLocation)
}
