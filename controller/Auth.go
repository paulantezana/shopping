package controller

import (
	"errors"
	"gorm.io/gorm"

	"github.com/paulantezana/shopping/models"
)

func validateIsAuthorized(DB *gorm.DB, userRoleId uint, key string) error {
	appAuthorization := models.AppAuthorization{}
	if err := DB.Table("user_role_authorizations").Select("app_authorizations.*").
		Joins("INNER JOIN app_authorizations on app_authorizations.id = user_role_authorizations.app_authorization_id").
		Where("user_role_authorizations.user_role_id = ? AND user_role_authorizations.state = true "+
			"AND app_authorizations.state = true AND app_authorizations.key = ?", userRoleId, key).Limit(1).Scan(&appAuthorization).Error; err != nil {
		return err
	}

	if appAuthorization.ID == 0 {
		return errors.New("Denied")
	}
	return nil
}
