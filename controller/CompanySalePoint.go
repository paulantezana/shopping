package controller

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_sale_point"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	companySalePoints := make([]companySalePointPaginateResponse, 0)

	// Find companySalePoints
	if err := DB.Table("company_sale_points").Select("company_sale_points.id, company_sale_points.description, company_sale_points.state, company_locals.social_reason as company_local").
		Joins("INNER JOIN company_locals on company_sale_points.company_local_id = company_locals.id").
		Where("company_sale_points.company_id = ? AND lower(company_sale_points.description) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
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
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_sale_point"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&companySalePoint, companySalePoint.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Series
	if err := DB.Where("company_sale_point_id = ?", companySalePoint.ID).Find(&companySalePoint.CompanySalePointSeries).Error; err != nil {
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
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_sale_point"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert companySalePoint in database
	companySalePoint.CreatedUserId = currentUser.ID
	companySalePoint.CompanyId = currentUser.CompanyId
	if err := DB.Create(&companySalePoint).Error; err != nil {
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
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_sale_point"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation companySalePoint exist
	aux := models.CompanySalePoint{ID: companySalePoint.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", companySalePoint.ID),
		})
	}

	// Update companySalePoint in database
	companySalePoint.UpdatedUserId = currentUser.ID
	if err := DB.Model(&companySalePoint).Updates(companySalePoint).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !companySalePoint.State {
		if err := DB.Model(companySalePoint).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}
	for _, item := range companySalePoint.CompanySalePointSeries {
		if item.IsDeleted {
			if err := DB.Delete(&item).Error; err != nil {
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
		} else {
			if err := DB.Save(&item).Error; err != nil {
				return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
			}
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
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	companySalePoint := models.CompanySalePoint{}
	if err := c.Bind(&companySalePoint); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_sale_point"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update companySalePoint in database
	if !companySalePoint.State {
		if err := DB.Model(companySalePoint).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(companySalePoint).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
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

// validateCompanyLocal --
func validateCompanySalePoint(companySalePoint models.CompanySalePoint) utilities.Response {
	response := utilities.Response{}
	if companySalePoint.Description == "" {
		response.Message += "Falta ingresar la descripción \n"
		return response
	}
	if len(companySalePoint.CompanySalePointSeries) == 0 {
		response.Message += "Falta ingresar el item \n"
		return response
	}
	for _, companySerie := range companySalePoint.CompanySalePointSeries {
		if len(companySerie.Serie) != 4 {
			response.Message += "La serie debe contener 4 digitos \n"
			return response
		}
		if companySerie.UtilDocumentTypeId == 0 {
			response.Message += "Especifique el tipo de documento \n"
			return response
		}
		cSerie := string(companySerie.Serie[0])
		if companySerie.UtilDocumentTypeId == 1 && companySerie.Contingency == false {
			if !(cSerie == "F") {
				response.Message += fmt.Sprintf("La serie %s es incorecto para este tipo de documento", companySerie.Serie)
				return response
			}
		}
		if companySerie.UtilDocumentTypeId == 2 && companySerie.Contingency == false {
			if !(cSerie == "B") {
				response.Message += fmt.Sprintf("La serie %s es incorecto para este tipo de documento", companySerie.Serie)
				return response
			}
		}
		if companySerie.UtilDocumentTypeId == 3 && companySerie.Contingency == false {
			if !(cSerie == "F" || cSerie == "B") {
				response.Message += fmt.Sprintf("La serie %s es incorecto para este tipo de documento", companySerie.Serie)
				return response
			}
		}
		if companySerie.UtilDocumentTypeId == 4 && companySerie.Contingency == false {
			if !(cSerie == "F" || cSerie == "B") {
				response.Message += fmt.Sprintf("La serie %s es incorecto para este tipo de documento", companySerie.Serie)
				return response
			}
		}
		if companySerie.Contingency {
			if !regexp.MustCompile("^[0-9]{4}$").MatchString(companySerie.Serie) {
				response.Message += fmt.Sprintf("La serie %s es incorecto para este tipo de documento", companySerie.Serie)
				return response
			}
		}
		if companySerie.UtilDocumentTypeId == 5 {
			if !(cSerie == "T") {
				response.Message += fmt.Sprintf("La serie %s es incorecto para este tipo de documento", companySerie.Serie)
				return response
			}
		}
	}

	response.Success = true
	return response
}
