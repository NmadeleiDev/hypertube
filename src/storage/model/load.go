package model

type LoadInfo struct {
	VideoFile	FileInfo
	SrtFile		FileInfo
	IsLoaded	bool
	InProgress	bool
}

type FileInfo struct {
	Name		string
	Length		int64
}
