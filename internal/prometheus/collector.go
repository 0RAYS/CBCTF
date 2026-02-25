package prometheus

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// CTFCollector 实现 prometheus.Collector 接口，通过 DB 实时查询避免状态漂移
type CTFCollector struct {
	contestTeamsDesc        *prometheus.Desc
	contestParticipantsDesc *prometheus.Desc
	victimsActiveDesc       *prometheus.Desc
	generatorPodsDesc       *prometheus.Desc
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
		generatorPodsDesc: prometheus.NewDesc(
			"cbctf_generator_pods_active_total",
			"Number of active generator pods (in-memory)",
			nil, nil,
		),
	}
}

func (c *CTFCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.contestTeamsDesc
	ch <- c.contestParticipantsDesc
	ch <- c.victimsActiveDesc
	ch <- c.generatorPodsDesc
}

func (c *CTFCollector) Collect(ch chan<- prometheus.Metric) {
	if db.DB == nil {
		return
	}
	contestRepo := db.InitContestRepo(db.DB)
	contests, _, _ := contestRepo.List(-1, -1)
	for _, contest := range contests {
		// 参赛选手数量
		ch <- prometheus.MustNewConstMetric(
			c.contestParticipantsDesc,
			prometheus.GaugeValue,
			float64(contestRepo.CountAssociation(contest, "Users")),
			fmt.Sprintf("%d", contest.ID),
		)
		// 参赛队伍数量
		ch <- prometheus.MustNewConstMetric(
			c.contestTeamsDesc,
			prometheus.GaugeValue,
			float64(contestRepo.CountAssociation(contest, "Teams")),
			fmt.Sprintf("%d", contest.ID),
		)
	}

	// 活跃靶机数
	count, _ := db.InitVictimRepo(db.DB).Count(db.CountOptions{Conditions: map[string]any{"deleted_at": nil}})
	ch <- prometheus.MustNewConstMetric(
		c.victimsActiveDesc,
		prometheus.GaugeValue,
		float64(count),
	)

	// 活跃 Generator Pod 数
	k8s.GeneratorMapMutex.RLock()
	genCount := 0
	for _, generators := range k8s.GeneratorMap {
		genCount += len(generators)
	}
	k8s.GeneratorMapMutex.RUnlock()
	ch <- prometheus.MustNewConstMetric(
		c.generatorPodsDesc,
		prometheus.GaugeValue,
		float64(genCount),
	)
}
