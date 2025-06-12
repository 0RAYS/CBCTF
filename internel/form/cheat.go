package form

type GetCheatsForm struct {
	Offset int    `form:"offset" json:"offset"`
	Limit  int    `form:"limit" json:"limit"`
	Type   string `form:"type" json:"type"`
}
