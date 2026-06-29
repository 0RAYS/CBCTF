package prometheus

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// CTFCollector implements prometheus.Collector and reads metrics from DB state.
type CTFCollector struct {
	contestTeamsDesc        *prometheus.Desc
	contestParticipantsDesc *prometheus.Desc
	victimsActiveDesc       *prometheus.Desc
	victimsDesc             *prometheus.Desc
}

func NewCTFCollector() *CTFCollector {
	return &CTFCollector{
		contestTeamsDesc: prometheus.NewDesc(
			"cbctf_contest_teams_total",
			"Number of teams per contest (DB-driven)",
			[]string{"contest_id"}, nil,
		),
		contestParticipantsDesc: prometheus.NewDesc(
			"cbctf_contest_participants_total",
			"Number of participants per contest (DB-driven)",
			[]string{"contest_id"}, nil,
		),
		victimsActiveDesc: prometheus.NewDesc(
			"cbctf_victims_active",
			"Number of running victim containers (DB-driven)",
			nil, nil,
		),
		victimsDesc: prometheus.NewDesc(
			"cbctf_victims",
			"Number of victim containers by status (DB-driven)",
			[]string{"status"}, nil,
		),
	}
}

func (c *CTFCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.contestTeamsDesc
	ch <- c.contestParticipantsDesc
	ch <- c.victimsActiveDesc
	ch <- c.victimsDesc
}

func (c *CTFCollector) Collect(ch chan<- prometheus.Metric) {
	if db.DB == nil {
		return
	}

	contestRepo := db.InitContestRepo(db.DB)
	contests, _, _ := contestRepo.List(-1, -1)
	contestIDL := make([]uint, 0, len(contests))
	for _, contest := range contests {
		contestIDL = append(contestIDL, contest.ID)
	}
	userCountMap, _ := contestRepo.CountUsersMap(contestIDL...)
	teamCountMap, _ := contestRepo.CountTeamsMap(contestIDL...)

	for _, contest := range contests {
		ch <- prometheus.MustNewConstMetric(
			c.contestParticipantsDesc,
			prometheus.GaugeValue,
			float64(userCountMap[contest.ID]),
			strconv.FormatUint(uint64(contest.ID), 10),
		)
		ch <- prometheus.MustNewConstMetric(
			c.contestTeamsDesc,
			prometheus.GaugeValue,
			float64(teamCountMap[contest.ID]),
			strconv.FormatUint(uint64(contest.ID), 10),
		)
	}

	type victimStatusCount struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}
	rows := make([]victimStatusCount, 0)
	_ = db.DB.Model(&model.Victim{}).Select("status, count(*) AS count").Group("status").Scan(&rows).Error
	running := int64(0)
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(c.victimsDesc, prometheus.GaugeValue, float64(row.Count), row.Status)
		if row.Status == model.RunningVictimStatus {
			running = row.Count
		}
	}
	ch <- prometheus.MustNewConstMetric(
		c.victimsActiveDesc,
		prometheus.GaugeValue,
		float64(running),
	)
}
