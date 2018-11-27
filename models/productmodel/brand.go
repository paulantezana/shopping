package productmodel

type Brand struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	WebSite      string `json:"web_site"`
	LogoLocation string `json:"logo_location"`
	CaptionImage string `json:"caption_image"`
	State        bool   `json:"state" gorm:"default:'true'"`
}
