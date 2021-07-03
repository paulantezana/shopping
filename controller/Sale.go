package controller

import (
    "crypto/sha256"
    "errors"
    "fmt"
    "github.com/jung-kurt/gofpdf"
    "gorm.io/gorm"
    "net/http"
    "strings"
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
	CompanySalePointId  uint `json:"company_sale_point_id"`
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

type newSaleResponse struct {
	InvoicePath      string             `json:"invoice_path"`
	NubefactResponse utilities.Response `json:"nubefact_response"`
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

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_new_sale"); err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: "unauthorized"})
	}

	err, salePointTokens := getTokensCompanySalePoint(DB, saleObj.CompanySalePointId, saleObj.UtilDocumentTypeId, false)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	sale := models.Sale{}
	sale.DateOfIssue = time.Now()
	sale.TimeOfIssue = time.Now()
	sale.DateOfDue = time.Now()

	errTransaction := DB.Transaction(func(TX *gorm.DB) error {
		// Save Customer
		if saleObj.CustomerId == 0 {
			customer := models.Customer{DocumentNumber: saleObj.DocumentNumber}
			if TX.Where("document_number = ?", customer.DocumentNumber).Where("company_id = ?", currentUser.CompanyId).First(&customer).RowsAffected == 0 {
				customer.DocumentNumber = saleObj.DocumentNumber
				customer.SocialReason = saleObj.SocialReason
				customer.UtilIdentityDocumentTypeId = saleObj.UtilIdentityDocumentTypeId
				customer.CompanyId = currentUser.CompanyId
				customer.Address = saleObj.FiscalAddress
				customer.Email = saleObj.Email
				customer.Phone = saleObj.Phone
				if err := TX.Create(&customer).Error; err != nil {
					return err
				}
				saleObj.CustomerId = customer.ID
			} else {
				saleObj.CustomerId = customer.ID
			}
		} else {
			customer := models.Customer{ID: saleObj.CustomerId}
			if TX.Where("company_id = ?", currentUser.CompanyId).First(&customer).RowsAffected == 0 {
				return err
			}
			customer.DocumentNumber = saleObj.DocumentNumber
			customer.SocialReason = saleObj.SocialReason
			customer.UtilIdentityDocumentTypeId = saleObj.UtilIdentityDocumentTypeId
			customer.CompanyId = currentUser.CompanyId
			customer.Address = saleObj.FiscalAddress
			customer.Email = saleObj.Email
			customer.Phone = saleObj.Phone

			if err := TX.Save(&customer).Error; err != nil {
				return err
			}
		}

		// Insert sale in database
		sale.Serie = salePointTokens.Serie
		sale.Number = salePointTokens.Number
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
		sale.UtilDocumentTypeId = saleObj.UtilDocumentTypeId
		sale.UtilCurrencyTypeId = saleObj.UtilCurrencyTypeId
		sale.UtilOperationTypeId = saleObj.UtilOperationTypeId
		sale.CompanyId = currentUser.CompanyId
		sale.CompanyWareHouseId = saleObj.CompanyWareHouseId
		sale.CompanySalePointId = saleObj.CompanySalePointId
		sale.UserId = currentUser.ID

		if err := TX.Create(&sale).Error; err != nil {
			return err
		}

		wareHouse := models.CompanyWareHouse{ID: saleObj.CompanyWareHouseId}
		if TX.First(&wareHouse).RowsAffected == 0 {
			return err
		}

		for _, item := range saleObj.Item {
			item.SaleId = sale.ID
			if err := TX.Create(&item).Error; err != nil {
				return err
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
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Nubefact
	var nuResponse utilities.Response
	var invoicePath string
	if sale.UtilDocumentTypeId == 1 || sale.UtilDocumentTypeId == 2 {
		nuResponse, err = sendDocumentNubefact(DB, sale.ID, salePointTokens.NuUrl, salePointTokens.NuToken)
		if err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else if sale.UtilDocumentTypeId == 6 {
	    company := models.Company{}
        if DB.Where("id = ?", currentUser.CompanyId).First(&company).RowsAffected == 0 {
            return err
        }
		invoicePath, err = buildSaleTicketPdf(DB, company, sale.ID)
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data: newSaleResponse{
			NubefactResponse: nuResponse,
			InvoicePath:      invoicePath,
		},
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
	Number uint
}

type tokensCompanySalePointSerie struct {
	Serie   string
	NuToken string
	NuUrl   string
}

type tokensCompanySalePoint struct {
	Serie   string
	Number  uint
	NuToken string
	NuUrl   string
}

func getTokensCompanySalePoint(DB *gorm.DB, companySalePointId uint, documentTypeId uint, contingency bool) (error, tokensCompanySalePoint) {
	tokensCompanySalePoints := tokensCompanySalePoint{}

	tokensCompanySalePointSeries := tokensCompanySalePointSerie{}
	if err := DB.Table("company_sale_point_series").Select("company_sale_point_series.serie, company_sale_points.nu_token, company_sale_points.nu_url").
		Joins("INNER JOIN company_sale_points ON company_sale_point_series.company_sale_point_id = company_sale_points.id").
		Where("company_sale_point_id = ?", companySalePointId).
		Where("util_document_type_id = ?", documentTypeId).
		Where("contingency = ?", contingency).
		Scan(&tokensCompanySalePointSeries).Error; err != nil {
		return err, tokensCompanySalePoints
	}

	if tokensCompanySalePointSeries.Serie == "" {
		return errors.New("No se encontró una seri configurada para este punto de venta"), tokensCompanySalePoints
	}

	tokensCompanySalePointNumbers := tokensCompanySalePointNumber{}
	if err := DB.Raw("SELECT CASE WHEN MAX(number) IS NULL THEN 0 ELSE MAX(number) END AS number FROM sales WHERE company_sale_point_id = ? AND util_document_type_id = ? AND serie = ?", companySalePointId, documentTypeId, tokensCompanySalePointSeries.Serie).
		Scan(&tokensCompanySalePointNumbers).Error; err != nil {
		return err, tokensCompanySalePoints
	}

	tokensCompanySalePoints.Serie = tokensCompanySalePointSeries.Serie
	tokensCompanySalePoints.NuUrl = tokensCompanySalePointSeries.NuUrl
	tokensCompanySalePoints.NuToken = tokensCompanySalePointSeries.NuToken
	tokensCompanySalePoints.Number = tokensCompanySalePointNumbers.Number + 1

	return nil, tokensCompanySalePoints
}

func sendDocumentNubefact(DB *gorm.DB, saleId uint, url string, token string) (utilities.Response, error) {
	docResponse := utilities.Response{}

	// Query document
	nDoc := NubefactDocument{}
	if err := DB.Debug().Raw("SELECT doc.nu_code as tipo_de_comprobante, sale.serie, sale.number as numero, ide_doc.nu_code as cliente_tipo_de_documento, "+
		" sale.document_number as cliente_numero_de_documento, sale.social_reason as cliente_denominacion, sale.fiscal_address as cliente_direccion, sale.email as cliente_email, "+
		" sale.date_of_issue as fecha_de_emision, sale.date_of_due as fecha_de_vencimiento, curren.nu_code as moneda, "+
		" sale.total_discount as total_descuento, sale.total_taxed as total_gravada, sale.total_unaffected as total_inafecta, sale.total_exonerated as total_exonerada, sale.total_igv as total_igv, sale.total_free as total_gratuita, sale.total, "+
		" sale.pdf_format as formato_de_pdf "+
		" FROM sales as sale"+
		" INNER JOIN util_document_types as doc ON sale.util_document_type_id = doc.id "+
		" INNER JOIN util_identity_document_types as ide_doc ON sale.util_identity_document_type_id = ide_doc.id "+
		" INNER JOIN util_currency_types as curren ON sale.util_currency_type_id = curren.id "+
		" INNER JOIN util_operation_types as ope ON sale.util_operation_type_id = ope.id "+
		" WHERE sale.id = ? LIMIT 1", saleId).Scan(&nDoc).Error; err != nil {
		return docResponse, err
	}

	nDoc.Operacion = "generar_comprobante"
	nDoc.SunatTransaction = "1"
	nDoc.TipoDeCambio = "1"
	nDoc.PorcentajeDeIgv = "18.00"
	nDoc.CodigoUnico = nDoc.Serie + nDoc.Numero
	nDoc.EnviarAutomaticamenteAlCliente = true
	nDoc.EnviarAutomaticamenteALaSunat = true

	nDocItem := make([]NubefactDocumentItem, 0)
	if err := DB.Raw("SELECT unit.nu_code as unidad_de_medida, sale_item.description as descripcion, sale_item.quantity as cantidad, "+
		" sale_item.unit_value as valor_unitario, sale_item.unit_price as precio_unitario, sale_item.discount as descuento, "+
		" sale_item.total_value as subtotal, afe.nu_code as tipo_de_igv, sale_item.igv as igv, sale_item.total as total "+
		" FROM sale_items as sale_item "+
		" INNER JOIN util_unit_measure_types as unit ON sale_item.util_unit_measure_type_id = unit.id "+
		" INNER JOIN util_affectation_igv_types as afe ON sale_item.util_affectation_igv_type_id = afe.id "+
		" WHERE sale_item.sale_id = ? ", saleId).Scan(&nDocItem).Error; err != nil {
		return docResponse, err
	}
	nDoc.Items = nDocItem

	// Send document
	res, err := NubefactSendDocument(nDoc, url, token, false)
	return res, err
}

type  buildSaleTicketTemplate struct {
    ID                uint  
    DateOfIssue       time.Time 
    DateOfSale        time.Time 
    Serie             string
    Number            string
    Observation       string
    CancelObservation string

    // Customer
    DocumentNumber             string
    SocialReason               string
    FiscalAddress              string
    Email                      string
    Phone                      string

    TotalUnaffected   float64
    TotalTaxed        float64
    TotalIgv          float64
    Total             float64

    DocumentName   string
    UserName       string
    CurrencySymbol string
    CurrencyDescription string
    IdentityDocumentName string
}

type  buildSaleItemTicketTemplate struct {
    ID uint `json:"id" gorm:"primaryKey"`
    ProductCode           string 
    Description           string 
    Quantity              float64 
    UnitValue             float64
    UnitPrice             float64 
    Discount float64 
    TotalValue float64
    Total      float64
    UnitMeasureCode string
}

func buildSaleTicketPdf(DB *gorm.DB, con models.Company,  saleId uint) (string, error) {
    // Get data sale
    sale := buildSaleTicketTemplate{}
    if err := DB.Table("sales").Select("sales.*, util_document_types.description as document_name, users.user_name as user_name, " +
        " util_currency_types.symbol as currency_symbol, util_currency_types.description as currency_description, " +
        " util_identity_document_types.description as identity_document_name").
        Joins("INNER JOIN util_document_types ON sales.util_document_type_id = util_document_types.id").
        Joins("INNER JOIN util_identity_document_types ON sales.util_identity_document_type_id = util_identity_document_types.id").
        Joins("INNER JOIN users ON sales.user_id = users.id").
        Joins("INNER JOIN util_currency_types ON sales.util_currency_type_id = util_currency_types.id").
        Where("sales.id = ?", saleId).Limit(1).Scan(&sale).Error; err != nil {
            return "", err
    }

    saleItems := make([]buildSaleItemTicketTemplate, 0)
    if err := DB.Table("sale_items").Select("sale_items.*, util_unit_measure_types.code as unit_measure_code").
        Joins("INNER JOIN util_unit_measure_types ON sale_items.util_unit_measure_type_id = util_unit_measure_types.id").
        Where("sale_items.sale_id = ?", saleId).Scan(&saleItems).Error; err != nil {
        return "", err
    }

    geographicalLocation := models.UtilGeographicalLocation{}
    if DB.Where("id = ?", con.UtilGeographicalLocationId).First(&geographicalLocation).RowsAffected == 0 {
        return "", nil
    }

	// Settings
	pageMargin := 3.0

	// Create PDF
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size:    gofpdf.SizeType{Wd: 74.1, Ht: 140},
	})
	pdf.SetMargins(pageMargin, pageMargin, pageMargin)
	pdf.AddUTF8Font("Calibri", "", "static/font/Calibri_Regular.ttf")
	pdf.AddUTF8Font("Calibri", "B", "static/font/Calibri_Bold.ttf")
	pdf.AddUTF8Font("Calibri", "I", "static/font/Calibri_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "BI", "static/font/Calibri_Bold_Italic.ttf")
	pdf.AddUTF8Font("Calibri", "L", "static/font/Calibri_Light.ttf")
	pdf.AddUTF8Font("Calibri", "LI", "static/font/Calibri_Light_Italic.ttf")

	//Settings
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	pageWidth -= leftMargin + rightMargin
	fontFamilyName := "Calibri"
	fontSize:=7.0
	//gutter := 2.0
    lineHeight := 2.5

    // Init
    pdf.AddPage()
    pdf.SetFont(fontFamilyName, "B", fontSize)
	//offsetTop := topMargin
    //pdf.SetY()

    // Header
    pdf.WriteAligned(0, lineHeight, strings.ToUpper(con.SocialReason), "C")
    pdf.Ln(lineHeight)

    pdf.SetFont(fontFamilyName, "", fontSize)
    pdf.WriteAligned(0, lineHeight, strings.ToUpper(con.Address), "C")
    pdf.Ln(lineHeight)

    pdf.WriteAligned(0, lineHeight, strings.ToUpper(geographicalLocation.District) + " - " + strings.ToUpper(geographicalLocation.Province) + " - " + strings.ToUpper(geographicalLocation.Department), "C")
    pdf.Ln(lineHeight)

    pdf.SetFont(fontFamilyName, "B", fontSize + 2)
    pdf.WriteAligned(0, lineHeight + 0.5, sale.DocumentName, "C")
    pdf.Ln(lineHeight)

    pdf.WriteAligned(0, lineHeight + 0.5, sale.Serie + "-" + sale.Number, "C")
    pdf.Ln(lineHeight * 2)

    pdf.SetFont(fontFamilyName, "B", fontSize)
    pdf.WriteAligned(0, lineHeight, "CLIENTE", "L")
    pdf.Ln(lineHeight)

    pdf.SetFont(fontFamilyName, "", fontSize)
    pdf.WriteAligned(0, lineHeight, sale.IdentityDocumentName + ": " + sale.DocumentNumber, "L")
    pdf.Ln(lineHeight)

    pdf.WriteAligned(0, lineHeight, sale.SocialReason, "L")
    pdf.Ln(lineHeight)

    pdf.SetFont(fontFamilyName, "B", fontSize)
    pdf.WriteAligned(0, lineHeight, "FECHA EMISIÓN: ", "L")
    pdf.SetFont(fontFamilyName, "", fontSize)
    pdf.WriteAligned(0, lineHeight, fmt.Sprintf("%d-%02d-%02d", sale.DateOfIssue.Year(), sale.DateOfIssue.Month(), sale.DateOfIssue.Day()), "L")
    pdf.Ln(lineHeight)

    pdf.SetFont(fontFamilyName, "B", fontSize)
    pdf.WriteAligned(0, lineHeight, "MONEDA: ", "L")
    pdf.SetFont(fontFamilyName, "", fontSize)
    pdf.WriteAligned(0, lineHeight, sale.CurrencyDescription, "L")
    pdf.Ln(lineHeight * 1.5)

    // Column widths
    w := []float64{8.0, 38.0, 10.0, 10.0}
    wSum := 0.0
    for _, v := range w {
        wSum += v
    }
    left := leftMargin + 1

    // 	Table header
    pdf.SetX(left)
    pdf.SetDrawColor(206, 206, 206)
    pdf.SetFont(fontFamilyName, "B", fontSize)
    for j, str := range []string{"CANT", "DESCRIPCIÓN", "P/U", "Total"} {
        pdf.CellFormat(w[j], lineHeight, str, "TB", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)

    // Table body
    pdf.SetFont(fontFamilyName, "", fontSize)
    for _, item := range saleItems {
        pdf.SetX(left)
        pdf.CellFormat(w[0], lineHeight, fmt.Sprintf("%.1f",item.Quantity), "TB", 0, "C", false, 0, "")
        pdf.CellFormat(w[1], lineHeight, item.Description, "TB", 0, "", false, 0, "")
        pdf.CellFormat(w[2], lineHeight, fmt.Sprintf("%.2f",item.UnitPrice), "TB", 0, "C", false, 0, "")
        pdf.CellFormat(w[3], lineHeight, fmt.Sprintf("%.2f",item.Total), "TB", 0, "C", false, 0, "")
        pdf.Ln(-1)
    }
    pdf.Ln(lineHeight / 2)

    // Footer
    pdf.WriteAligned(0, lineHeight, fmt.Sprintf("Total: %s %.2f", sale.CurrencySymbol, sale.Total), "R")
    pdf.Ln(lineHeight)

    pdf.WriteAligned(0, lineHeight, fmt.Sprintf("User: %s", sale.UserName), "L")
    pdf.Ln(lineHeight)

    // Set file name
	cc := sha256.Sum256([]byte(fmt.Sprintf("%d-", saleId)))
	pwd := fmt.Sprintf("%x", cc)
	fileName := fmt.Sprintf("static/temp/%s.pdf", pwd)

	// Remove all files
    err := utilities.ClearTempFolder(false)
    if err != nil {
        return "", err
    }

    // Save file
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
