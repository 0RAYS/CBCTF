package prometheus

import (
	"CBCTF/internal/db"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// CTFCollector implements prometheus.Collector and reads metrics from DB state.
type CTFCollector struct {
	contestTeamsDesc        *prometheus.Desc
	contestParticipantsDesc *prometheus.Desc
	victimsActiveDesc       *prometheus.Desc
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
			"cbctf_victims_active_total",
			"Number of active victim containers (DB-driven)",
			nil, nil,
		),
	}
}

func (c *CTFCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.contestTeamsDesc
	ch <- c.contestParticipantsDesc
	ch <- c.victimsActiveDesc
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
			fmt.Sprintf("%d", contest.ID),
		)
		ch <- prometheus.MustNewConstMetric(
			c.contestTeamsDesc,
			prometheus.GaugeValue,
			float64(teamCountMap[contest.ID]),
			fmt.Sprintf("%d", contest.ID),
		)
	}

	count, _ := db.InitVictimRepo(db.DB).Count(db.CountOptions{Conditions: map[string]any{"deleted_at": nil}})
	ch <- prometheus.MustNewConstMetric(
		c.victimsActiveDesc,
		prometheus.GaugeValue,
		float64(count),
	)
}
