package prometheus

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

type contestSnapshot struct {
	ID           uint
	Status       string
	Teams        int64
	Participants int64
}

type statusCount struct {
	Status string
	Count  int64
}

type challengeCount struct {
	ContestID uint
	Type      string
	Count     int64
}

type ctfSnapshot struct {
	contests        []contestSnapshot
	contestTotals   []statusCount
	challenges      []challengeCount
	victimCounts    []statusCount
	generatorCounts []statusCount
}

// CTFCollector implements prometheus.Collector and reads metrics from DB state.
type CTFCollector struct {
	collectionSuccessDesc   *prometheus.Desc
	contestTeamsDesc        *prometheus.Desc
	contestParticipantsDesc *prometheus.Desc
	contestStatusDesc       *prometheus.Desc
	contestTotalDesc        *prometheus.Desc
	challengeCountDesc      *prometheus.Desc
	victimsActiveDesc       *prometheus.Desc
	victimsDesc             *prometheus.Desc
	generatorStatusDesc     *prometheus.Desc

	cache cacheMemo[ctfSnapshot]
}

func NewCTFCollector() *CTFCollector {
	return &CTFCollector{
		collectionSuccessDesc: prometheus.NewDesc("cbctf_ctf_collection_success", "Whether a CTF database metric section was collected", []string{"section"}, nil),
		contestTeamsDesc: prometheus.NewDesc(
			"cbctf_contest_teams",
			"Number of teams in an unfinished contest (DB-driven)",
			[]string{"contest_id"}, nil,
		),
		contestParticipantsDesc: prometheus.NewDesc(
			"cbctf_contest_participants",
			"Number of participants in an unfinished contest (DB-driven)",
			[]string{"contest_id"}, nil,
		),
		contestStatusDesc: prometheus.NewDesc(
			"cbctf_contest_status",
			"Whether an unfinished contest is in the given state (DB-driven)",
			[]string{"contest_id", "status"}, nil,
		),
		contestTotalDesc: prometheus.NewDesc(
			"cbctf_contests",
			"Number of contests by state (DB-driven)",
			[]string{"status"}, nil,
		),
		challengeCountDesc: prometheus.NewDesc(
			"cbctf_challenges",
			"Number of challenges in a contest by type (DB-driven)",
			[]string{"contest_id", "challenge_type"}, nil,
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
		generatorStatusDesc: prometheus.NewDesc(
			"cbctf_generators",
			"Number of dynamic challenge generators by status (DB-driven)",
			[]string{"status"}, nil,
		),
	}
}

func (c *CTFCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.collectionSuccessDesc
	ch <- c.contestTeamsDesc
	ch <- c.contestParticipantsDesc
	ch <- c.contestStatusDesc
	ch <- c.contestTotalDesc
	ch <- c.challengeCountDesc
	ch <- c.victimsActiveDesc
	ch <- c.victimsDesc
	ch <- c.generatorStatusDesc
}

func (c *CTFCollector) Collect(ch chan<- prometheus.Metric) {
	if db.DB == nil {
		for _, section := range []string{"contests", "challenges", "victims", "generators"} {
			ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, 0, section)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	root := db.DB.WithContext(ctx)

	snapshot, fresh := c.cache.get(func() (ctfSnapshot, bool) {
		return c.collectSnapshot(root)
	})
	for _, section := range []string{"contests", "challenges", "victims", "generators"} {
		success := float64(1)
		if !fresh {
			success = 0
		}
		ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, success, section)
	}

	for _, contest := range snapshot.contests {
		contestID := strconv.FormatUint(uint64(contest.ID), 10)
		ch <- prometheus.MustNewConstMetric(c.contestTeamsDesc, prometheus.GaugeValue, float64(contest.Teams), contestID)
		ch <- prometheus.MustNewConstMetric(c.contestParticipantsDesc, prometheus.GaugeValue, float64(contest.Participants), contestID)
		ch <- prometheus.MustNewConstMetric(c.contestStatusDesc, prometheus.GaugeValue, 1, contestID, contest.Status)
	}
	for _, total := range snapshot.contestTotals {
		ch <- prometheus.MustNewConstMetric(c.contestTotalDesc, prometheus.GaugeValue, float64(total.Count), total.Status)
	}
	for _, challenge := range snapshot.challenges {
		ch <- prometheus.MustNewConstMetric(
			c.challengeCountDesc,
			prometheus.GaugeValue,
			float64(challenge.Count),
			strconv.FormatUint(uint64(challenge.ContestID), 10),
			challenge.Type,
		)
	}
	for _, count := range snapshot.generatorCounts {
		ch <- prometheus.MustNewConstMetric(c.generatorStatusDesc, prometheus.GaugeValue, float64(count.Count), count.Status)
	}

	running := int64(0)
	for _, count := range snapshot.victimCounts {
		ch <- prometheus.MustNewConstMetric(c.victimsDesc, prometheus.GaugeValue, float64(count.Count), count.Status)
		if count.Status == model.RunningVictimStatus {
			running = count.Count
		}
	}
	ch <- prometheus.MustNewConstMetric(c.victimsActiveDesc, prometheus.GaugeValue, float64(running))
}

func (c *CTFCollector) collectSnapshot(root *gorm.DB) (ctfSnapshot, bool) {
	snapshot := ctfSnapshot{}
	var rows []model.Contest
	if err := root.Model(&model.Contest{}).
		Select("id", "start", "duration", "paused", "paused_at", "hidden").
		Find(&rows).Error; err != nil {
		return snapshot, false
	}
	totals := make(map[string]int64, 4)
	live := make([]model.Contest, 0, len(rows))
	for _, contest := range rows {
		status := contestStatus(contest)
		totals[status]++
		if status != contestStatusOver {
			live = append(live, contest)
		}
	}
	for _, status := range []string{contestStatusComing, contestStatusRunning, contestStatusPaused, contestStatusOver} {
		snapshot.contestTotals = append(snapshot.contestTotals, statusCount{Status: status, Count: totals[status]})
	}

	contestIDs := make([]uint, 0, len(live))
	for _, contest := range live {
		contestIDs = append(contestIDs, contest.ID)
	}
	teamCounts, teamOK := db.InitContestRepo(root).CountTeamsMap(contestIDs...)
	userCounts, userOK := db.InitContestRepo(root).CountUsersMap(contestIDs...)
	if !teamOK.OK || !userOK.OK {
		return snapshot, false
	}
	for _, contest := range live {
		snapshot.contests = append(snapshot.contests, contestSnapshot{
			ID:           contest.ID,
			Status:       contestStatus(contest),
			Teams:        teamCounts[contest.ID],
			Participants: userCounts[contest.ID],
		})
	}

	var challengeRows []challengeCount
	if err := root.Model(&model.Challenge{}).
		Select("contest_id, type AS type, count(*) AS count").
		Where("contest_id IS NOT NULL").
		Group("contest_id, type").Scan(&challengeRows).Error; err != nil {
		return snapshot, false
	}
	snapshot.challenges = challengeRows

	var victimRows []statusCount
	if err := root.Model(&model.Victim{}).Select("status, count(*) AS count").Group("status").Scan(&victimRows).Error; err != nil {
		return snapshot, false
	}
	snapshot.victimCounts = victimRows

	var generatorRows []statusCount
	if err := root.Model(&model.Generator{}).Select("status, count(*) AS count").Group("status").Scan(&generatorRows).Error; err != nil {
		return snapshot, false
	}
	snapshot.generatorCounts = generatorRows

	return snapshot, true
}

const (
	contestStatusComing  = "coming"
	contestStatusRunning = "running"
	contestStatusPaused  = "paused"
	contestStatusOver    = "over"
)

// contestStatus mirrors model.Contest.Status with locale-independent values,
// so the label set stays stable regardless of the request language.
func contestStatus(contest model.Contest) string {
	switch {
	case time.Now().After(contest.End()):
		return contestStatusOver
	case time.Now().Before(contest.Start):
		return contestStatusComing
	case contest.Paused:
		return contestStatusPaused
	default:
		return contestStatusRunning
	}
}
