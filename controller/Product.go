package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
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

// PaginateProduct function get all products
func PaginateProductSearch(c echo.Context) error {
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
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	products := make([]models.Product, 0)

	// Find products
	if err := DB.Table("products").Select("products.*, kardexes.stock").
		Joins("LEFT JOIN kardexes ON products.id = kardexes.product_id AND kardexes.company_ware_house_id = ? AND kardexes.is_last = true", request.WareHouseId).
		Joins("LEFT JOIN util_unit_measure_types ON products.purchase_util_unit_measure_type_id = util_unit_measure_types.id").
		Where("products.company_id = ? AND (lower(products.title) LIKE lower(?) OR lower(products.barcode) LIKE lower(?)) AND products.state = true", currentUser.CompanyId, "%"+request.Search+"%", "%"+request.Search+"%").
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
	defer DB.Close()

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
	defer DB.Close()

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
	defer DB.Close()

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
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation product exist
	aux := models.Product{ID: product.ID}
	if DB.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", product.ID),
		})
	}

	// Update product in database
	product.UpdatedUserId = currentUser.ID
	product.State = aux.State
	if err := DB.Model(&product).Update(product).Error; err != nil {
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
	defer DB.Close()

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
