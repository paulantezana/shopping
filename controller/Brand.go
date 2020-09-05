package controller

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "github.com/paulantezana/shopping/models"
    "github.com/paulantezana/shopping/provider"
    "github.com/paulantezana/shopping/utilities"
    "net/http"
)

// PaginateBrand function get all brands
func PaginateBrand(c echo.Context) error {
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
    defer DB.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_brand"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Pagination calculate
    offset := request.Validate()

    // Check the number of matches
    var total uint
    brands := make([]models.Brand, 0)

    // Find users
    if err := DB.Where("company_id = ? AND lower(name) LIKE lower(?)", currentUser.CompanyId, "%"+request.Search+"%").
        Order("id desc").Offset(offset).Limit(request.PageSize).Find(&brands).
        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:  true,
        Data:     brands,
        Total:    total,
        Current:  request.CurrentPage,
        PageSize: request.PageSize,
    })
}

// GetBrandByID function get brand by id
func GetBrandByID(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    brand := models.Brand{}
    if err := c.Bind(&brand); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // Get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_brand"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Execute instructions
    if err := DB.First(&brand, brand.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    brand,
    })
}

// CreateBrand function create new brand
func CreateBrand(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    brand := models.Brand{}
    if err := c.Bind(&brand); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_brand"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Insert brand in database
    brand.CreatedUserId = currentUser.ID
    brand.CompanyId = currentUser.CompanyId
    if err := DB.Create(&brand).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    brand.ID,
        Message: fmt.Sprintf("La marca %s se registro exitosamente", brand.Name),
    })
}

// UpdateBrand function update current brand
func UpdateBrand(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    brand := models.Brand{}
    if err := c.Bind(&brand); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_brand"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Validation brand exist
    aux := models.Brand{ID: brand.ID}
    if DB.First(&aux).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se encontró el registro con id %d", brand.ID),
        })
    }

    // Update brand in database
    brand.UpdatedUserId = currentUser.ID
    if err := DB.Model(&brand).Update(brand).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }
    if !brand.State {
        if err := DB.Model(brand).UpdateColumn("state", false).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: "La marca se actualizó correctamente",
        Data:    brand.ID,
    })
}

// UpdateStateBrand function update current brand
func UpdateStateBrand(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

    // Get data request
    brand := models.Brand{}
    if err := c.Bind(&brand); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := provider.GetConnection()
    defer DB.Close()

    // Validate Auth
    if err := validateIsAuthorized(DB, currentUser.UserRoleId, "maintenance_brand"); err != nil {
        return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
    }

    // Update brand in database
    if !brand.State {
        if err := DB.Model(brand).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
    } else {
        if err := DB.Model(brand).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
            return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
        }
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: "La marca se actualizó correctamente",
        Data:    brand.ID,
    })
}