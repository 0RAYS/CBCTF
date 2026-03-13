package model

import "time"

const (
	CheckCheatCronJob           = "CheckCheat"
	UpdateFlagScoreCronJob      = "UpdateFlagScore"
	StopUnCtrlGeneratorCronJob  = "StopUnCtrlGenerator"
	ClearSubmissionMutexCronJob = "ClearSubmissionMutex"
	ClearCheatMutexCronJob      = "ClearCheatMutex"
	ClearJoinTeamMutexCronJob   = "ClearJoinTeamMutes"
	UpdateTeamRankingCronJob    = "UpdateTeamRanking"
	UpdateUserRankingCronJob    = "UpdateUserRanking"
	CollectSystemMetricsCronJob = "CollectSystemMetrics"
	ClearEmptyTeamCronJob       = "ClearEmptyTeam"
	CloseTimeoutVictimsCronJob  = "CloseTimeoutVictims"
	CloseUnCtrlVictimsCronJob   = "CloseUnCtrlVictims"
)

var CronJobs = []CronJob{
	{Name: CollectSystemMetricsCronJob, Description: "收集系统监控指标", Schedule: time.Second},
	{Name: CloseTimeoutVictimsCronJob, Description: "关闭运行超时的靶机实例", Schedule: time.Minute},
	{Name: CloseUnCtrlVictimsCronJob, Description: "清理数据库外仍在运行的失控靶机实例", Schedule: 10 * time.Minute},
	{Name: ClearEmptyTeamCronJob, Description: "清理没有成员的空队伍", Schedule: 5 * time.Minute},
	{Name: UpdateFlagScoreCronJob, Description: "重算比赛题目 Flag 分数和解题人数", Schedule: 5 * time.Minute},
	{Name: UpdateUserRankingCronJob, Description: "全量刷新用户得分和排名", Schedule: 3 * time.Hour},
	{Name: UpdateTeamRankingCronJob, Description: "全量刷新队伍得分和排名", Schedule: 5 * time.Minute},
	{Name: StopUnCtrlGeneratorCronJob, Description: "清理未受数据库管控的附件生成器 Pod", Schedule: 10 * time.Minute},
	{Name: ClearSubmissionMutexCronJob, Description: "清理解题提交锁缓存", Schedule: 10 * time.Minute},
	{Name: CheckCheatCronJob, Description: "扫描并分析比赛作弊事件", Schedule: 10 * time.Minute},
	{Name: ClearCheatMutexCronJob, Description: "清理作弊检测锁缓存", Schedule: 10 * time.Minute},
	{Name: ClearJoinTeamMutexCronJob, Description: "清理队伍加入锁缓存", Schedule: 10 * time.Minute},
}

type CronJob struct {
	Name        string        `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string        `json:"description"`
	Schedule    time.Duration `gorm:"not null" json:"schedule"`
	Last        time.Time     `json:"last"`
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
	return []string{"id", "name", "description", "schedule", "last"}
}
