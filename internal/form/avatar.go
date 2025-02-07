package form

// DeleteAvatarForm for delete files
type DeleteAvatarForm struct {
	Force   bool     `form:"force" json:"force"`
	FilesID []string `form:"file_ids" json:"file_ids"`
}
