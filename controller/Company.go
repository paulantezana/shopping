package controller

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// CompanyRequestId --
type CompanyRequestId struct {
	Company                       models.Company                       `json:"company"`
	UtilGeographicalLocationShort models.UtilGeographicalLocationShort `json:"util_geographical_location_short"`
}

// GetCompanyByID function get Company by id
func GetCompanyByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	Company := models.Company{}
	if err := c.Bind(&Company); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_company"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	CompanyRequest := CompanyRequestId{}
	if err := DB.First(&Company, Company.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	CompanyRequest.Company = Company
	if err := DB.Raw("SELECT id, code, concat(department, '-', province, '-', district) as description  FROM util_geographical_locations WHERE id = ?", Company.UtilGeographicalLocationId).
		Scan(&CompanyRequest.UtilGeographicalLocationShort).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    CompanyRequest,
	})
}

// GetFirstCompany function get Company by id
func GetFirstCompany(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_company"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	companyRequest := CompanyRequestId{}
	if err := DB.Where("id = ?", currentUser.CompanyId).First(&companyRequest.Company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if err := DB.Raw("SELECT id, code, concat(department, '-', province, '-', district) as description FROM util_geographical_locations WHERE id = ?", companyRequest.Company.UtilGeographicalLocationId).
		Scan(&companyRequest.UtilGeographicalLocationShort).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    companyRequest,
	})
}

// UpdateCompany function update current Company
func UpdateCompany(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Validate
	valid := ValidateCompany(company)
	if !valid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: valid.Message,
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_company"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation company exist
	aux := models.Company{ID: company.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", company.ID),
		})
	}

	// Update company in database
	company.UpdatedUserId = currentUser.ID
	if err := DB.Model(&company).Updates(company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !company.State {
		if err := DB.Model(company).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "La empresa se actualizó correctamente",
		Data:    company.ID,
	})
}

// UploadLogoCompany function update current Company
func UploadLogoCompany(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	companyId := c.FormValue("id")
	company := models.Company{}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_company"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation Company exist
	if DB.Where("id = ?", companyId).First(&company).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", companyId),
		})
	}

	// Source
	file, err := c.FormFile("logo")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Validate
	isValid := utilities.ValidateUploadFile(file, 100, []string{"JPG", "PNG", "SVG"})
	if !isValid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: isValid.Message,
		})
	}

	// Destination
	ccc := sha256.Sum256([]byte(strconv.Itoa(int(company.ID))))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	logoSRC := "static/company/" + name
	dst, err := os.Create(logoSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	company.Logo = "/" + logoSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update Company in database
	company.UpdatedUserId = currentUser.ID
	if err := DB.Model(&company).Updates(company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El logo se subio correctamente",
		Data:    company.ID,
	})
}

// UploadLogoLargeCompany function update current Company
func UploadLogoLargeCompany(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	companyId := c.FormValue("id")
	company := models.Company{}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_company"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation Company exist
	if DB.Where("id = ?", companyId).First(&company).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", companyId),
		})
	}

	// Source
	file, err := c.FormFile("logo_large")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Validate
	isValid := utilities.ValidateUploadFile(file, 100, []string{"JPG", "PNG", "SVG"})
	if !isValid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: isValid.Message,
		})
	}

	// Destination
	ccc := sha256.Sum256([]byte(strconv.Itoa(int(company.ID))))
	name := fmt.Sprintf("%xlarge%s", ccc, filepath.Ext(file.Filename))
	logoSRC := "static/company/" + name
	dst, err := os.Create(logoSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	company.LogoLarge = "/" + logoSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update Company in database
	company.UpdatedUserId = currentUser.ID
	if err := DB.Model(&company).Updates(company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El logo se subio correctamente",
		Data:    company.ID,
	})
}

// ValidateCompany - validate
func ValidateCompany(Company models.Company) utilities.Response {
	response := utilities.Response{}
	if Company.DocumentNumber == "" {
		response.Message += "Falta ingresar el número del documento \n"
		return response
	}
	if Company.SocialReason == "" {
		response.Message += "Falta ingresar el codigo \n"
		return response
	}
	if Company.CommercialReason == "" {
		response.Message += "Falta ingresar el nombre de sucursal \n"
		return response
	}

	response.Success = true
	return response
}
