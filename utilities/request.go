package utilities

import "github.com/paulantezana/shopping/config"

type Request struct {
    Search      string `json:"search"`
    CurrentPage uint   `json:"current_page"`
    Limit       uint   `json:"limit"`
    IDs            []uint `json:"i_ds"`
}

func (r *Request) Validate() uint {
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