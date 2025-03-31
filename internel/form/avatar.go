package form

// DeleteFileForm for delete files
type DeleteFileForm struct {
	FilesID []string `form:"file_ids" json:"file_ids"`
}
