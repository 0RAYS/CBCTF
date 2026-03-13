package model

type CronJob struct {
	Name        string `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	Schedule    string `gorm:"not null" json:"schedule"`
	Status      string `json:"status"`
	BaseModel
}

func (c CronJob) TableName() string {
	return "cron_jobs"
}

func (c CronJob) ModelName() string {
	return "CronJob"
}

func (c CronJob) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c CronJob) UniqueFields() []string {
	return []string{"id", "name"}
}

func (c CronJob) QueryFields() []string {
	return []string{"id", "name", "description", "schedule", "status"}
}
