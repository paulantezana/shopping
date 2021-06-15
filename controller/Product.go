package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
	"net/http"
)

// PaginateProduct function get all products
func PaginateProduct(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	products := make([]models.Product, 0)

	// Find products
	if err := DB.Table("products").Select("products.*, kardexes.stock").
		Joins("LEFT JOIN kardexes ON products.id = kardexes.product_id AND kardexes.company_ware_house_id = ? AND kardexes.is_last = true", request.WareHouseId).
		Where("products.company_id = ? AND lower(products.title) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
		Order("products.id desc").Offset(offset).Limit(request.PageSize).Scan(&products).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     products,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

type productSeeker struct {
	ID uint `json:"id"`

	Url             string  `json:"url"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	LongDescription string  `json:"long_description"`
	Barcode         string  `json:"barcode"`
	IsService       bool    `json:"is_service"`
	Location        string  `json:"location"`
	StockMin        float32 `json:"stock_min"`
	StockMax        float32 `json:"stock_max"`

	InternalUse   bool    `json:"internal_use"`
	Favourite     bool    `json:"favourite"`
	PurchasePrice float64 `json:"purchase_price"`

	Lot    bool    `json:"lot"`    // Lote
	Bulk   bool    `json:"bulk"`   // Granel
	Recipe bool    `json:"recipe"` // Receta medica
	Weight float32 `json:"weight"` // Peso

	SalePrice1 float64 `json:"sale_price_1"`
	SalePrice2 float64 `json:"sale_price_2"`
	SalePrice3 float64 `json:"sale_price_3"`
	SalePrice4 float64 `json:"sale_price_4"`

	WholeSale1 float64 `json:"whole_sale_1"`
	WholeSale2 float64 `json:"whole_sale_2"`
	WholeSale3 float64 `json:"whole_sale_3"`
	WholeSale4 float64 `json:"whole_sale_4"`

	PurchaseUtilUnitMeasureTypeId uint    `json:"purchase_util_unit_measure_type_id"`
	SaleUtilUnitMeasureTypeId     uint    `json:"sale_util_unit_measure_type_id"`
	Factor                        float32 `json:"factor"`

	State bool `json:"state" gorm:"default: true"`

	Stock                          float32 `json:"stock"`
	PurchaseUnitMeasureDescription string  `json:"purchase_unit_measure_description"`
	SaleUnitMeasureDescription     string  `json:"sale_unit_measure_description"`
	SearchText                     string  `json:"search_text"`
}

// PaginateProductSeekerSearch function get all products
func PaginateProductSeekerSearch(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	productSeekers := make([]productSeeker, 0)

	// Find products
	if err := DB.Table("products").Select("products.*, kardexes.stock, pur_unit.code as purchase_unit_measure_description, inv_unit.code  as sale_unit_measure_description").
		Joins("INNER JOIN util_unit_measure_types as pur_unit ON products.purchase_util_unit_measure_type_id = pur_unit.id").
		Joins("INNER JOIN util_unit_measure_types as inv_unit ON products.sale_util_unit_measure_type_id = inv_unit.id").
		Joins("LEFT JOIN kardexes ON products.id = kardexes.product_id AND kardexes.company_ware_house_id = ? AND kardexes.is_last = true", request.WareHouseId).
		Where("products.company_id = ? AND (lower(products.title) LIKE lower(?) OR lower(products.barcode) LIKE lower(?)) AND products.state = true", currentUser.CompanyId, "%"+request.Search+"%", "%"+request.Search+"%").
		Order("products.id desc").Offset(offset).Limit(request.PageSize).Scan(&productSeekers).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     productSeekers,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetProductByID function get product by id
func GetProductByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&product, product.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product,
	})
}

// GetProductSearch function get product by id
func GetProductSearch(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Check the number of matches
	products := make([]models.Product, 0)

	// Find users
	if err := DB.Raw("SELECT * FROM (SELECT *, CONCAT(barcode, ' ', title) as search_text FROM products) as product_aux "+
		" WHERE product_aux.company_id = ? AND lower(product_aux.search_text) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
		Scan(&products).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    products,
	})
}

// CreateProduct function create new product
func CreateProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert product in database
	product.CreatedUserId = currentUser.ID
	product.CompanyId = currentUser.CompanyId
	if err := DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product.ID,
		Message: fmt.Sprintf("El producto %s se registro exitosamente", product.Title),
	})
}

// UpdateProduct function update current product
func UpdateProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation product exist
	aux := models.Product{ID: product.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", product.ID),
		})
	}

	// Update product in database
	product.UpdatedUserId = currentUser.ID
	product.State = aux.State
	if err := DB.Model(&product).Updates(product).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !product.State {
		if err := DB.Model(&product).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El producto se actualizó correctamente",
		Data:    product.ID,
	})
}

// UpdateStateProduct function update current product
func UpdateStateProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update product in database
	if !product.State {
		if err := DB.Model(product).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(product).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El producto se actualizó correctamente",
		Data:    product.ID,
	})
}

// GetProductPurchaseByCode function get product by id
func GetProductSeekerByCode(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	DB.Table("products").Select("products.*, util_unit_measure_types.description as purchase_unit_measure_description").
		Joins("LEFT JOIN util_unit_measure_types ON products.purchase_util_unit_measure_type_id = util_unit_measure_types.id").
		Where("products.company_id = ? AND lower(products.barcode) = lower(?) AND products.state = true", currentUser.CompanyId, product.Barcode).
		Limit(1).Scan(&product)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product,
	})
}
