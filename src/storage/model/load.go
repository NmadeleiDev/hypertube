package model

type LoadInfo struct {
	VideoFile	FileInfo
	IsLoaded	bool
	InProgress	bool
}

type FileInfo struct {
	Name		string
	Length		int64
}
