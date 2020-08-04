package controller

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// Login login app
func Login(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)

	// Validate user and email
	if !DB.Where("user_name = ? and password = ?", user.UserName, pwd).First(&user).RecordNotFound() {
		if !DB.Where("email = ? and password = ?", user.UserName, pwd).First(&user).RecordNotFound() {
			return c.JSON(http.StatusOK, utilities.Response{
				Message: "El nombre de usuario o contraseña es incorecta",
			})
		}
	}

	// Check state user
	if !user.State {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: "No autorizado",
		})
	}

	// Prepare response data
	user.Password = ""

	// get token key
	token := utilities.GenerateJWT(user)

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", user.UserName),
		Data:    token,
	})
}

// GetUserByToken function get user by token
func GetUserByToken(c echo.Context) error {
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	db := provider.GetConnection()
	defer db.Close()

	user := models.User{}
	if err := db.First(&user, currentUser.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	user.Password = ""

	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user,
	})
}

// PaginateUser function get all users
func PaginateUser(c echo.Context) error {
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

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	users := make([]models.User, 0)

	// Find users
	if err := DB.Where("lower(user_name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").Offset(offset).Limit(request.PageSize).Find(&users).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     users,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetUserByID function get user by id
func GetUserByID(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	db := provider.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user,
	})
}

// CreateUser function create new user
func CreateUser(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Default empty values
	// if user.UserRoleID == 0 {
	// 	user.UserRoleID = 6
	// }

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Insert user in database
	user.CreatedUserId = currentUser.ID
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user.ID,
		Message: fmt.Sprintf("El usuario %s se registro exitosamente", user.UserName),
	})
}

// UpdateUser function update current user
func UpdateUser(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validation user exist
	aux := models.User{ID: user.ID}
	if db.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}
	if !user.State {
		if aux.Freeze {
			return c.JSON(http.StatusOK, utilities.Response{
				Message: fmt.Sprintf("El usuario %s está protegido por el sistema y no se permite deshavilitar", user.UserName),
			})
		}
	}

	// Update user in database
	user.UpdatedUserId = currentUser.ID
	if err := db.Model(&user).Update(user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !user.State {
		if err := db.Model(user).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El usuario se actualizó correctamente",
		Data:    user.ID,
	})
}

// UpdateStateUser function update current user
func UpdateStateUser(c echo.Context) error {
    // Get user token authenticate
    tUser := c.Get("user").(*jwt.Token)
    claims := tUser.Claims.(*utilities.Claim)
    currentUser := claims.User

	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validate
	aux := models.User{}
	db.First(&aux, user.ID)
	if aux.Freeze {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El usuario %s está protegido por el sistema y no se permite deshavilitar", user.UserName),
		})
	}

	// Update
	if !user.State {
		if err := db.Model(user).UpdateColumn("state", false).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := db.Model(user).UpdateColumn("state", true).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El usuario se actualizó correctamente",
		Data:    user.ID,
	})
}
