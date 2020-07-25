package controller

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// PaginateCompanyLocal function get all companylocals
func PaginateCompanyLocal(c echo.Context) error {
	// Get companyLocal token authenticate
	// companyLocal := c.Get("companyLocal").(*jwt.Token)
	// claims := companyLocal.Claims.(*utilities.Claim)
	// currentCompanyLocal := claims.CompanyLocal

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
	companyLocals := make([]models.CompanyLocal, 0)

	// Find companyLocals
	if err := DB.Where("social_reason LIKE ?", "%"+request.Search+"%").
		Order("id desc").Offset(offset).Limit(request.PageSize).Find(&companyLocals).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     companyLocals,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// CompanyLocalRequestId --
type CompanyLocalRequestId struct {
	CompanyLocal                  models.CompanyLocal                  `json:"company_local"`
	UtilGeographicalLocationShort models.UtilGeographicalLocationShort `json:"util_geographical_location_short"`
}

// GetCompanyLocalByID function get companyLocal by id
func GetCompanyLocalByID(c echo.Context) error {
	// Get data request
	companyLocal := models.CompanyLocal{}
	if err := c.Bind(&companyLocal); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Execute instructions
	companyLocalRequest := CompanyLocalRequestId{}
	if err := DB.First(&companyLocal, companyLocal.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	companyLocalRequest.CompanyLocal = companyLocal
	if err := DB.Raw("SELECT id, code, concat(department, '-', province, '-', district) as description  FROM util_geographical_locations WHERE id = ?", companyLocal.UtilGeographicalLocationId).
		Scan(&companyLocalRequest.UtilGeographicalLocationShort).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if err := DB.Where("company_local_id = ?", companyLocal.ID).Find(&companyLocalRequest.CompanyLocal.CompanySeries).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companyLocalRequest,
	})
}

// CreateCompanyLocal function create new companyLocal
func CreateCompanyLocal(c echo.Context) error {
	// Get data request
	companyLocal := models.CompanyLocal{}
	if err := c.Bind(&companyLocal); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Validate
	valid := ValidateCompanyLocal(companyLocal)
	if !valid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: valid.Message,
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Insert companyLocal in database
	if err := DB.Create(&companyLocal).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companyLocal.ID,
		Message: fmt.Sprintf("El sucursal %s se registro exitosamente", companyLocal.SocialReason),
	})
}

// UpdateCompanyLocal function update current companyLocal
func UpdateCompanyLocal(c echo.Context) error {
	// Get data request
	companyLocal := models.CompanyLocal{}
	if err := c.Bind(&companyLocal); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Validate
	valid := ValidateCompanyLocal(companyLocal)
	if !valid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: valid.Message,
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validation companyLocal exist
	aux := models.CompanyLocal{ID: companyLocal.ID}
	if DB.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", companyLocal.ID),
		})
	}

	// Update companyLocal in database
	if err := DB.Model(&companyLocal).Update(companyLocal).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !companyLocal.State {
		if err := DB.Model(companyLocal).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El usuario se actualizó correctamente",
		Data:    companyLocal.ID,
	})
}

// DeleteCompanyLocal function delete companyLocal by id
func DeleteCompanyLocal(c echo.Context) error {
	// Get data request
	companyLocal := models.CompanyLocal{}
	if err := c.Bind(&companyLocal); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate
	if companyLocal.ID == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: "No se especificó la clave del registro",
		})
	}

	// Delete companyLocal in database
	if err := DB.Unscoped().Delete(&companyLocal).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    companyLocal.ID,
		Message: fmt.Sprintf("El sucursal %s, se elimino correctamente", companyLocal.SocialReason),
	})
}

func ValidateCompanyLocal(companyLocal models.CompanyLocal) utilities.Response {
	response := utilities.Response{}
	if companyLocal.SocialReason == "" {
		response.Message += "Falta ingresar el codigo \n"
		return response
	}
	if companyLocal.SocialReason == "" {
		response.Message += "Falta ingresar el nombre de sucursal \n"
		return response
	}
	if companyLocal.Address == "" {
		response.Message += "Falta ingresar el dirección \n"
		return response
	}
	if len(companyLocal.CompanySeries) == 0 {
		response.Message += "Falta ingresar el item \n"
		return response
	}
	for _, companySerie := range companyLocal.CompanySeries {
		cSerie := string(companySerie.Serie[0])
		if len(companySerie.Serie) != 4 {
			response.Message += "La serie debe contener 4 digitos \n"
			return response
		}
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
			if !regexp.MustCompile("^[0-9]{4}$").MatchString(cSerie) {
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
