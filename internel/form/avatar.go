package form

// DeleteFileForm for delete files
type DeleteFileForm struct {
	FilesID []string `form:"file_id" json:"file_id" banding:"required"`
}
