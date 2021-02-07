package postgres

func (d *manager) LoadedFilesTablePath() string  {
	return d.schemaName + "." + d.loadedFilesTable
}
