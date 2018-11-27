package productcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/shopping/config"
    "github.com/paulantezana/shopping/models/productmodel"
    "github.com/paulantezana/shopping/utilities"
    "net/http"
)

func GetBrandsPaginate(c echo.Context) error {
    // Get data request
    request := utilities.RequestPaginate{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // Get connection
    db := config.GetConnection()
    defer db.Close()

    // Pagination calculate
    offset := request.Validate()

    // Execute instructions
    var total uint
    brands := make([]productmodel.Brand, 0)

    // Query in database
    if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
        Order("id asc").
        Offset(offset).Limit(request.Limit).Find(&brands).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        brands,
        Total:       total,
        CurrentPage: request.CurrentPage,
        Limit:       request.Limit,
    })
}

func CreateBrand(c echo.Context) error {
    // Get data request
    brand := productmodel.Brand{}
    if err := c.Bind(&brand); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert brands in database
    if err := db.Create(&brand).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    brand.ID,
        Message: fmt.Sprintf("La marca %s se registro correctamente", brand.Name),
    })
}

func UpdateBrand(c echo.Context) error {
    // Get data request
    brand := productmodel.Brand{}
    if err := c.Bind(&brand); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Update brand in database
    rows := db.Model(&brand).Update(brand).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", brand.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    brand.ID,
        Message: fmt.Sprintf("Los datos del la marca %s se actualizaron correctamente", brand.Name),
    })
}

func DeleteBrand(c echo.Context) error {
    // Get data request
    brand := productmodel.Brand{}
    if err := c.Bind(&brand); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Delete brand in database
    if err := db.Delete(&brand).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    brand.ID,
        Message: fmt.Sprintf("La marca %s se elimino correctamente", brand.Name),
    })
}

func MultipleDeleteBrand(c echo.Context) error {
    // Get data request
    deleteRequest := utilities.DeleteRequest{}
    if err := c.Bind(&deleteRequest); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    tx := db.Begin()
    for _, value := range deleteRequest.Ids {
        brand := productmodel.Brand{
            ID: value,
        }

        // Delete brand in database
        if err := tx.Delete(&brand).Error; err != nil {
            tx.Rollback()
            return c.JSON(http.StatusOK, utilities.Response{
                Success: false,
                Message: fmt.Sprintf("%s", err),
            })
        }
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: fmt.Sprintf("Sel eliminaron %d registros", len(deleteRequest.Ids)),
    })
}
