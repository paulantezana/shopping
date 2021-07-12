package controller

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// PaginateIncomeType function get all incomeTypes
func PaginateIncomeType(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_income_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	incomeTypes := make([]models.IncomeType, 0)

	// Find incomeTypes
    if err := DB.Where("company_id = ? AND lower(description) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
        Order("id desc").Offset(offset).Limit(request.PageSize).Find(&incomeTypes).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     incomeTypes,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
		Message: "Consulta exitosa",
	})
}

// GetAllIncomeType function get all types
func GetAllIncomeType(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get connection
    DB := provider.GetConnection()
    // defer db.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_income_type"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Check the number of matches
    incomeTypes := make([]models.IncomeType, 0)

    // Find users
    if err := DB.Where("company_id = ?", currentUser.CompanyId).Find(&incomeTypes).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    incomeTypes,
        Message: "Consulta exitosa",
    })
}

// GetIncomeTypeByID function get incomeType by id
func GetIncomeTypeByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	incomeType := models.IncomeType{}
	if err := c.Bind(&incomeType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_income_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&incomeType, incomeType.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    incomeType,
        Message: "Consulta exitosa",
	})
}

// CreateIncomeType function create new incomeType
func CreateIncomeType(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	incomeType := models.IncomeType{}
	if err := c.Bind(&incomeType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_income_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert incomeType in database
	incomeType.CreatedUserId = currentUser.ID
	incomeType.CompanyId = currentUser.CompanyId
	if err := DB.Create(&incomeType).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    incomeType.ID,
		Message: fmt.Sprintf("El tipo de ingreso %s se registro exitosamente", incomeType.Description),
	})
}

// UpdateIncomeType function update current incomeType
func UpdateIncomeType(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	incomeType := models.IncomeType{}
	if err := c.Bind(&incomeType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_income_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation incomeType exist
	aux := models.IncomeType{ID: incomeType.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", incomeType.ID),
		})
	}

	// Update incomeType in database
	incomeType.UpdatedUserId = currentUser.ID
	if err := DB.Model(&incomeType).Updates(incomeType).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !incomeType.State {
		if err := DB.Model(incomeType).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El tipo de ingreso se actualizó correctamente",
		Data:    incomeType.ID,
	})
}

// UpdateStateIncomeType function update current incomeType
func UpdateStateIncomeType(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	incomeType := models.IncomeType{}
	if err := c.Bind(&incomeType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_income_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update incomeType in database
	if !incomeType.State {
		if err := DB.Model(incomeType).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(incomeType).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El tipo de ingreso se actualizó correctamente",
		Data:    incomeType.ID,
	})
}
