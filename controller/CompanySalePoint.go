package controller

import (
	"fmt"
	"net/http"

    "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

type companySalePointPaginateResponse struct {
	ID             uint   `json:"id"`
	CompanyLocalId uint   `json:"company_local_id"`
	Description    string `json:"description"`
	CompanyLocal   string `json:"company_local"`
	State          bool   `json:"state"`
}

// PaginateCompanySalePoint function get all companySalePoints
func PaginateCompanySalePoint(c echo.Context) error {
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

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	companySalePoints := make([]companySalePointPaginateResponse, 0)

	// Find companySalePoints
	if err := DB.Table("company_sale_points").Select("company_sale_points.id, company_sale_points.description, company_sale_points.state, company_locals.social_reason as company_local").
		Joins("INNER JOIN company_locals on company_sale_points.company_local_id = company_locals.id").
		Where("lower(company_sale_points.description) LIKE lower(?)", "%"+request.Search+"%").
		Order("company_sale_points.id desc").Offset(offset).Limit(request.PageSize).Scan(&companySalePoints).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     companySalePoints,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetCompanySalePointByID function get companySalePoint by id
func GetCompanySalePointByID(c echo.Context) error {
	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&companySalePoint, companySalePoint.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companySalePoint,
	})
}

// CreateCompanySalePoint function create new companySalePoint
func CreateCompanySalePoint(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Default empty values
	// if companySalePoint.CompanySalePointRoleID == 0 {
	// 	companySalePoint.CompanySalePointRoleID = 6
	// }

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Insert companySalePoint in database
	companySalePoint.CreatedUserId = currentUser.ID
	if err := db.Create(&companySalePoint).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companySalePoint.ID,
		Message: fmt.Sprintf("El punto de venta %s se registro exitosamente", companySalePoint.Description),
	})
}

// UpdateCompanySalePoint function update current companySalePoint
func UpdateCompanySalePoint(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation companySalePoint exist
	aux := models.CompanySalePoint{ID: companySalePoint.ID}
	if db.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", companySalePoint.ID),
		})
	}

	// Update companySalePoint in database
	companySalePoint.UpdatedUserId = currentUser.ID
	if err := db.Model(&companySalePoint).Update(companySalePoint).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !companySalePoint.State {
		if err := db.Model(companySalePoint).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El punto de venta se actualizó correctamente",
		Data:    companySalePoint.ID,
	})
}

// UpdateStateCompanySalePoint function update current companySalePoint
func UpdateStateCompanySalePoint(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update companySalePoint in database
	if !companySalePoint.State {
		if err := db.Model(companySalePoint).UpdateColumn("state", false).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := db.Model(companySalePoint).UpdateColumn("state", true).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El punto de venta se actualizó correctamente",
		Data:    companySalePoint.ID,
	})
}
