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

// PaginateExpenseType function get all expenseTypes
func PaginateExpenseType(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_expense_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	expenseTypes := make([]models.ExpenseType, 0)

	// Find expenseTypes
    if err := DB.Where("company_id = ? AND lower(description) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
        Order("id desc").Offset(offset).Limit(request.PageSize).Find(&expenseTypes).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     expenseTypes,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
        Message: "Consulta exitosa",
	})
}

// GetAllExpenseType function get all types
func GetAllExpenseType(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get connection
    DB := provider.GetConnection()
    // defer db.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_expense_type"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Check the number of matches
    expenseTypes := make([]models.ExpenseType, 0)

    // Find users
    if err := DB.Where("company_id = ?", currentUser.CompanyId).Find(&expenseTypes).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    expenseTypes,
        Message: "Consulta exitosa",
    })
}


// GetExpenseTypeByID function get expenseType by id
func GetExpenseTypeByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	expenseType := models.ExpenseType{}
	if err := c.Bind(&expenseType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_expense_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&expenseType, expenseType.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    expenseType,
        Message: "Consulta exitosa",
	})
}

// CreateExpenseType function create new expenseType
func CreateExpenseType(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	expenseType := models.ExpenseType{}
	if err := c.Bind(&expenseType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_expense_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert expenseType in database
	expenseType.CreatedUserId = currentUser.ID
	expenseType.CompanyId = currentUser.CompanyId
	if err := DB.Create(&expenseType).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    expenseType.ID,
		Message: fmt.Sprintf("El tipo de ingreso %s se registro exitosamente", expenseType.Description),
	})
}

// UpdateExpenseType function update current expenseType
func UpdateExpenseType(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	expenseType := models.ExpenseType{}
	if err := c.Bind(&expenseType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_expense_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation expenseType exist
	aux := models.ExpenseType{ID: expenseType.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", expenseType.ID),
		})
	}

	// Update expenseType in database
	expenseType.UpdatedUserId = currentUser.ID
	if err := DB.Model(&expenseType).Updates(expenseType).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !expenseType.State {
		if err := DB.Model(expenseType).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El tipo de ingreso se actualizó correctamente",
		Data:    expenseType.ID,
	})
}

// UpdateStateExpenseType function update current expenseType
func UpdateStateExpenseType(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	expenseType := models.ExpenseType{}
	if err := c.Bind(&expenseType); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_expense_type"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update expenseType in database
	if !expenseType.State {
		if err := DB.Model(expenseType).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(expenseType).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El tipo de ingreso se actualizó correctamente",
		Data:    expenseType.ID,
	})
}
