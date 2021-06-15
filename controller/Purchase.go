package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

type newPurchase struct {
	DateOfIssue       time.Time `json:"date_of_issue"`
	DateOfPurchase    time.Time `json:"date_of_purchase"`
	Serie             string    `json:"serie"`
	Number            string    `json:"number"`
	Guide             string    `json:"guide" `
	Observation       string    `json:"observation"`
	CancelObservation string    `json:"cancel_observation"`
	TotalUnaffected   float64   `json:"total_unaffected"`
	TotalTaxed        float64   `json:"total_taxed"`
	TotalIgv          float64   `json:"total_igv"`
	Total             float64   `json:"total"`

	ProviderId         uint `json:"provider_id"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`
	UtilCurrencyTypeId uint `json:"util_currency_type_id"`
	UtilDocumentTypeId uint `json:"util_document_type_id"`
	CompanyId          uint `json:"company_id"`

	Item []models.PurchaseItem `json:"item"`
	Pay  struct{}              `json:"pay"`
}

// GetPurchaseItemByPurchaseID function create new purchase
func GetPurchaseItemByPurchaseID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	purchaseObj := models.Purchase{}
	if err := c.Bind(&purchaseObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "purchase_new_purchase"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Find PurchaseItem
	purchaseItem := make([]models.PurchaseItem, 0)
	if err := DB.Where("purchase_id = ?", purchaseObj.ID).Find(&purchaseItem).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    purchaseItem,
	})
}

// NewPurchase function create new purchase
func NewPurchase(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	purchaseObj := newPurchase{}
	if err := c.Bind(&purchaseObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	if validateResponse := validateNewPurchase(purchaseObj); validateResponse.Success == false {
		return c.JSON(http.StatusOK, validateResponse)
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// start transaction
	TX := DB.Begin()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "purchase_new_purchase"); err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: "unauthorized"})
	}

	// Insert purchase in database
	purchase := models.Purchase{}
	purchase.DateOfIssue = time.Now()
	purchase.DateOfPurchase = purchaseObj.DateOfPurchase
	purchase.Guide = purchaseObj.Guide
	purchase.Observation = purchaseObj.Guide
	purchase.TotalUnaffected = purchaseObj.TotalUnaffected
	purchase.TotalTaxed = purchaseObj.TotalTaxed
	purchase.TotalIgv = purchaseObj.TotalIgv
	purchase.Total = purchaseObj.Total
	purchase.ProviderId = purchaseObj.ProviderId
	purchase.CompanyWareHouseId = purchaseObj.CompanyWareHouseId
	purchase.UtilCurrencyTypeId = purchaseObj.UtilCurrencyTypeId
	purchase.UtilDocumentTypeId = purchaseObj.UtilDocumentTypeId
	purchase.CreatedUserId = currentUser.ID
	purchase.CompanyId = currentUser.CompanyId
	purchase.UserId = currentUser.ID
	if err := TX.Create(&purchase).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	wareHouse := models.CompanyWareHouse{ID: purchaseObj.CompanyWareHouseId}
	if TX.First(&wareHouse).RowsAffected == 0 {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro del almacen con id %d", wareHouse.ID),
		})
	}

	for _, item := range purchaseObj.Item {
		item.PurchaseId = purchase.ID
		if err := TX.Create(&item).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}

		// Update Kardex
		kardexAux := models.Kardex{}
		TX.Where("product_id = ? AND company_ware_house_id = ? AND is_last = true", item.ProductId, wareHouse.ID).First(&kardexAux)
		TX.Model(&models.Kardex{}).Where("product_id = ?", item.ProductId).Where("company_ware_house_id = ?", wareHouse.ID).Update("is_last", false)

		kardex := models.Kardex{}
		kardex.DateOfIssue = time.Now()
		kardex.Quantity = item.Quantity
		kardex.UnitPrice = item.UnitPrice
		kardex.Total = item.Total
		kardex.Origin = "Compra"
		kardex.Destination = wareHouse.Description
		kardex.Description = item.Description
		kardex.DocumentDescription = "Compra"
		kardex.UserId = currentUser.ID
		kardex.ProductId = item.ProductId
		kardex.CompanyWareHouseId = wareHouse.ID
		kardex.Stock = kardexAux.Stock + item.Quantity
		kardex.IsLast = true
		kardex.IsIncome = true
		if err := TX.Create(&kardex).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    purchase.ID,
		Message: fmt.Sprintf("La compra se realizó exitosamente"),
	})
}

// CancelPurchase function create new purchase
func CancelPurchase(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	purchaseObj := models.Purchase{}
	if err := c.Bind(&purchaseObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// start transaction
	TX := DB.Begin()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "purchase_new_purchase"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	purchase := models.Purchase{}
	if err := DB.First(&purchase, purchaseObj.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query Ware House
	wareHouse := models.CompanyWareHouse{ID: purchase.CompanyWareHouseId}
	if DB.First(&wareHouse).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro del almacen con id %d", wareHouse.ID),
		})
	}

	// Find PurchaseItem
	purchaseItem := make([]models.PurchaseItem, 0)
	if err := DB.Where("purchase_id = ?", purchase.ID).Find(&purchaseItem).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	for _, item := range purchaseItem {
		// Update Kardex
		kardexAux := models.Kardex{}
		TX.Where("product_id = ? AND company_ware_house_id = ? AND is_last = true", item.ProductId, purchase.CompanyWareHouseId).First(&kardexAux)
		TX.Model(&models.Kardex{}).Where("product_id = ?", item.ProductId).Where("company_ware_house_id = ?", wareHouse.ID).Update("is_last", false)

		kardex := models.Kardex{}
		kardex.DateOfIssue = time.Now()
		kardex.Quantity = item.Quantity
		kardex.UnitPrice = item.UnitPrice
		kardex.Total = item.Total
		kardex.Origin = wareHouse.Description
		kardex.Destination = "Anular compra"
		kardex.Description = item.Description
		kardex.DocumentDescription = "Anular compra"
		kardex.UserId = currentUser.ID
		kardex.ProductId = item.ProductId
		kardex.CompanyWareHouseId = purchase.CompanyWareHouseId
		kardex.Stock = kardexAux.Stock - item.Quantity
		kardex.IsLast = true
		kardex.IsIncome = false
		if err := TX.Create(&kardex).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Update purchase
	if err := TX.Exec("UPDATE purchases SET state = false, cancel_observation = ? WHERE id = ?", purchaseObj.CancelObservation, purchase.ID).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    purchase.ID,
		Message: fmt.Sprintf("La compra se anuló exitosamente"),
	})
}

type purchasePage struct {
	ID                uint      `json:"id"`
	DateOfIssue       time.Time `json:"date_of_issue"`
	DateOfPurchase    time.Time `json:"date_of_purchase"`
	Serie             string    `json:"serie"`
	Number            string    `json:"number"`
	Observation       string    `json:"observation"`
	CancelObservation string    `json:"cancel_observation"`
	TotalUnaffected   float64   `json:"total_unaffected"`
	TotalTaxed        float64   `json:"total_taxed"`
	TotalIgv          float64   `json:"total_igv"`
	Total             float64   `json:"total"`

	ProviderId         uint `json:"provider_id"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`
	UtilCurrencyTypeId uint `json:"util_currency_type_id"`
	UtilDocumentTypeId uint `json:"util_document_type_id"`
	CompanyId          uint `json:"company_id"`
	UserId             uint `json:"user_id"`
	State              bool `json:"state"`

	DocumentName   string `json:"document_name"`
	ProviderName   string `json:"provider_name"`
	UserName       string `json:"user_name"`
	CurrencySymbol string `json:"currency_symbol"`
}

// PaginateProduct function get all products
func PaginatePurchase(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "purchase_report"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	purchases := make([]purchasePage, 0)

	// Find purchases
	if err := DB.Table("purchases").Select("purchases.*, providers.social_reason as provider_name, util_document_types.description as document_name, users.user_name as user_name, util_currency_types.symbol as currency_symbol").
		Joins("INNER JOIN providers ON purchases.provider_id = providers.id").
		Joins("INNER JOIN util_document_types ON purchases.util_document_type_id = util_document_types.id").
		Joins("INNER JOIN users ON purchases.user_id = users.id").
		Joins("INNER JOIN util_currency_types ON purchases.util_currency_type_id = util_currency_types.id").
		Where("purchases.company_ware_house_id = ? AND  (purchases.date_of_issue BETWEEN ? AND ?)", request.WareHouseId, request.StartDate, request.EndDate).
		Order("purchases.id desc").Offset(offset).Limit(request.PageSize).Scan(&purchases).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     purchases,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// validateNewPurchase function validate
func validateNewPurchase(p newPurchase) utilities.Response {
	res := utilities.Response{}
	res.Success = true

	if p.ProviderId == 0 {
		res.Message = "Falta especificar el proveedor | "
		res.Success = false
	}
	if p.UtilDocumentTypeId == 0 {
		res.Message += "Falta especificar el tipo de comprobante | "
		res.Success = false
	}
	if p.UtilCurrencyTypeId == 0 {
		res.Message += "Falta especificar la moneda | "
		res.Success = false
	}
	if len(p.Item) == 0 {
		res.Message += "Ingrese al menos un artículo | "
		res.Success = false
	}
	if p.Total == 0 {
		res.Message += "El total no puede ser igual a cero | "
		res.Success = false
	}

	return res
}
