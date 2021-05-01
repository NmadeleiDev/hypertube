package model

type DataResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type LoaderTaskResponse struct {
	Status bool        `json:"status"`
	Data	struct{
		IsLoaded	bool	`json:"isLoaded"`
		IsLoading	bool	`json:"isLoading"`
		FileName	string		`json:"fileName"`
	}
}
