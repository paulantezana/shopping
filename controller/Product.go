package controller

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type paginateProductResponse struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Url             string  `json:"url"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	LongDescription string  `json:"long_description"`
	Barcode         string  `json:"barcode"`
	IsService       bool    `json:"is_service"`
	Location        string  `json:"location"`
	StockMin        float32 `json:"stock_min"`
	StockMax        float32 `json:"stock_max"`

	InternalUse   bool    `json:"internal_use"`
	Favourite     bool    `json:"favourite"`
	PurchasePrice float64 `json:"purchase_price"`

	Lot    bool    `json:"lot"`    // Lote
	Bulk   bool    `json:"bulk"`   // Granel
	Recipe bool    `json:"recipe"` // Receta medica
	Weight float32 `json:"weight"` // Peso

	SalePrice1 float64 `json:"sale_price_1"`
	SalePrice2 float64 `json:"sale_price_2"`
	SalePrice3 float64 `json:"sale_price_3"`
	SalePrice4 float64 `json:"sale_price_4"`

	WholeSale1 float64 `json:"whole_sale_1"`
	WholeSale2 float64 `json:"whole_sale_2"`
	WholeSale3 float64 `json:"whole_sale_3"`
	WholeSale4 float64 `json:"whole_sale_4"`

	PurchaseUtilUnitMeasureTypeId uint    `json:"purchase_util_unit_measure_type_id"`
	SaleUtilUnitMeasureTypeId     uint    `json:"sale_util_unit_measure_type_id"`
	Factor                        float32 `json:"factor"`

	CompanyId  uint `json:"company_id"`
	CategoryId uint `json:"category_id"`
	State      bool `json:"state" gorm:"default: true"`

	// AUX
	Stock float32 `json:"stock"`
}

// PaginateProduct function get all products
func PaginateProduct(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	products := make([]paginateProductResponse, 0)

	// Find products
	if err := DB.Table("products").Select("products.*, kardexes.stock").
		Joins("LEFT JOIN kardexes ON products.id = kardexes.product_id AND kardexes.company_ware_house_id = ? AND kardexes.is_last = true", request.WareHouseId).
		Where("products.company_id = ? AND lower(products.title) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
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

type productSeeker struct {
	ID uint `json:"id"`

	Url             string  `json:"url"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	LongDescription string  `json:"long_description"`
	Barcode         string  `json:"barcode"`
	IsService       bool    `json:"is_service"`
	Location        string  `json:"location"`
	StockMin        float32 `json:"stock_min"`
	StockMax        float32 `json:"stock_max"`

	InternalUse   bool    `json:"internal_use"`
	Favourite     bool    `json:"favourite"`
	PurchasePrice float64 `json:"purchase_price"`

	Lot    bool    `json:"lot"`    // Lote
	Bulk   bool    `json:"bulk"`   // Granel
	Recipe bool    `json:"recipe"` // Receta medica
	Weight float32 `json:"weight"` // Peso

	SalePrice1 float64 `json:"sale_price_1"`
	SalePrice2 float64 `json:"sale_price_2"`
	SalePrice3 float64 `json:"sale_price_3"`
	SalePrice4 float64 `json:"sale_price_4"`

	WholeSale1 float64 `json:"whole_sale_1"`
	WholeSale2 float64 `json:"whole_sale_2"`
	WholeSale3 float64 `json:"whole_sale_3"`
	WholeSale4 float64 `json:"whole_sale_4"`

	PurchaseUtilUnitMeasureTypeId uint    `json:"purchase_util_unit_measure_type_id"`
	SaleUtilUnitMeasureTypeId     uint    `json:"sale_util_unit_measure_type_id"`
	Factor                        float32 `json:"factor"`

	State bool `json:"state" gorm:"default: true"`

	Stock                          float32 `json:"stock"`
	PurchaseUnitMeasureDescription string  `json:"purchase_unit_measure_description"`
	SaleUnitMeasureDescription     string  `json:"sale_unit_measure_description"`
	SearchText                     string  `json:"search_text"`
}

// PaginateProductSeekerSearch function get all products
func PaginateProductSeekerSearch(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total int64
	productSeekers := make([]productSeeker, 0)

	// Find products
	if err := DB.Table("products").Select("products.*, kardexes.stock, pur_unit.code as purchase_unit_measure_description, inv_unit.code  as sale_unit_measure_description").
		Joins("INNER JOIN util_unit_measure_types as pur_unit ON products.purchase_util_unit_measure_type_id = pur_unit.id").
		Joins("INNER JOIN util_unit_measure_types as inv_unit ON products.sale_util_unit_measure_type_id = inv_unit.id").
		Joins("LEFT JOIN kardexes ON products.id = kardexes.product_id AND kardexes.company_ware_house_id = ? AND kardexes.is_last = true", request.WareHouseId).
		Where("products.company_id = ? AND (lower(products.title) LIKE lower(?) OR lower(products.barcode) LIKE lower(?)) AND products.state = true", currentUser.CompanyId, "%"+request.Search+"%", "%"+request.Search+"%").
		Order("products.id desc").Offset(offset).Limit(request.PageSize).Scan(&productSeekers).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     productSeekers,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetProductByID function get product by id
func GetProductByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&product, product.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product,
	})
}

// GetProductSearch function get product by id
func GetProductSearch(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Check the number of matches
	products := make([]models.Product, 0)

	// Find users
	if err := DB.Raw("SELECT * FROM (SELECT *, CONCAT(barcode, ' ', title) as search_text FROM products) as product_aux "+
		" WHERE product_aux.company_id = ? AND lower(product_aux.search_text) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
		Scan(&products).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    products,
	})
}

// CreateProduct function create new product
func CreateProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert product in database
	product.CreatedUserId = currentUser.ID
	product.CompanyId = currentUser.CompanyId
	if err := DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product.ID,
		Message: fmt.Sprintf("El producto %s se registro exitosamente", product.Title),
	})
}

// UpdateProduct function update current product
func UpdateProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation product exist
	aux := models.Product{ID: product.ID}
	if DB.First(&aux).RowsAffected == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", product.ID),
		})
	}

	// Update product in database
	product.UpdatedUserId = currentUser.ID
	product.State = aux.State
	if err := DB.Model(&product).Updates(product).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !product.State {
		if err := DB.Model(&product).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El producto se actualizó correctamente",
		Data:    product.ID,
	})
}

// UpdateStateProduct function update current product
func UpdateStateProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update product in database
	if !product.State {
		if err := DB.Model(product).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(product).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El producto se actualizó correctamente",
		Data:    product.ID,
	})
}

// GetProductSeekerByCode function get product by id
func GetProductSeekerByCode(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "operation_product"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	DB.Table("products").Select("products.*, util_unit_measure_types.description as purchase_unit_measure_description").
		Joins("LEFT JOIN util_unit_measure_types ON products.purchase_util_unit_measure_type_id = util_unit_measure_types.id").
		Where("products.company_id = ? AND lower(products.barcode) = lower(?) AND products.state = true", currentUser.CompanyId, product.Barcode).
		Limit(1).Scan(&product)

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product,
	})
}

// ImportProduct function update current Company
func ImportProduct(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Read form fields
	wareHouseIdTemp := c.FormValue("company_ware_house_id")
	companyWareHouseId, err := strconv.ParseUint(wareHouseIdTemp, 10, 32)
	if err != nil {
		return c.JSON(http.StatusConflict, utilities.Response{Message: "unauthorized"})
	}
	updateStock := c.FormValue("update_stock")
	updatePrice := c.FormValue("update_price")
	updateCategory := c.FormValue("update_category")

	// get connection
	DB := provider.GetConnection()
	// defer db.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_import"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Source
	file, err := c.FormFile("excel_file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Validate
	isValid := utilities.ValidateUploadFile(file, 5000, []string{"XLSX", "LSX"})
	if !isValid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: isValid.Message,
		})
	}

	// Destination
	auxDir := "temp/" + file.Filename
	dst, err := os.Create(auxDir)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// ---------------------
	// Read File whit Excel
	// ---------------------
	excel, err := excelize.OpenFile(auxDir)
	if err != nil {
		return err
	}

	// Prepare
	ignoreCols := 1
	counter := 0

	currentWareHouse := models.CompanyWareHouse{}
	if DB.Where("id = ?", companyWareHouseId).First(&currentWareHouse).RowsAffected == 0 {
		return c.JSON(http.StatusConflict, utilities.Response{Message: "unauthorized"})
	}

	err = DB.Transaction(func(TX *gorm.DB) error {
		// Get all the rows in the student.
		rows := excel.GetRows("Sheet0")
		if len(rows) == 0 {
			return errors.New(fmt.Sprintf("no se encontró ningún registro en la hoja 'Sheet0'"))
		}

		for k, row := range rows {
			if k >= ignoreCols {
				// Validate required fields
				//if row[0] == "" || row[1] == "" {
				//    break
				//}
				barcode := strings.TrimSpace(row[0])
				title := strings.TrimSpace(row[1])
				stockMin, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 32)
				if err != nil {
					return err
				}
				stockMax, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 32)
				if err != nil {
					return err
				}
				purchasePrice, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
				if err != nil {
					return err
				}
				unitMeasurePurchase := strings.TrimSpace(row[5])
				factor, err := strconv.ParseFloat(strings.TrimSpace(row[6]), 32)
				if err != nil {
					return err
				}
				unitMeasureSale := strings.TrimSpace(row[7])
				price1, err := strconv.ParseFloat(strings.TrimSpace(row[8]), 64)
				if err != nil {
					return err
				}
				price2, err := strconv.ParseFloat(strings.TrimSpace(row[9]), 64)
				if err != nil {
					return err
				}
				wholeSale2, err := strconv.ParseFloat(strings.TrimSpace(row[10]), 64)
				if err != nil {
					return err
				}
				price3, err := strconv.ParseFloat(strings.TrimSpace(row[11]), 64)
				if err != nil {
					return err
				}
				wholeSale3, err := strconv.ParseFloat(strings.TrimSpace(row[12]), 64)
				if err != nil {
					return err
				}
				price4, err := strconv.ParseFloat(strings.TrimSpace(row[13]), 64)
				if err != nil {
					return err
				}
				wholeSale4, err := strconv.ParseFloat(strings.TrimSpace(row[14]), 64)
				if err != nil {
					return err
				}
				stock, err := strconv.ParseFloat(strings.TrimSpace(row[15]), 64)
				if err != nil {
					return err
				}
				weight, err := strconv.ParseFloat(strings.TrimSpace(row[16]), 32)
				if err != nil {
					return err
				}
				category := strings.TrimSpace(row[17])
				recipie := strings.TrimSpace(row[18])
				bulk := strings.TrimSpace(row[19])

				// Insert product
				product := models.Product{
					Barcode:   barcode,
					Title:     title,
					Weight:    float32(weight),
					Recipe:    recipie == "S",
					Bulk:      bulk == "S",
					CompanyId: currentUser.CompanyId,
				}

				if updatePrice == "true" {
					umPurchaseAux := models.UtilUnitMeasureType{}
					if unitMeasurePurchase == "" {
						unitMeasurePurchase = "NIU"
					}
					if TX.Where("code = ?", unitMeasurePurchase).First(&umPurchaseAux).RowsAffected == 0 {
						return errors.New("unida de medida de compra no reconocida")
					}

					umSaleAux := models.UtilUnitMeasureType{}
					if unitMeasureSale == "" {
						unitMeasureSale = "NIU"
					}
					if TX.Where("code = ?", unitMeasureSale).First(&umSaleAux).RowsAffected == 0 {
						return errors.New("unida de medida de venta no reconocida")
					}

					product.StockMin = float32(stockMin)
					product.StockMax = float32(stockMax)
					product.PurchasePrice = float64(purchasePrice)
					product.PurchaseUtilUnitMeasureTypeId = umPurchaseAux.ID
					product.Factor = float32(factor)
					product.SaleUtilUnitMeasureTypeId = umSaleAux.ID
					product.SalePrice1 = float64(price1)
					product.WholeSale1 = 0
					product.SalePrice2 = float64(price2)
					product.WholeSale2 = float64(wholeSale2)
					product.SalePrice3 = float64(price3)
					product.WholeSale3 = float64(wholeSale3)
					product.SalePrice4 = float64(price4)
					product.WholeSale4 = float64(wholeSale4)
				}

				if updateCategory == "true" {
					categoryAux := models.Category{}
					if category == "" {
						category = "General"
					}
					if TX.Where("name = ?", category).Where("company_id = ?", currentUser.CompanyId).First(&categoryAux).RowsAffected == 0 {
						categoryAux.Name = category
						categoryAux.CompanyId = currentUser.CompanyId
						if err := TX.Create(&categoryAux).Error; err != nil {
							return err
						}
					}
					product.CategoryId = categoryAux.ID
				}

				// Find If exist
				productAux := models.Product{}
				if TX.Where("barcode", product.Barcode).Where("company_id", currentUser.CompanyId).First(&productAux).RowsAffected != 0 {
					product.ID = productAux.ID
					if err := TX.Model(&product).Updates(product).Error; err != nil {
						return err
					}
				} else {
					if err := TX.Create(&product).Error; err != nil {
						return err
					}
				}

				if updateStock == "true" {
					// Update Kardex
					kardexAux := models.Kardex{}
					TX.Where("product_id = ? AND company_ware_house_id = ? AND is_last = true", product.ID, companyWareHouseId).First(&kardexAux)
					TX.Model(&models.Kardex{}).Where("product_id = ?", product.ID).Where("company_ware_house_id = ?", companyWareHouseId).Update("is_last", false)

					kardex := models.Kardex{}
					kardex.DateOfIssue = time.Now()
					kardex.Quantity = stock
					kardex.UnitPrice = product.SalePrice1
					kardex.Total = product.SalePrice1 * stock
					kardex.Origin = "Importar"
					kardex.Destination = currentWareHouse.Description
					kardex.Description = product.Title
					kardex.DocumentDescription = "Importar"
					kardex.UserId = currentUser.ID
					kardex.ProductId = product.ID
					kardex.CompanyWareHouseId = uint(companyWareHouseId)
					kardex.Stock = kardexAux.Stock + stock
					kardex.IsLast = true
					kardex.IsIncome = true
					if err := TX.Create(&kardex).Error; err != nil {
						return err
					}
				}

				counter++
			}
		}

		return nil
	})
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "La importación se realizó exitosamente",
		Data:    counter,
	})
}
