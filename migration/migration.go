package migration

import (
	"github.com/paulantezana/shopping/config"
	"github.com/paulantezana/shopping/models"
    "github.com/paulantezana/shopping/models/productmodel"
)

// migration function
func Migrate() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
        productmodel.Category{},
        productmodel.Product{},
        productmodel.ProductCategory{},
        productmodel.Brand{},
        productmodel.UnitMeasure{},
        productmodel.Commentary{},
        productmodel.ProductRelationship{},
        productmodel.Image{},
        productmodel.Presentation{},
        productmodel.Variant{},
        productmodel.Alternative{},

		// General
		models.Country{},
		models.Level1{},
		models.Level2{},
		models.Level3{},
		models.LabelLocation{},
		models.GeographicLocation{},
		models.GeneralData{},
		models.Representative{},
		models.GeneralSetting{},
	)
	db.Model(&productmodel.ProductCategory{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.ProductCategory{}).AddForeignKey("category_id", "categories(id)", "RESTRICT", "RESTRICT")

	db.Model(&productmodel.Product{}).AddForeignKey("unit_measure_id", "unit_measures(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.Product{}).AddForeignKey("brand_id", "brands(id)", "RESTRICT", "RESTRICT")

	db.Model(&productmodel.Commentary{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.ProductRelationship{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.Image{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.Presentation{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.Variant{}).AddForeignKey("presentation_id", "presentations(id)", "RESTRICT", "RESTRICT")
	db.Model(&productmodel.Alternative{}).AddForeignKey("variant_id", "variants(id)", "RESTRICT", "RESTRICT")

	// Generals
	db.Model(&models.Level1{}).AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Level2{}).AddForeignKey("level1_id", "level1(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Level3{}).AddForeignKey("level2_id", "level2(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.LabelLocation{}).AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.GeographicLocation{}).AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.GeographicLocation{}).AddForeignKey("level1_id", "level1(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.GeneralData{}).AddForeignKey("geographic_location_id", "geographic_locations(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Representative{}).AddForeignKey("general_data_id", "general_data(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.GeneralSetting{}).AddForeignKey("general_data_id", "general_data(id)", "RESTRICT", "RESTRICT")

}
