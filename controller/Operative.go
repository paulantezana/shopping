package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
	"net/http"
)

type operativeRequest struct {
	SaleConf     models.SaleConf     `json:"sale_conf"`
	PurchaseConf models.PurchaseConf `json:"purchase_conf"`
}

// GetOperative get
func GetOperative(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_operative"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Query
	saleConf := models.SaleConf{}
	if err := DB.First(&saleConf, currentUser.CompanyId).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query
	purchaseConf := models.PurchaseConf{}
	if err := DB.First(&purchaseConf, currentUser.CompanyId).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	or := operativeRequest{}
	or.SaleConf = saleConf
	or.PurchaseConf = purchaseConf

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    or,
	})
}

// SaveOperative save
func SaveOperative(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	ope := operativeRequest{}
	if err := c.Bind(&ope); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_operative"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update sale conf
	if err := DB.Model(&ope.SaleConf).Updates(ope.SaleConf).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Update purchase conf
	if err := DB.Model(&ope.PurchaseConf).Updates(ope.PurchaseConf).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Message: "Las configuraciones de operatividad se guardaron exitosamente",
	})
}
