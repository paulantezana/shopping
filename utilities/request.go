package utilities

import "github.com/paulantezana/shopping/config"

// Request Type
// 2 = primordial
// 1 = minimal
// 0 = all
type RequestPaginate struct {
	Search      string `json:"search"`
	CurrentPage uint   `json:"current_page"`
	Limit       uint   `json:"limit"`
	Type        uint   `json:"query"`
}

func (r *RequestPaginate) Validate() uint {
	con := config.GetConfig()
	if r.Limit == 0 {
		r.Limit = con.Global.Paginate
	}
	if r.CurrentPage == 0 {
		r.CurrentPage = 1
	}
	offset := r.Limit*r.CurrentPage - r.Limit
	return offset
}

// DeleteRequest use in multiple deletes
type DeleteRequest struct {
	Ids []uint `json:"ids"`
}
