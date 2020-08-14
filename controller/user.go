package controller

import (
    "bytes"
    "crypto/sha256"
	"fmt"
    "html/template"
    "math/rand"
    "net/http"
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

type forgotSearchEmailTemplate struct {
    UserName string `json:"user_name" gorm:"type:varchar(64); not null"`
    Email    string `json:"email" gorm:"type:varchar(64); not null"`
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
    userForgot.CreatedUserId =  user.ID

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
        Data: userForgot,
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
    db := provider.GetConnection()
    defer db.Close()

    // Validation user exist
    aux := models.User{ID: user.ID}
    if db.First(&aux).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
        })
    }

    // Hash password
    cc := sha256.Sum256([]byte(user.Password))
    pwd := fmt.Sprintf("%x", cc)
    user.Password = pwd

    // Update user in database
    if err := db.Model(user).UpdateColumn("password", user.Password).UpdateColumn("updated_user_id",currentUser.ID).Error; err != nil {
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
