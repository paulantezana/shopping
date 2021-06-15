package controller

import (
	"fmt"
    "gorm.io/gorm"
    "net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

type newSale struct {
	ChangeType        string `json:"change_type"`
	VehiclePlate      string `json:"vehicle_plate"`
	Term              string `json:"term"`
	PdfFormat         string `json:"pdf_format"`
	Guide             string `json:"guide" `
	Observation       string `json:"observation"`
	CancelObservation string `json:"cancel_observation"`

	DocumentNumber             string `json:"document_number"`
	SocialReason               string `json:"social_reason"`
	FiscalAddress              string `json:"fiscal_address"`
	Email                      string `json:"email"`
	Phone                      string `json:"phone"`
	EmailSend                  bool   `json:"email_send"`
	PhoneSend                  bool   `json:"phone_send"`
	CustomerId                 uint   `json:"customer_id"`
	UtilIdentityDocumentTypeId uint   `json:"util_identity_document_type_id"`

	TotalFree          float64 `json:"total_free"`
	TotalExportation   float64 `json:"total_exportation"`
	TotalDiscount      float64 `json:"total_discount"`
	TotalExonerated    float64 `json:"total_exonerated"`
	TotalUnaffected    float64 `json:"total_unaffected"`
	TotalTaxed         float64 `json:"total_taxed"`
	TotalIgv           float64 `json:"total_igv"`
	TotalCharge        float64 `json:"total_charge"`
	TotalValue         float64 `json:"total_value"`
	TotalPlasticBagTax float64 `json:"total_plastic_bag_tax"`
	Total              float64 `json:"total"`

	UtilDocumentTypeId  uint `json:"util_document_type_id"`
	UtilCurrencyTypeId  uint `json:"util_currency_type_id"`
	UtilOperationTypeId uint `json:"util_operation_type_id"`
	CompanyWareHouseId  uint `json:"company_ware_house_id"`
	UserId              uint `json:"user_id" gorm:"not null"`

	State bool `json:"state" gorm:"default: true"`

	Item []models.SaleItem `json:"item"`
	Pay  struct{}          `json:"pay"`
}

// GetSaleItemBySaleID function create new sale
func GetSaleItemBySaleID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	saleObj := models.Sale{}
	if err := c.Bind(&saleObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_new_sale"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Find SaleItem
	saleItem := make([]models.SaleItem, 0)
	if err := DB.Where("sale_id = ?", saleObj.ID).Find(&saleItem).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    saleItem,
	})
}

// NewSale function create new sale
func NewSale(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	saleObj := newSale{}
	if err := c.Bind(&saleObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	if validateResponse := validateNewSale(saleObj); validateResponse.Success == false {
		return c.JSON(http.StatusOK, validateResponse)
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// start transaction
	TX := DB.Begin()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_new_sale"); err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: "unauthorized"})
	}

	// Insert sale in database
	sale := models.Sale{}
	sale.DateOfIssue = time.Now()
	sale.TimeOfIssue = time.Now()
	sale.DateOfDue = time.Now()
	sale.Serie = ""
	sale.Number = 0
	sale.ChangeType = ""
	sale.VehiclePlate = saleObj.VehiclePlate
	sale.Term = saleObj.Term
	sale.PdfFormat = saleObj.PdfFormat
	sale.Guide = saleObj.Guide
	sale.Observation = saleObj.Observation
	sale.DocumentNumber = saleObj.DocumentNumber
	sale.SocialReason = saleObj.SocialReason
	sale.FiscalAddress = saleObj.FiscalAddress
	sale.Email = saleObj.Email
	sale.Phone = saleObj.Phone
	sale.EmailSend = saleObj.EmailSend
	sale.PhoneSend = saleObj.PhoneSend
	sale.CustomerId = saleObj.CustomerId
	sale.UtilIdentityDocumentTypeId = saleObj.UtilIdentityDocumentTypeId
	sale.TotalFree = saleObj.TotalFree
	sale.TotalExportation = saleObj.TotalExportation
	sale.TotalDiscount = saleObj.TotalDiscount
	sale.TotalExonerated = saleObj.TotalExonerated
	sale.TotalUnaffected = saleObj.TotalUnaffected
	sale.TotalTaxed = saleObj.TotalTaxed
	sale.TotalIgv = saleObj.TotalIgv
	sale.TotalCharge = saleObj.TotalCharge
	sale.TotalValue = saleObj.TotalValue
	sale.TotalPlasticBagTax = saleObj.TotalPlasticBagTax
	sale.Total = saleObj.Total
	sale.UtilDocumentTypeId = saleObj.UtilCurrencyTypeId
	sale.UtilCurrencyTypeId = saleObj.UtilCurrencyTypeId
	sale.UtilOperationTypeId = saleObj.UtilOperationTypeId
	sale.CompanyId = currentUser.CompanyId
	sale.CompanyWareHouseId = saleObj.CompanyWareHouseId
	sale.UserId = currentUser.ID

	if err := TX.Create(&sale).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	wareHouse := models.CompanyWareHouse{ID: saleObj.CompanyWareHouseId}
	if TX.First(&wareHouse).RowsAffected == 0 {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro del almacen con id %d", wareHouse.ID),
		})
	}

	for _, item := range saleObj.Item {
		item.SaleId = sale.ID
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
		kardex.Origin = "Venta"
		kardex.Destination = wareHouse.Description
		kardex.Description = item.Description
		kardex.DocumentDescription = "Venta"
		kardex.UserId = currentUser.ID
		kardex.ProductId = item.ProductId
		kardex.CompanyWareHouseId = wareHouse.ID
		kardex.Stock = kardexAux.Stock - item.Quantity
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
		Data:    sale.ID,
		Message: "La venta se realizó exitosamente",
	})
}

// CancelSale function create new sale
func CancelSale(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	saleObj := models.Sale{}
	if err := c.Bind(&saleObj); err != nil {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_new_sale"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	sale := models.Sale{}
	if err := DB.First(&sale, saleObj.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query Ware House
	wareHouse := models.CompanyWareHouse{ID: sale.CompanyWareHouseId}
	if DB.First(&wareHouse).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro del almacen con id %d", wareHouse.ID),
		})
	}

	// Find SaleItem
	saleItem := make([]models.SaleItem, 0)
	if err := DB.Where("sale_id = ?", sale.ID).Find(&saleItem).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	for _, item := range saleItem {
		// Update Kardex
		kardexAux := models.Kardex{}
		TX.Where("product_id = ? AND company_ware_house_id = ? AND is_last = true", item.ProductId, sale.CompanyWareHouseId).First(&kardexAux)
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
		kardex.CompanyWareHouseId = sale.CompanyWareHouseId
		kardex.Stock = kardexAux.Stock + item.Quantity
		kardex.IsLast = true
		kardex.IsIncome = false
		if err := TX.Create(&kardex).Error; err != nil {
			TX.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Update sale
	if err := TX.Exec("UPDATE sales SET state = false, cancel_observation = ? WHERE id = ?", saleObj.CancelObservation, sale.ID).Error; err != nil {
		TX.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Commit transaction
	TX.Commit()

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    sale.ID,
		Message: fmt.Sprintf("La compra se anuló exitosamente"),
	})
}

type salePage struct {
	ID                uint      `json:"id"`
	DateOfIssue       time.Time `json:"date_of_issue"`
	DateOfSale        time.Time `json:"date_of_sale"`
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

// PaginateSale function get all products
func PaginateSale(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_report"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	sales := make([]salePage, 0)

	// Find sales
	if err := DB.Table("sales").Select("sales.*, util_document_types.description as document_name, users.user_name as user_name, util_currency_types.symbol as currency_symbol").
		Joins("INNER JOIN util_document_types ON sales.util_document_type_id = util_document_types.id").
		Joins("INNER JOIN users ON sales.user_id = users.id").
		Joins("INNER JOIN util_currency_types ON sales.util_currency_type_id = util_currency_types.id").
		Where("sales.company_ware_house_id = ? AND  (sales.date_of_issue BETWEEN ? AND ?)", request.WareHouseId, request.StartDate, request.EndDate).
		Order("sales.id desc").Offset(offset).Limit(request.PageSize).Scan(&sales).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     sales,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// validateNewSale function validate
func validateNewSale(p newSale) utilities.Response {
	res := utilities.Response{}
	res.Success = true

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



type tokensCompanySalePointNumber struct {
    Number          uint
}

type tokensCompanySalePointSerie struct {
    Serie           string
    NubefactUrlPath string
    NubefactToken   string
}

type tokensCompanySalePoint struct {
    Serie           string
    Number          uint
    NubefactUrlPath string
    NubefactToken   string
}

func getTokensCompanyLocal(DB gorm.DB, companySalePointId uint, documentTypeId uint, contingency bool) (error, tokensCompanySalePoint) {
    tokensCompanySalePoints := tokensCompanySalePoint{}

    tokensCompanySalePointSeries := tokensCompanySalePointSerie{}
    if err := DB.Table("company_sale_point_series").Select("serie").
        Where("id = ?", companySalePointId).
        Where("util_document_type_id = ?", documentTypeId).
        Where("contingency = ?", contingency).
        Scan(&tokensCompanySalePointSeries).Error; err != nil {
        return err, tokensCompanySalePoints
    }

    tokensCompanySalePointNumbers := tokensCompanySalePointNumber{}
    if err := DB.Table("sales").Select("serie").
        Where("company_local_id = ?", companySalePointId).
        Where("util_document_type_id = ?", documentTypeId).
        Where("contingency = ?", contingency).
        Scan(&tokensCompanySalePointNumbers).Error; err != nil {
        return err, tokensCompanySalePoints
    }

    return nil, tokensCompanySalePoints
}
