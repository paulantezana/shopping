package controller

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
	"github.com/paulantezana/shopping/utilities"
)

// PaginateUserRole function get all userRoles
func PaginateUserRole(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Pagination calculate
	offset := request.Validate()

	// Check the number of matches
	var total uint
	userRoles := make([]models.UserRole, 0)

	// Find userRoles
	if err := DB.Where("lower(description) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").Offset(offset).Limit(request.PageSize).Find(&userRoles).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:  true,
		Data:     userRoles,
		Total:    total,
		Current:  request.CurrentPage,
		PageSize: request.PageSize,
	})
}

// GetAllUserRole function get all userRoles
func GetAllUserRole(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Find userRoles
	userRoles := make([]models.UserRole, 0)
	if err := DB.Where("state = true").Find(&userRoles).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    userRoles,
	})
}

// GetUserRoleByID function get userRole by id
func GetUserRoleByID(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	userRole := models.UserRole{}
	if err := c.Bind(&userRole); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&userRole, userRole.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    userRole,
	})
}

// CreateUserRole function create new userRole
func CreateUserRole(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	userRole := models.UserRole{}
	if err := c.Bind(&userRole); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Insert userRole in database
	userRole.CreatedUserId = currentUser.ID
	if err := DB.Create(&userRole).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    userRole.ID,
		Message: fmt.Sprintf("El rol %s se registro exitosamente", userRole.Description),
	})
}

// UpdateUserRole function update current userRole
func UpdateUserRole(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	userRole := models.UserRole{}
	if err := c.Bind(&userRole); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation userRole exist
	aux := models.UserRole{ID: userRole.ID}
	if DB.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", userRole.ID),
		})
	}

	// Update userRole in database
	userRole.UpdatedUserId = currentUser.ID
	if err := DB.Model(&userRole).Update(userRole).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !userRole.State {
		if err := DB.Model(userRole).UpdateColumn("state", false).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El rol se actualizó correctamente",
		Data:    userRole.ID,
	})
}

// UpdateStateUserRole function update current userRole
func UpdateStateUserRole(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	userRole := models.UserRole{}
	if err := c.Bind(&userRole); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update userRole in database
	if !userRole.State {
		if err := DB.Model(userRole).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(userRole).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El rol se actualizó correctamente",
		Data:    userRole.ID,
	})
}

// getAppAuthorizationUserRoleResponse --
type getAppAuthorizationUserRoleResponse struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Key         string `json:"key" gorm:"not null; type:varchar(64)"`
	Type        uint   `json:"type" gorm:"default: true"` // 0 = menu, 1 = action
	Action      string `json:"action" gorm:"not null; type:varchar(128)"`
	Description string `json:"description" gorm:"not null; type:varchar(255)"`
	ParentId    uint   `json:"parent_id" gorm:"default: 0"`
	State       bool   `json:"state" gorm:"default: true"`
	AuthId      uint   `json:"auth_id"`
	AuthState   bool   `json:"auth_state"`
}

// UpdateStateUserRole function update current userRole
func GetAppAuthorizationByUserRole(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	userRole := models.UserRole{}
	if err := c.Bind(&userRole); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update userRole in database
	authorizations := make([]getAppAuthorizationUserRoleResponse, 0)
	if err := DB.Raw("SElECT au.*, ura.id as auth_id, ura.state as auth_state FROM app_authorizations au "+
		" LEFT JOIN ( "+
		" SELECT id, app_authorization_id, state FROM user_role_authorizations WHERE user_role_authorizations.user_role_id = ? "+
		" ) as ura on au.id = ura.app_authorization_id "+
		" WHERE au.state = true ", userRole.ID).Scan(&authorizations).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    authorizations,
	})
}

// UpdateUserRoleAppAuthorization function update current userRole
func UpdateUserRoleAppAuthorization(c echo.Context) error {
	// Get user token authenticate
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	userRole := models.UserRole{}
	if err := c.Bind(&userRole); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user_rol"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update userRole in database
	userRole.UpdatedUserId = currentUser.ID
	if err := DB.Model(&userRole).Update(userRole).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Los permisos del rol %s se actualizaron correctamente ", userRole.Description),
		Data:    userRole.ID,
	})
}
