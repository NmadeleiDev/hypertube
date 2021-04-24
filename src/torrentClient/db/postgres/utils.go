package postgres

func (d *manager) PartsTablePathForFile(fileId string) string  {
	return d.schemaName + "." + tableNamePrefix + fileId
}

func (d *manager) PartsTableNameForFile(fileId string) string  {
	return tableNamePrefix + fileId
}

func (d *manager) LoadedFilesTablePath() string  {
	return "hypertube.loaded_files"
}

