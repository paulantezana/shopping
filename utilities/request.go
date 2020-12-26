package utilities

import (
    "github.com/paulantezana/shopping/provider"
    "time"
)

type Request struct {
	Search      string `json:"search"`
	CurrentPage uint   `json:"current_page"`
	PageSize    uint   `json:"page_size"`
	WareHouseId uint   `json:"ware_house_id"`
    StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	IDs         []uint `json:"i_ds"`
}

func (r *Request) Validate() uint {
	con := provider.GetConfig()
	if r.PageSize == 0 {
		r.PageSize = con.Global.PageLimit
	}
	if r.CurrentPage == 0 {
		r.CurrentPage = 1
	}
	offset := r.PageSize*r.CurrentPage - r.PageSize
	return offset
}
