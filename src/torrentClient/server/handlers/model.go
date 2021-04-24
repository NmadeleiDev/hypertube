package handlers

type DataResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

