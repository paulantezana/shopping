package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// GetAllUtilAdditionalLegendType --
func GetAllUtilAdditionalLegendType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilAdditionalLegendType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilCatAffectationIgvType --
func GetAllUtilCatAffectationIgvType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilCatAffectationIgvType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilCreditDebitType --
func GetAllUtilCreditDebitType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilCreditDebitType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilCurrencyType --
func GetAllUtilCurrencyType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilCurrencyType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilDocumentType --
func GetAllUtilDocumentType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilDocumentType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilGeographicalLocation --
func GetAllUtilGeographicalLocation(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilGeographicalLocation, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetSearchUtilGeographicalLocation --
func GetSearchUtilGeographicalLocation(c echo.Context) error {
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	DB := provider.GetConnection()
	defer DB.Close()

	utilGeographicalLocationShorts := make([]models.UtilGeographicalLocationShort, 0)

	// Find users
	if err := DB.Raw("SELECT * FROM (SELECT id, code, concat(department, '-', province, '-', district) as description  FROM util_geographical_locations) as geo WHERE description LIKE ?", "%"+request.Search+"%").
		Scan(&utilGeographicalLocationShorts).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    utilGeographicalLocationShorts,
	})
}

// GetAllUtilIdentityDocumentType --
func GetAllUtilIdentityDocumentType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilIdentityDocumentType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilOperationType --
func GetAllUtilOperationType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilOperationType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilPerceptionType --
func GetAllUtilPerceptionType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilPerceptionType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilProductType --
func GetAllUtilProductType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilProductType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilSubjectDetractionType --
func GetAllUtilSubjectDetractionType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilSubjectDetractionType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilSystemIscType --
func GetAllUtilSystemIscType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilSystemIscType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilTransferReasonType --
func GetAllUtilTransferReasonType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilTransferReasonType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilTransportModeType --
func GetAllUtilTransportModeType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilTransportModeType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilTributeType --
func GetAllUtilTributeType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilTributeType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}

// GetAllUtilUnitMeasureType --
func GetAllUtilUnitMeasureType(c echo.Context) error {
	DB := provider.GetConnection()
	defer DB.Close()

	util := make([]models.UtilUnitMeasureType, 0)

	// Find users
	if err := DB.Find(&util).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    util,
	})
}
