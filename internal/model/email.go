package model

type Email struct {
	SmtpID  uint   `json:"smtp_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Success bool   `json:"success"`
	BaseModel
}
