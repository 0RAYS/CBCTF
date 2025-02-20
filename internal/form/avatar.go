package form

// DeleteFileForm for delete files
type DeleteFileForm struct {
	Force   bool     `form:"force" json:"force"`
	FilesID []string `form:"file_ids" json:"file_ids"`
}
