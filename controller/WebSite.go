package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// ForgotSearch function forgot user search
func WebSiteCard(c echo.Context) error {
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

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	products := make([]models.Product, 0)

	// Find products
	if err := DB.Table("products").Select("products.*, product_ware_houses.stock").
		Joins("LEFT JOIN product_ware_houses ON products.id = product_ware_houses.product_id AND product_ware_houses.ware_house_id = ?", request.WareHouseId).
		Where("products.company_id = ? AND lower(products.title) LIKE lower(?)", request.CompanyId, "%"+request.Search+"%").
		Order("products.id desc").Offset(offset).Limit(request.PageSize).Scan(&products).
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
