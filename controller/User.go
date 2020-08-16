package controller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

// CreateUser function create new user
func RegisterUser(c echo.Context) error {
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

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Insert user in database
	user.State = true
	user.UserRoleId = 1
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

// forgotSearchEmailTemplate struct template
type forgotSearchEmailTemplate struct {
	UserName  string `json:"user_name" gorm:"type:varchar(64); not null"`
	Email     string `json:"email" gorm:"type:varchar(64); not null"`
	SecretKey string `json:"secret_key"`
}

// ForgotSearch function forgot user search
func ForgotSearch(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validations
	if err := DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("Tu búsqueda no arrojó ningún resultado. Vuelve a intentarlo con otros datos."),
		})
	}

	// Generate key validation
	userForgot := models.UserForgot{}
	rand.Seed(time.Now().UnixNano())
	secretKey := rand.Int31()
	userForgot.SecretKey = fmt.Sprint(secretKey)[:7]
	userForgot.UserId = user.ID
	userForgot.CreatedUserId = user.ID

	// Update key
	userForgotUpdate := models.UserForgot{
		UserId: user.ID,
	}
	if err := DB.Model(&userForgotUpdate).UpdateColumn("used", true).UpdateColumn("secret_key", "").UpdateColumn("updated_user_id", user.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Update database
	if err := DB.Create(&userForgot).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Query Database Get Settings
	con := models.Company{}
	DB.First(&con)

	// SEND EMAIL get html template
	t, err := template.ParseFiles("templates/email.html")
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// SEND EMAIL new buffer
	forgotTemplate := forgotSearchEmailTemplate{}
	forgotTemplate.UserName = user.UserName
	forgotTemplate.SecretKey = userForgot.SecretKey
	forgotTemplate.Email = user.Email
	buf := new(bytes.Buffer)
	err = t.Execute(buf, forgotTemplate)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// SEND EMAIL
	err = provider.SendEmail(
		con.CommercialReason,
		user.Email,
		fmt.Sprint(secretKey)+" es el código de recuperación de tu cuenta",
		buf.String(),
	)
	if err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response success api service
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user,
	})
}

// ForgotValidate function forgot user validate
func ForgotValidate(c echo.Context) error {
	userForgot := models.UserForgot{}
	if err := c.Bind(&userForgot); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}

	// get connection
	db := provider.GetConnection()
	defer db.Close()

	// Validations
	if err := db.Where("user_id = ? AND secret_key = ? AND used = ?", userForgot.UserId, userForgot.SecretKey, false).Last(&userForgot).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El número %s que ingresaste no coincide con tu código de seguridad. Vuelve a intentarlo", userForgot.SecretKey),
		})
	}

	// Response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    userForgot,
	})
}

// ForgotChange function forgot password change
func ForgotChange(c echo.Context) error {
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
	currentUser := models.User{}
	if err := db.Where("id = ?", user.ID).First(&currentUser).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontro ningun registro con el id %d", user.ID),
		})
	}

	// Encrypted old password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Update
	if err := db.Model(&user).Update(user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Update key
	userForgot := models.UserForgot{
		UserId: user.ID,
	}
	if err := db.Model(&userForgot).UpdateColumn("used", true).UpdateColumn("secret_key", "").UpdateColumn("updated_user_id", user.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Response data
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user,
		Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", currentUser.UserName),
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

// configResponse struct
type configResponse struct {
	AdminMenu []models.AppAuthorization `json:"admin_menu"`
}

// GetMenuAdminByUserId get admin menu
func GetMenuAdminByUserId(c echo.Context) error {
	tUser := c.Get("user").(*jwt.Token)
	claims := tUser.Claims.(*utilities.Claim)
	currentUser := claims.User

	db := provider.GetConnection()
	defer db.Close()

	appAuthorizations := make([]models.AppAuthorization, 0)
	if err := db.Table("user_role_authorizations").Select("app_authorizations.*").
		Joins("INNER JOIN app_authorizations on app_authorizations.id = user_role_authorizations.app_authorization_id").
		Where("user_role_authorizations.user_role_id = ? AND user_role_authorizations.state = true AND app_authorizations.state = true", currentUser.UserRoleId).
		Scan(&appAuthorizations).
		Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data: configResponse{
			AdminMenu: appAuthorizations,
		},
	})
}

// PaginateUser function get all users
func PaginateUser(c echo.Context) error {
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
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

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

	// Get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Execute instructions
	if err := DB.First(&user, user.ID).Error; err != nil {
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

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Insert user in database
	user.CreatedUserId = currentUser.ID
	if err := DB.Create(&user).Error; err != nil {
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation user exist
	aux := models.User{ID: user.ID}
	if DB.First(&aux).RecordNotFound() {
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
	if err := DB.Model(&user).Update(user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}
	if !user.State {
		if err := DB.Model(user).UpdateColumn("state", false).Error; err != nil {
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

// UpdateUser function update current user
func ChangePasswordUser(c echo.Context) error {
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validation user exist
	aux := models.User{ID: user.ID}
	if DB.First(&aux).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Update user in database
	if err := DB.Model(user).UpdateColumn("password", user.Password).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "La contraseña se cambio correctamente",
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Validate
	aux := models.User{}
	DB.First(&aux, user.ID)
	if aux.Freeze {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El usuario %s está protegido por el sistema y no se permite deshavilitar", user.UserName),
		})
	}

	// Update
	if !user.State {
		if err := DB.Model(user).UpdateColumn("state", false).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := DB.Model(user).UpdateColumn("state", true).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
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

// UploadLogoCompany function update current Company
func UploadAvatarUser(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// get connection
	DB := provider.GetConnection()
	defer DB.Close()

	// Source
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Validate
	isValid := utilities.ValidateUploadFile(file, 100, []string{"JPG", "PNG", "SVG"})
	if !isValid.Success {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: isValid.Message,
		})
	}

	// Destination
	ccc := sha256.Sum256([]byte(strconv.Itoa(int(currentUser.ID))))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	logoSRC := "static/user/" + name
	dst, err := os.Create(logoSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	currentUser.Avatar = "/" + logoSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update Company in database
	currentUser.UpdatedUserId = currentUser.ID
	if err := DB.Model(&currentUser).Update(currentUser).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El avatar se subio correctamente",
		Data:    currentUser.ID,
	})
}

type salePointByUserResponse struct {
	ID             uint   `json:"id"`
	Description    string `json:"description"`
	CompanyLocalId uint   `json:"company_local_id"`
	State          bool   `json:"state"`
	AuthId         uint   `json:"auth_id"`
	AuthState      bool   `json:"auth_state"`
}

// SaveSalePointByUserId function update current user
func LoadSalePointByUserId(c echo.Context) error {
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update
	salePoints := make([]salePointByUserResponse, 0)
	if err := DB.Raw("SELECT csp.*,  usp.id as auth_id, usp.state as auth_state  FROM company_sale_points as csp "+
		" LEFT JOIN ( "+
		" SELECT id, company_sale_point_id, state FROM user_sale_points WHERE user_sale_points.user_id = ? "+
		" ) as usp ON csp.id = usp.company_sale_point_id "+
		" WHERE csp.state = true ", user.ID).Scan(&salePoints).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El usuario se actualizó correctamente",
		Data:    salePoints,
	})
}

// SaveSalePointByUserId function update current user
func SaveSalePointByUserId(c echo.Context) error {
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
	DB := provider.GetConnection()
	defer DB.Close()

	// Validate Auth
	if err := validateIsAuthorized(DB, currentUser.UserRoleId, "setting_user"); err != nil {
		return c.JSON(http.StatusForbidden, utilities.Response{Message: "unauthorized"})
	}

	// Update
    for _, point := range user.UserSalePoints {
        if point.ID == 0 {
            point.CreatedUserId = currentUser.ID
            if err := DB.Create(&point).Error; err != nil {
                return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
            }
        } else {
            if err := DB.Model(point).UpdateColumn("state", point.State).UpdateColumn("updated_user_id", currentUser.ID).Error; err != nil {
                return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
            }
        }
    }

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: "El usuario se actualizó correctamente",
		Data:    user.ID,
	})
}
