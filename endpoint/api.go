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
	ar.POST("/user/changePassword", controller.ChangePasswordUser)
	ar.POST("/user/updateState", controller.UpdateStateUser)

    // User Role
    ar.POST("/user/role/by/id", controller.GetUserRoleByID)
    ar.POST("/user/role/paginate", controller.PaginateUserRole)
    ar.POST("/user/role/all", controller.GetAllUserRole)
    ar.POST("/user/role/create", controller.CreateUserRole)
    ar.POST("/user/role/update", controller.UpdateUserRole)
    ar.POST("/user/role/updateState", controller.UpdateStateUserRole)
    ar.POST("/user/role/appAuthorization/by/userRoleId", controller.GetAppAuthorizationByUserRole)
    ar.POST("/user/role/appAuthorization/update", controller.UpdateUserRoleAppAuthorization)

	// Company
	ar.GET("/company/first", controller.GetFirstCompany)
	ar.POST("/company/by/id", controller.GetCompanyByID)
	ar.POST("/company/update", controller.UpdateCompany)
	ar.POST("/company/uploadLogo", controller.UploadLogoCompany)
	ar.POST("/company/uploadLogoLarge", controller.UploadLogoLargeCompany)

	// Company Local
	ar.POST("/company/local/by/id", controller.GetCompanyLocalByID)
	ar.POST("/company/local/paginate", controller.PaginateCompanyLocal)
	ar.POST("/company/local/all", controller.GetAllCompanyLocal)
	ar.POST("/company/local/create", controller.CreateCompanyLocal)
	ar.POST("/company/local/update", controller.UpdateCompanyLocal)
	ar.POST("/company/local/updateState", controller.UpdateStateCompanyLocal)

	// Company WareHouse
	ar.POST("/company/warehouse/by/id", controller.GetCompanyWareHouseByID)
	ar.POST("/company/warehouse/paginate", controller.PaginateCompanyWareHouse)
	ar.POST("/company/warehouse/create", controller.CreateCompanyWareHouse)
	ar.POST("/company/warehouse/update", controller.UpdateCompanyWareHouse)
	ar.POST("/company/warehouse/updateState", controller.UpdateStateCompanyWareHouse)

	// Company SalePoint
	ar.POST("/company/salePoint/by/id", controller.GetCompanySalePointByID)
	ar.POST("/company/salePoint/paginate", controller.PaginateCompanySalePoint)
	ar.POST("/company/salePoint/create", controller.CreateCompanySalePoint)
	ar.POST("/company/salePoint/update", controller.UpdateCompanySalePoint)
	ar.POST("/company/salePoint/updateState", controller.UpdateStateCompanySalePoint)

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
