package model

type DataResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type LoaderTaskResponse struct {
	Status bool        `json:"status"`
	Data	struct{
		IsLoaded	bool	`json:"isLoaded"`
		Key			string	`json:"key"`
		LoadedPiecesTable string	`json:"loadedPiecesTable"`
		FileName	string		`json:"fileName"`
	}
}
