package endpoint

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/paulantezana/shopping/controller"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// PublicApi function public urls
func PublicApi(e *echo.Echo) {
	pb := e.Group("/api/v1/public")

	pb.POST("/user/login", controller.Login)
	pb.POST("/user/logout", controller.Logout)
	pb.POST("/user/register", controller.RegisterUser)
	pb.POST("/user/forgot/search", controller.ForgotSearch)
	pb.POST("/user/forgot/validate", controller.ForgotValidate)
	pb.POST("/user/forgot/change", controller.ForgotChange)

	pb.POST("/site/card", controller.WebSiteCard)
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
	ar.POST("/user/upload/avatar", controller.UploadAvatarUser)
	ar.POST("/user/changePassword", controller.ChangePasswordUser)
	ar.POST("/user/updateState", controller.UpdateStateUser)
	ar.POST("/user/load/localAuth/by/id", controller.LoadLocalAuthByUserId)
	ar.POST("/user/save/localAuth/by/id", controller.SaveLocalAuthByUserId)

	// User Role
	ar.POST("/user/role/by/id", controller.GetUserRoleByID)
	ar.POST("/user/role/paginate", controller.PaginateUserRole)
	ar.POST("/user/role/all", controller.GetAllUserRole)
	ar.POST("/user/role/create", controller.CreateUserRole)
	ar.POST("/user/role/update", controller.UpdateUserRole)
	ar.POST("/user/role/updateState", controller.UpdateStateUserRole)
	ar.POST("/user/role/appAuthorization/by/userRoleId", controller.GetAppAuthorizationByUserRole)
	ar.POST("/user/role/appAuthorization/update", controller.UpdateUserRoleAppAuthorization)
	ar.GET("/user/role/get/config", controller.GetMenuAdminByUserId)

	// Company
	ar.GET("/company/first", controller.GetFirstCompany)
	ar.POST("/company/by/id", controller.GetCompanyByID)
	ar.POST("/company/update", controller.UpdateCompany)
	ar.POST("/company/uploadLogo", controller.UploadLogoCompany)
	ar.POST("/company/uploadLogoLarge", controller.UploadLogoLargeCompany)
	ar.GET("/operative/get", controller.GetOperative)
	ar.POST("/operative/save", controller.SaveOperative)

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

    // Income type
    ar.POST("/incomeType/by/id", controller.GetIncomeTypeByID)
    ar.POST("/incomeType/paginate", controller.PaginateIncomeType)
    ar.POST("/incomeType/all", controller.GetAllIncomeType)
    ar.POST("/incomeType/create", controller.CreateIncomeType)
    ar.POST("/incomeType/update", controller.UpdateIncomeType)
    ar.POST("/incomeType/updateState", controller.UpdateStateIncomeType)

    // Expense type
    ar.POST("/expenseType/by/id", controller.GetExpenseTypeByID)
    ar.POST("/expenseType/paginate", controller.PaginateExpenseType)
    ar.POST("/expenseType/all", controller.GetAllExpenseType)
    ar.POST("/expenseType/create", controller.CreateExpenseType)
    ar.POST("/expenseType/update", controller.UpdateExpenseType)
    ar.POST("/expenseType/updateState", controller.UpdateStateExpenseType)

	// patter
	ar.POST("/category/by/id", controller.GetCategoryByID)
	ar.POST("/category/all", controller.GetAllCategory)
	ar.POST("/category/paginate", controller.PaginateCategory)
	ar.POST("/category/create", controller.CreateCategory)
	ar.POST("/category/update", controller.UpdateCategory)
	ar.POST("/category/updateState", controller.UpdateStateCategory)

	// Product
	ar.POST("/product/by/id", controller.GetProductByID)
	ar.POST("/product/search", controller.GetProductSearch)
	ar.POST("/product/paginate", controller.PaginateProduct)
	ar.POST("/product/create", controller.CreateProduct)
	ar.POST("/product/update", controller.UpdateProduct)
	ar.POST("/product/updateState", controller.UpdateStateProduct)
	ar.POST("/product/seeker/by/code", controller.GetProductSeekerByCode)
	ar.POST("/product/seeker/paginate/search", controller.PaginateProductSeekerSearch)
	ar.POST("/product/import", controller.ImportProduct)

	// provider
	ar.POST("/provider/by/id", controller.GetProviderByID)
	ar.POST("/provider/all", controller.GetAllProvider)
	ar.POST("/provider/search", controller.GetSearchProvider)
	ar.POST("/provider/paginate", controller.PaginateProvider)
	ar.POST("/provider/create", controller.CreateProvider)
	ar.POST("/provider/update", controller.UpdateProvider)
	ar.POST("/provider/updateState", controller.UpdateStateProvider)

	// Purchase
	ar.POST("/purchase/item/by/purchase/id", controller.GetPurchaseItemByPurchaseID)
	ar.POST("/purchase/newPurchase", controller.NewPurchase)
	ar.POST("/purchase/cancel", controller.CancelPurchase)
	ar.POST("/purchase/paginate", controller.PaginatePurchase)

	// Customer
	ar.POST("/customer/by/id", controller.GetCustomerByID)
	ar.POST("/customer/all", controller.GetAllCustomer)
	ar.POST("/customer/search", controller.GetSearchCustomer)
	ar.POST("/customer/paginate", controller.PaginateCustomer)
	ar.POST("/customer/create", controller.CreateCustomer)
	ar.POST("/customer/update", controller.UpdateCustomer)
	ar.POST("/customer/updateState", controller.UpdateStateCustomer)

	// Sale
	ar.POST("/sale/item/by/sale/id", controller.GetSaleItemBySaleID)
	ar.POST("/sale/newSale", controller.NewSale)
	ar.POST("/sale/cancel", controller.CancelSale)
	ar.POST("/sale/paginate", controller.PaginateSale)

	// Kardex
	ar.POST("/kardex/paginate", controller.PaginateKardex)

	// Utils all
	ar.GET("/util/additionalLegendType/all", controller.GetAllUtilAdditionalLegendType)
	ar.GET("/util/affectationIgvType/all", controller.GetAllUtilAffectationIgvType)
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

	// Custom
	ar.POST("/util/documentType/all/sale", controller.GetAllUtilDocumentTypeSale)

	// Utils Search
	ar.POST("/util/geographicalLocation/search", controller.GetSearchUtilGeographicalLocation)
	ar.POST("/util/productType/search", controller.GetAllUtilGeographicalLocation)
	ar.POST("/util/productType/search", controller.GetAllUtilGeographicalLocation)
	ar.POST("/util/queryPeru", controller.QueryDocument)
}
