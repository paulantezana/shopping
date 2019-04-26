package companyController

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/shopping/config"
    "github.com/paulantezana/shopping/models"
    "github.com/paulantezana/shopping/utilities"
    "html/template"
    "io"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
)

type loginResponse struct {
    User     interface{}   `json:"user"`
    Token    interface{}   `json:"token"`
}

// Login login app
func Login(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Hash password
    cc := sha256.Sum256([]byte(personal.Password))
    pwd := fmt.Sprintf("%x", cc)

    // Query database
    if DB.Where("user = ? and password = ?", personal.User, pwd).First(&personal).RecordNotFound() {
        if DB.Where("email = ? and password = ?", personal.User, pwd).First(&personal).RecordNotFound() {
            return c.JSON(http.StatusOK, utilities.Response{
                Message: "El nombre de usuario o contraseña es incorecta",
            })
        }
    }

    // Check state user
    //if !personal.State {
    //    return c.NoContent(http.StatusForbidden)
    //}

    // Prepare response data
    personal.Password = ""
    personal.Key = ""

    // Insert new Session
    //session := models.Session{
    //    UserID:       user.ID,
    //    LastActivity: time.Now(),
    //}
    //if err := DB.Create(&session).Error; err != nil {
    //    return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    //}

    // get token key
    token := utilities.GenerateJWT(personal)

    // Login success
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: fmt.Sprintf("Bienvenido al sistema %s", personal.User),
        Data: loginResponse{
            User:  personal,
            Token: token,
        },
    })
}

// Login password login check
func LoginPasswordCheck(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Hash password
    cc := sha256.Sum256([]byte(personal.Password))
    pwd := fmt.Sprintf("%x", cc)

    // Hash password
    if DB.Where("id = ? and password = ?", personal.ID, pwd).First(&personal).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("Contraseña incorrecta")})
    }

    // Check state user
    //if !user.State {
    //    return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("No tiene permisos para realizar ningún tipo de acción.")})
    //}

    // Prepare response data
    personal.Password = ""
    personal.Key = ""

    // Login success
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: fmt.Sprintf("La verificación de su contraseña fue exitosamente. usuario: %s", personal.User),
        Data:    personal.ID,
    })
}

// ForgotSearch function forgot user search
func ForgotSearch(c echo.Context) error {
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // Get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Validations
    if err := DB.Where("email = ?", personal.Email).First(&personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("Tu búsqueda no arrojó ningún resultado. Vuelve a intentarlo con otros datos."),
        })
    }

    // Generate key validation
    key := (int)(rand.Float32() * 10000000)
    personal.Key = fmt.Sprint(key)

    // Update database
    if err := DB.Model(&personal).Update(personal).Error; err != nil {
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
    buf := new(bytes.Buffer)
    err = t.Execute(buf, personal)
    if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // SEND EMAIL
    err = config.SendEmail(
        con.CompanyName,
        personal.Email,
        fmt.Sprint(key)+" es el código de recuperación de tu cuenta",
        buf.String(),
    )
    if err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Response success api service
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    personal.ID,
    })
}

// ForgotValidate function forgot user validate
func ForgotValidate(c echo.Context) error {
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Validations
    if err := DB.Where("id = ? AND key = ?", personal.ID, personal.Key).First(&personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("El número %s que ingresaste no coincide con tu código de seguridad. Vuelve a intentarlo", personal.Key),
        })
    }

    // Response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    personal.ID,
    })
}

// ForgotChange function forgot password change
func ForgotChange(c echo.Context) error {
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Validate
    currentUser := models.Personal{}
    if err := DB.Where("id = ?", personal.ID).First(&currentUser).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se encontro ningun registro con el id %d", personal.ID),
        })
    }

    // Encrypted old password
    cc := sha256.Sum256([]byte(personal.Password))
    pwd := fmt.Sprintf("%x", cc)
    personal.Password = pwd

    // Update
    if err := DB.Model(&personal).Update(personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Update key
    if err := DB.Model(&personal).UpdateColumn("key", "").Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Response data
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    personal.ID,
        Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", currentUser.User),
    })
}

// GetUsers function get all users
//func GetUsers(c echo.Context) error {
//    // Get user token authenticate
//    user := c.Get("user").(*jwt.Token)
//    claims := user.Claims.(*utilities.Claim)
//    currentUser := claims.User
//
//    // Get data request
//    request := utilities.Request{}
//    if err := c.Bind(&request); err != nil {
//        return c.JSON(http.StatusBadRequest, utilities.Response{
//            Message: "La estructura no es válida",
//        })
//    }
//
//    // Get connection
//    DB := config.GetConnection()
//    defer DB.Close()
//
//    // Pagination calculate
//    offset := request.Validate()
//
//    // Check the number of matches
//    var total uint
//    personals := make([]models.Personal, 0)
//
//    // Find users
//    if err := DB.Where("user_name LIKE ? AND role_id >= ?", "%"+request.Search+"%", currentUser.RoleID).
//        Order("id desc").Offset(offset).Limit(request.Limit).Find(&personals).
//        Offset(-1).Limit(-1).Count(&total).Error; err != nil {
//        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
//    }
//
//    // find
//    for i := range users {
//        student := models.Student{}
//        DB.First(&student, models.Student{UserID: users[i].ID})
//        if student.ID >= 1 {
//            users[i].UserName = student.FullName
//        } else {
//            teacher := models.Teacher{}
//            DB.First(&teacher, models.Teacher{UserID: users[i].ID})
//            if teacher.ID >= 1 {
//                users[i].UserName = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
//            }
//        }
//    }
//
//    // Return response
//    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
//        Success:     true,
//        Data:        users,
//        Total:       total,
//        CurrentPage: request.CurrentPage,
//        Limit:       request.Limit,
//    })
//}

// GetUsers function get all users
type searchUsersResponse struct {
    ID       uint   `json:"id"`
    UserName string `json:"user_name"`
    Avatar   string `json:"avatar"`
}

//func SearchUsers(c echo.Context) error {
//    // Get user token authenticate
//    user := c.Get("user").(*jwt.Token)
//    claims := user.Claims.(*utilities.Claim)
//    currentUser := claims.User
//
//    // Get data request
//    request := utilities.Request{}
//    if err := c.Bind(&request); err != nil {
//        return c.JSON(http.StatusBadRequest, utilities.Response{
//            Message: "La estructura no es válida",
//        })
//    }
//
//    // Get connection
//    DB := config.GetConnection()
//    defer DB.Close()
//
//    // Find users
//    users := make([]searchUsersResponse, 0)
//    if err := DB.Raw("SELECT id, user_name, avatar FROM users "+
//        "WHERE lower(user_name) LIKE lower(?) "+
//        "OR id IN (SELECT user_id FROM teachers WHERE lower(first_name) LIKE lower(?) LIMIT 20) "+
//        "OR id IN (SELECT user_id FROM students WHERE lower(full_name) LIKE lower(?) LIMIT 20) "+
//        "LIMIT 30", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%").Scan(&users).Error; err != nil {
//        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
//    }
//
//    // Queries
//    newUsers := make([]searchUsersResponse, 0)
//    for i := range users {
//        // ignore current user
//        if users[i].ID == currentUser.ID {
//            break
//        }
//
//        // search user
//        nUser := searchUsersResponse{
//            ID:       users[i].ID,
//            UserName: users[i].UserName,
//            Avatar:   users[i].Avatar,
//        }
//
//        // Query current student Name
//        student := models.Student{}
//        DB.First(&student, models.Student{UserID: nUser.ID})
//        if student.ID >= 1 {
//            nUser.UserName = student.FullName
//        } else {
//            teacher := models.Teacher{}
//            DB.First(&teacher, models.Teacher{UserID: nUser.ID})
//            if teacher.ID >= 1 {
//                nUser.UserName = fmt.Sprintf("%s %s", teacher.FirstName, teacher.LastName)
//            }
//        }
//
//        // append child
//        newUsers = append(newUsers, nUser)
//    }
//
//    // Return response
//    return c.JSON(http.StatusCreated, utilities.Response{
//        Success: true,
//        Data:    newUsers,
//    })
//}

// GetUserByID function get user by id
func GetUserByID(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // Get connection
    db := config.GetConnection()
    defer db.Close()

    // Execute instructions
    if err := db.First(&personal, personal.ID).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    personal,
    })
}

// CreateUser function create new user
func CreateUser(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // Default empty values
    //if user.RoleID == 0 {
    //    user.RoleID = 6
    //}

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Hash password
    cc := sha256.Sum256([]byte(personal.Password))
    pwd := fmt.Sprintf("%x", cc)
    personal.Password = pwd

    // Insert user in database
    if err := DB.Create(&personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success: true,
        Data:    personal.ID,
        Message: fmt.Sprintf("El usuario %s se registro exitosamente", personal.User),
    })
}

// UpdateUser function update current user
func UpdateUser(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Validation user exist
    aux := models.Personal{ID: personal.ID}
    if db.First(&aux).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se encontró el registro con id %d", personal.ID),
        })
    }

    // Update user in database
    if err := db.Model(&personal).Update(personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }
    //if !personal.State {
    //    if err := db.Model(personal).UpdateColumn("state", false).Error; err != nil {
    //        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    //    }
    //}

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    personal.ID,
    })
}

// DeleteUser function delete user by id
func DeleteUser(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Validate
    DB.First(&personal, personal.ID)
    if personal.Freeze {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("El usuario %s está protegido por el sistema y no se permite eliminar", personal.User),
        })
    }

    // Delete user in database
    if err := DB.Delete(&personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    personal.ID,
        Message: fmt.Sprintf("El usuario %s, se elimino correctamente", personal.User),
    })
}

// UploadAvatarUser function upload avatar user
func UploadAvatarUser(c echo.Context) error {
    // Read form fields
    idUser := c.FormValue("id")
    personal := models.Personal{}

    // get connection
    DB := config.GetConnection()
    defer DB.Close()

    // Validation user exist
    if DB.First(&personal, "id = ?", idUser).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se encontró el registro con id %d", personal.ID),
        })
    }

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

    // Destination
    ccc := sha256.Sum256([]byte(string(personal.ID)))
    name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
    avatarSRC := "static/profiles/" + name
    dst, err := os.Create(avatarSRC)
    if err != nil {
        return err
    }
    defer dst.Close()
    personal.Avatar = avatarSRC

    // Copy
    if _, err = io.Copy(dst, src); err != nil {
        return err
    }

    // Update database user
    if err := DB.Model(&personal).Update(personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    // Return response
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Data:    personal,
        Message: fmt.Sprintf("El avatar del usuario %s, se subió correctamente", personal.User),
    })
}

// ResetPasswordUser function reset password
func ResetPasswordUser(c echo.Context) error {
    // Get data request
    personal := models.Personal{}
    if err := c.Bind(&personal); err != nil {
        return c.JSON(http.StatusBadRequest, utilities.Response{
            Message: "La estructura no es válida",
        })
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Validation user exist
    if db.First(&personal, "id = ?", personal.ID).RecordNotFound() {
        return c.JSON(http.StatusOK, utilities.Response{
            Message: fmt.Sprintf("No se encontró el registro con id %d", personal.ID),
        })
    }

    // Set new password
    cc := sha256.Sum256([]byte(fmt.Sprintf("%d%s", personal.ID, personal.User)))
    pwd := fmt.Sprintf("%x", cc)
    personal.Password = pwd

    // Update user in database
    if err := db.Model(&personal).Update(personal).Error; err != nil {
        return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
    }

    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: fmt.Sprintf("La contraseña del usuario se cambio exitosamente. ahora su numevacontraseña es %d%s", personal.ID, personal.User),
    })
}

// ChangePasswordUser function change password user
//func ChangePasswordUser(c echo.Context) error {
//    // Get data request
//    user := models.Personal{}
//    if err := c.Bind(&user); err != nil {
//        return c.JSON(http.StatusBadRequest, utilities.Response{
//            Message: "La estructura no es válida",
//        })
//    }
//
//    // get connection
//    db := config.GetConnection()
//    defer db.Close()
//
//    // Validation user exist
//    aux := models.Personal{ID: user.ID}
//    if db.First(&aux, "id = ?", aux.ID).RecordNotFound() {
//        return c.JSON(http.StatusOK, utilities.Response{
//            Message: fmt.Sprintf("No se encontró el registro con id %d", aux.ID),
//        })
//    }
//
//    // Change password
//    if len(user.Password) > 0 {
//        // Validate empty length old password
//        if len(user.OldPassword) == 0 {
//            return c.JSON(http.StatusOK, utilities.Response{
//                Message: "Ingrese la contraseña antigua",
//            })
//        }
//
//        // Hash old password
//        ccc := sha256.Sum256([]byte(user.OldPassword))
//        old := fmt.Sprintf("%x", ccc)
//
//        // validate old password
//        if db.Where("password = ?", old).First(&aux).RecordNotFound() {
//            return c.JSON(http.StatusOK, utilities.Response{
//                Message: "La contraseña antigua es incorrecta",
//            })
//        }
//
//        // Set and hash new password
//        cc := sha256.Sum256([]byte(user.Password))
//        pwd := fmt.Sprintf("%x", cc)
//        user.Password = pwd
//    }
//
//    // Update user in database
//    if err := db.Model(&user).Update(user).Error; err != nil {
//        return err
//    }
//
//    return c.JSON(http.StatusOK, utilities.Response{
//        Success: true,
//        Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", aux.User),
//    })
//}
