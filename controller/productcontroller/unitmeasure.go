package productcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/shopping/config"
    "github.com/paulantezana/shopping/models"
    "github.com/paulantezana/shopping/utilities"
    "net/http"
)

func GetUnitMeasuresPaginate(c echo.Context) error {
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
    unitMeasures := make([]models.UnitMeasure, 0)

    // Query in database
    if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
        Order("id asc").
        Offset(offset).Limit(request.Limit).Find(&unitMeasures).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        unitMeasures,
        Total:       total,
        CurrentPage: request.CurrentPage,
        Limit:       request.Limit,
    })
}

func CreateUnitMeasure(c echo.Context) error {
    // Get data request
    unitMeasure := models.UnitMeasure{}
    if err := c.Bind(&unitMeasure); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert unitMeasures in database
    if err := db.Create(&unitMeasure).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    unitMeasure.ID,
        Message: fmt.Sprintf("La unidad de medida %s se registro correctamente", unitMeasure.Name),
    })
}

func UpdateUnitMeasure(c echo.Context) error {
    // Get data request
    unitMeasure := models.UnitMeasure{}
    if err := c.Bind(&unitMeasure); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Update unitMeasure in database
    rows := db.Model(&unitMeasure).Update(unitMeasure).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", unitMeasure.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    unitMeasure.ID,
        Message: fmt.Sprintf("Los datos del la unidad de medida %s se actualizaron correctamente", unitMeasure.Name),
    })
}

func DeleteUnitMeasure(c echo.Context) error {
    // Get data request
    unitMeasure := models.UnitMeasure{}
    if err := c.Bind(&unitMeasure); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Delete unitMeasure in database
    if err := db.Delete(&unitMeasure).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    unitMeasure.ID,
        Message: fmt.Sprintf("La unidad de medida %s se elimino correctamente", unitMeasure.Name),
    })
}

func MultipleDeleteUnitMeasure(c echo.Context) error {
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
        unitMeasure := models.UnitMeasure{
            ID: value,
        }

        // Delete unitMeasure in database
        if err := tx.Delete(&unitMeasure).Error; err != nil {
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
