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

type companyWareHousePaginateResponse struct {
	ID             uint   `json:"id"`
	CompanyLocalId uint   `json:"company_local_id"`
	Description    string `json:"description"`
	CompanyLocal   string `json:"company_local"`
	State          bool   `json:"state"`
}

// PaginateCompanyWareHouse function get all companyWareHouses
func PaginateCompanyWareHouse(c echo.Context) error {
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
	companyWareHouses := make([]companyWareHousePaginateResponse, 0)

	// Find companyWareHouses
	if err := DB.Table("company_ware_houses").Select("company_ware_houses.id, company_ware_houses.description, company_ware_houses.state, company_locals.social_reason as company_local").
		Joins("INNER JOIN company_locals on company_ware_houses.company_local_id = company_locals.id").
		Where("lower(company_ware_houses.description) LIKE lower(?)", "%"+request.Search+"%").
		Order("company_ware_houses.id desc").Offset(offset).Limit(request.PageSize).Scan(&companyWareHouses).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     companyWareHouses,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetCompanyWareHouseByID function get companyWareHouse by id
func GetCompanyWareHouseByID(c echo.Context) error {
	// Get data request
	companyWareHouse := models.CompanyWareHouse{}
	if err := c.Bind(&companyWareHouse); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&companyWareHouse, companyWareHouse.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companyWareHouse,
	})
}

// CreateCompanyWareHouse function create new companyWareHouse
func CreateCompanyWareHouse(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	companyWareHouse := models.CompanyWareHouse{}
	if err := c.Bind(&companyWareHouse); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Default empty values
	// if companyWareHouse.CompanyWareHouseRoleID == 0 {
	// 	companyWareHouse.CompanyWareHouseRoleID = 6
	// }

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Insert companyWareHouse in database
	companyWareHouse.CreatedUserId = currentUser.ID
	if err := db.Create(&companyWareHouse).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companyWareHouse.ID,
		Message: fmt.Sprintf("El almacen %s se registro exitosamente", companyWareHouse.Description),
	})
}

// UpdateCompanyWareHouse function update current companyWareHouse
func UpdateCompanyWareHouse(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	companyWareHouse := models.CompanyWareHouse{}
	if err := c.Bind(&companyWareHouse); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation companyWareHouse exist
	aux := models.CompanyWareHouse{ID: companyWareHouse.ID}
	if db.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", companyWareHouse.ID),
		})
	}

	// Update companyWareHouse in database
	companyWareHouse.UpdatedUserId = currentUser.ID
	if err := db.Model(&companyWareHouse).Update(companyWareHouse).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !companyWareHouse.State {
		if err := db.Model(companyWareHouse).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El almacen se actualizó correctamente",
		Data:    companyWareHouse.ID,
	})
}

// UpdateStateCompanyWareHouse function update current companyWareHouse
func UpdateStateCompanyWareHouse(c echo.Context) error {
    // Get user token authenticate
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	companyWareHouse := models.CompanyWareHouse{}
	if err := c.Bind(&companyWareHouse); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Update companyWareHouse in database
	if !companyWareHouse.State {
		if err := db.Model(companyWareHouse).UpdateColumn("state", false).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := db.Model(companyWareHouse).UpdateColumn("state", true).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El almacen se actualizó correctamente",
		Data:    companyWareHouse.ID,
	})
}
