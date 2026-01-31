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

func (e Email) ModelName() string {
	return "Email"
}

func (e Email) GetBaseModel() BaseModel {
	return e.BaseModel
}

func (e Email) UniqueFields() []string {
	return []string{"id"}
}

func (e Email) QueryFields() []string {
	return []string{"id", "from", "to", "subject", "content"}
}
