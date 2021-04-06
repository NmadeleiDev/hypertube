package model

type DataResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

