package productcontroller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/shopping/config"
    "github.com/paulantezana/shopping/models/productmodel"
    "github.com/paulantezana/shopping/utilities"
    "net/http"
)

func GetCategoriesPaginate(c echo.Context) error {
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
    categories := make([]productmodel.Category, 0)

    // Query in database
    if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
        Order("id asc").
        Offset(offset).Limit(request.Limit).Find(&categories).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return err
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        categories,
        Total:       total,
        CurrentPage: request.CurrentPage,
        Limit:       request.Limit,
    })
}

func CreateCategory(c echo.Context) error {
    // Get data request
    category := productmodel.Category{}
    if err := c.Bind(&category); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert categories in database
    if err := db.Create(&category).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    category.ID,
        Message: fmt.Sprintf("La categoria %s se registro correctamente", category.Name),
    })
}

func UpdateCategory(c echo.Context) error {
    // Get data request
    category := productmodel.Category{}
    if err := c.Bind(&category); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Update category in database
    rows := db.Model(&category).Update(category).RowsAffected
    if rows == 0 {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", category.ID),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    category.ID,
        Message: fmt.Sprintf("Los datos del la categoria %s se actualizaron correctamente", category.Name),
    })
}

func DeleteCategory(c echo.Context) error {
    // Get data request
    category := productmodel.Category{}
    if err := c.Bind(&category); err != nil {
        return err
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Delete category in database
    if err := db.Delete(&category).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Success: false,
            Message: fmt.Sprintf("%s", err),
        })
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    category.ID,
        Message: fmt.Sprintf("La categoria %s se elimino correctamente", category.Name),
    })
}

func MultipleDeleteCategory(c echo.Context) error {
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
        category := productmodel.Category{
            ID: value,
        }

        // Delete category in database
        if err := tx.Delete(&category).Error; err != nil {
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