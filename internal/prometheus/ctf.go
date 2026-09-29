package prometheus

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

// CTFCollector implements prometheus.Collector and reads metrics from DB state.
type CTFCollector struct {
	collectionSuccessDesc   *prometheus.Desc
	contestTeamsDesc        *prometheus.Desc
	contestParticipantsDesc *prometheus.Desc
	victimsActiveDesc       *prometheus.Desc
	victimsDesc             *prometheus.Desc
}

func NewCTFCollector() *CTFCollector {
	return &CTFCollector{
		collectionSuccessDesc: prometheus.NewDesc("cbctf_ctf_collection_success", "Whether all CTF database metrics were collected", nil, nil),
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
	ch <- c.collectionSuccessDesc
	ch <- c.contestTeamsDesc
	ch <- c.contestParticipantsDesc
	ch <- c.victimsActiveDesc
	ch <- c.victimsDesc
}

func (c *CTFCollector) Collect(ch chan<- prometheus.Metric) {
	success := float64(1)
	defer func() { ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, success) }()
	if db.DB == nil {
		success = 0
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	root := db.DB.WithContext(ctx)
	contestRepo := db.InitContestRepo(root)
	contests, _, contestRet := contestRepo.List(-1, -1)
	if !contestRet.OK {
		success = 0
	}
	contestIDL := make([]uint, 0, len(contests))
	for _, contest := range contests {
		contestIDL = append(contestIDL, contest.ID)
	}
	userCountMap, userRet := contestRepo.CountUsersMap(contestIDL...)
	teamCountMap, teamRet := contestRepo.CountTeamsMap(contestIDL...)
	if !userRet.OK || !teamRet.OK {
		success = 0
	}

	for _, contest := range contests {
		if userRet.OK {
			ch <- prometheus.MustNewConstMetric(
				c.contestParticipantsDesc,
				prometheus.GaugeValue,
				float64(userCountMap[contest.ID]),
				strconv.FormatUint(uint64(contest.ID), 10),
			)
		}
		if teamRet.OK {
			ch <- prometheus.MustNewConstMetric(
				c.contestTeamsDesc,
				prometheus.GaugeValue,
				float64(teamCountMap[contest.ID]),
				strconv.FormatUint(uint64(contest.ID), 10),
			)
		}
	}

	type victimStatusCount struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}
	rows := make([]victimStatusCount, 0)
	if err := root.Model(&model.Victim{}).Select("status, count(*) AS count").Group("status").Scan(&rows).Error; err != nil {
		success = 0
		return
	}
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
