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

// PaginateCategory function get all categorys
func PaginateCategory(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_category"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	categorys := make([]models.Category, 0)

	// Find users
	if err := DB.Where("company_id = ? AND lower(name) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
		Order("id desc").Offset(offset).Limit(request.PageSize).Find(&categorys).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     categorys,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetAllCategory function get all categorys
func GetAllCategory(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_category"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Check the number of matches
	categories := make([]models.Category, 0)

	// Find users
	if err := DB.Where("company_id = ?", currentUser.CompanyId).Find(&categories).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    categories,
	})
}

// GetCategoryByID function get category by id
func GetCategoryByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_category"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&category, category.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    category,
	})
}

// CreateCategory function create new category
func CreateCategory(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_category"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert category in database
	category.CreatedUserId = currentUser.ID
	category.CompanyId = currentUser.CompanyId
	if err := DB.Create(&category).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    category.ID,
		Message: fmt.Sprintf("La categoria %s se registro exitosamente", category.Name),
	})
}

// UpdateCategory function update current category
func UpdateCategory(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_category"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation category exist
	aux := models.Category{ID: category.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", category.ID),
		})
	}

	// Update category in database
	category.UpdatedUserId = currentUser.ID
	if err := DB.Model(&category).Updates(category).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !category.State {
		if err := DB.Model(category).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "La categoria se actualizó correctamente",
		Data:    category.ID,
	})
}

// UpdateStateCategory function update current category
func UpdateStateCategory(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_category"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update category in database
	if !category.State {
		if err := DB.Model(category).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(category).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "La categoria se actualizó correctamente",
		Data:    category.ID,
	})
}
