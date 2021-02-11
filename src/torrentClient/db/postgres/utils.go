package postgres

func (d *manager) PartsTablePathForFile(fileId string) string  {
	return d.schemaName + "." + fileId
}

func (d *manager) LoadedFilesTablePath() string  {
	return "hypertube.loaded_files"
}

