package utilities

// SuccessMessage
type SuccessMessage struct {
	Message string `json:"message"`
}

// ErrorMessage
type ErrorMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Field   string `json:"field"`
}

// Response data
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ResponseAdvance data advance
type ResponseAdvance struct {
	Success  bool             `json:"success"`
	Messages []SuccessMessage `json:"messages"`
	Errors   []ErrorMessage   `json:"errors"`
}

// ResponsePaginate data pagination
type ResponsePaginate struct {
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Total    uint        `json:"total"`
	Current  uint        `json:"current"`
	PageSize uint        `json:"page_size"`
}
