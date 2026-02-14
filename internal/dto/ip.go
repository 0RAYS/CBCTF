package dto

type SearchIP struct {
	IP string `form:"ip" json:"ip" binding:"required,ip|cidr"`
}
