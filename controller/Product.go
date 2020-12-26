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
	if err := DB.Table("products").Select("products.*, product_ware_houses.stock").
		Joins("LEFT JOIN product_ware_houses ON products.id = product_ware_houses.product_id AND product_ware_houses.ware_house_id = ?", request.WareHouseId).
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
	if err := DB.Table("products").Select("products.*, product_ware_houses.stock").
		Joins("LEFT JOIN product_ware_houses ON products.id = product_ware_houses.product_id AND product_ware_houses.ware_house_id = ?", request.WareHouseId).
		Where("products.company_id = ? AND (lower(products.title) LIKE lower(?) OR lower(products.barcode) LIKE lower(?) OR lower(products.barcode_aux) LIKE lower(?)) AND products.state = true", currentUser.CompanyId, "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%").
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

// GetProductByCode function get product by id
func GetProductByCode(c echo.Context) error {
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
	DB.Where("lower(barcode) = lower(?) AND state = true", product.Barcode).First(&product)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product,
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
