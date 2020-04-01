package migration

import (
	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
)

// Migrate function
func Migrate() {
	db := provider.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.UtilAdditionalLegendType{},
		&models.UtilCatAffectationIgvType{},
		&models.UtilCreditDebitType{},
		&models.UtilCurrencyType{},
		&models.UtilDocumentType{},
		&models.UtilGeographicalLocation{},
		&models.UtilIdentityDocumentType{},
		&models.UtilOperationType{},
		&models.UtilPerceptionType{},
		&models.UtilProductType{},
		&models.UtilSubjectDetractionType{},
		&models.UtilSystemIscType{},
		&models.UtilTransferReasonType{},
		&models.UtilTransportModeType{},
		&models.UtilTributeType{},
		&models.UtilUnitMeasureType{},
	)
}
