package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

type kardexPage struct {
	ID                  uint      `json:"id"`
	DateOfIssue         time.Time `json:"date_of_issue"`
	Quantity            float64   `json:"quantity"`
	UnitPrice           float64   `json:"unit_price"`
	Total               float64   `json:"total"`
	Stock               float64   `json:"stock"`
	Origin              string    `json:"origin"`
	Destination         string    `json:"destination"`
	Description         string    `json:"description"`
	DocumentDescription string    `json:"document_description"`
	IsLast              bool      `json:"is_last"`
	IsIncome            bool      `json:"is_income"`

	UserId             uint `json:"user_id"`
	ProductId          uint `json:"product_id"`
	CompanyWareHouseId uint `json:"company_ware_house_id"`

	UserName string `json:"user_name"`
}

// PaginateKardex function get all products
func PaginateKardex(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "inventory_kardex"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	products := make([]kardexPage, 0)

	// Find products
	if err := DB.Table("kardexes").Select("kardexes.*, users.user_name").
		Joins("INNER JOIN users ON kardexes.user_id = users.id").
		Where("kardexes.company_ware_house_id = ? AND (kardexes.date_of_issue BETWEEN ? AND ?)", request.WareHouseId, request.StartDate, request.EndDate).
		Order("kardexes.id asc").Offset(offset).Limit(request.PageSize).Scan(&products).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     products,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}
