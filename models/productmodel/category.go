package productmodel

type Category struct {
	ID                   uint   `json:"id" gorm:"primary_key"`
	Name                 string `json:"name"`
	ParentCategoryID     uint   `json:"parent_category_id"`
	DisplayOrderProducts uint   `json:"display_order_products"`
	ShowProductsIn       uint   `json:"show_products_in"`
	NumberColumns        uint   `json:"number_columns"`
	TitleCategorySeo     string `json:"title_category_seo"`
	UrlCategorySeo       string `json:"url_category_seo"`
	MetaTagsSeo          string `json:"meta_tags_seo"`
	HeaderPage           string `json:"header_page"`
	FootPage             string `json:"foot_page"`
	Position             uint   `json:"position"`
	LocationLogo         uint   `json:"location_logo"`
	ShowWeb              bool   `json:"show_web"`
}
