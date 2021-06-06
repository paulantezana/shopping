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

// PaginateCustomer function get all customers
func PaginateCustomer(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	customers := make([]models.Customer, 0)

	// Find users
	if err := DB.Where("company_id = ? AND lower(social_reason) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
		Order("id desc").Offset(offset).Limit(request.PageSize).Find(&customers).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     customers,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetAllCustomer function get all customers
func GetAllCustomer(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Check the number of matches
	customers := make([]models.Customer, 0)

	// Find users
	if err := DB.Where("company_id = ?", currentUser.CompanyId).Find(&customers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customers,
	})
}

// GetAllCustomer function get all customers
func GetSearchCustomer(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// request
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Check the number of matches
	customers := make([]models.Customer, 0)

	// Find users
	if err := DB.Raw("SELECT * FROM (SELECT *, CONCAT(document_number, ' ', social_reason) as search_text FROM customers) as customer_aux "+
		" WHERE customer_aux.company_id = ? AND lower(customer_aux.search_text) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
		Scan(&customers).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customers,
	})
}

// GetCustomerByID function get customer by id
func GetCustomerByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	customerObj := models.Customer{}
	if err := c.Bind(&customerObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&customerObj, customerObj.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customerObj,
	})
}

// CreateCustomer function create new customer
func CreateCustomer(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	customerObj := models.Customer{}
	if err := c.Bind(&customerObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert customer in database
	customerObj.CreatedUserId = currentUser.ID
	customerObj.CompanyId = currentUser.CompanyId
	if err := DB.Create(&customerObj).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customerObj.ID,
		Message: fmt.Sprintf("El cliente %s se registro exitosamente", customerObj.SocialReason),
	})
}

// UpdateCustomer function update current customer
func UpdateCustomer(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	customerObj := models.Customer{}
	if err := c.Bind(&customerObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation customer exist
	aux := models.Customer{ID: customerObj.ID}
	if DB.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", customerObj.ID),
		})
	}

	// Update customer in database
	customerObj.UpdatedUserId = currentUser.ID
	if err := DB.Model(&customerObj).Update(customerObj).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !customerObj.State {
		if err := DB.Model(customerObj).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
        Message: fmt.Sprintf("El cliente %s se actualizó correctamente", customerObj.SocialReason),
		Data:    customerObj.ID,
	})
}

// UpdateStateCustomer function update current customer
func UpdateStateCustomer(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	customerObj := models.Customer{}
	if err := c.Bind(&customerObj); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "sale_customer"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update customer in database
	if !customerObj.State {
		if err := DB.Model(customerObj).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(customerObj).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El cliente se actualizó correctamente",
		Data:    customerObj.ID,
	})
}
