package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
	"net/http"
)

//
type webSiteCardResponse struct {
	ID uint `json:"id"`

	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	//LongDescription string  `json:"long_description"`
	Barcode string `json:"barcode"`
	//IsService       bool    `json:"is_service"`
	//Location        string  `json:"location"`
	//StockMin        float32 `json:"stock_min"`
	//StockMax        float32 `json:"stock_max"`

	//InternalUse   bool    `json:"internal_use"`
	//Favourite     bool    `json:"favourite"`

	//Lot    bool    `json:"lot"`                      // Lote
	//Bulk   bool    `json:"bulk"`                     // Granel
	//Recipe bool    `json:"recipe"`                   // Receta medica
	//Weight float32 `json:"weight" gorm:"default: 0"` // Peso

	SalePrice1 float64 `json:"sale_price_1"`
	//SalePrice2 float64 `json:"sale_price_2"`
	//SalePrice3 float64 `json:"sale_price_3"`
	//SalePrice4 float64 `json:"sale_price_4"`

	WholeSale1 float64 `json:"whole_sale_1"`
	//WholeSale2 float64 `json:"whole_sale_2"`
	//WholeSale3 float64 `json:"whole_sale_3"`
	//WholeSale4 float64 `json:"whole_sale_4"`

	//PurchaseUtilUnitMeasureTypeId uint    `json:"purchase_util_unit_measure_type_id"`
	//SaleUtilUnitMeasureTypeId     uint    `json:"sale_util_unit_measure_type_id"`
	//Factor                        float32 `json:"factor"`

	//CompanyId  uint `json:"company_id"`
	//CategoryId uint `json:"category_id"`
}

// WebSiteCard function forgot user search
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
	products := make([]webSiteCardResponse, 0)

	// Find products
	if err := DB.Table("products").Select("products.id, products.url, products.title, products.description, products.barcode, products.sale_price1, products.whole_sale1").
		Where("lower(products.title) LIKE lower(?)", "%"+request.Search+"%").
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
