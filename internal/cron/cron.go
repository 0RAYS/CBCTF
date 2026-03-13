package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type taskDefinition struct {
	name string
	run  func() model.RetVal
}

var (
	c           *cron.Cron
	taskEntries sync.Map
	taskMap     = map[string]taskDefinition{
		model.CollectSystemMetricsCronJob: {name: model.CollectSystemMetricsCronJob, run: collectSystemMetricsTask},
		model.CloseTimeoutVictimsCronJob:  {name: model.CloseTimeoutVictimsCronJob, run: closeTimeoutVictimsTask},
		model.CloseUnCtrlVictimsCronJob:   {name: model.CloseUnCtrlVictimsCronJob, run: closeUnCtrlVictimsTask},
		model.ClearEmptyTeamCronJob:       {name: model.ClearEmptyTeamCronJob, run: clearEmptyTeamTask},
		model.UpdateFlagScoreCronJob:      {name: model.UpdateFlagScoreCronJob, run: updateFlagScoreTask},
		model.UpdateUserRankingCronJob:    {name: model.UpdateUserRankingCronJob, run: updateUserRankingTask},
		model.UpdateTeamRankingCronJob:    {name: model.UpdateTeamRankingCronJob, run: updateTeamRankingTask},
		model.StopUnCtrlGeneratorCronJob:  {name: model.StopUnCtrlGeneratorCronJob, run: stopUnCtrlGeneratorTask},
		model.ClearSubmissionMutexCronJob: {name: model.ClearSubmissionMutexCronJob, run: clearSubmissionMutexTask},
		model.CheckCheatCronJob:           {name: model.CheckCheatCronJob, run: checkCheatTask},
		model.ClearCheatMutexCronJob:      {name: model.ClearCheatMutexCronJob, run: clearCheatMutexTask},
		model.ClearJoinTeamMutexCronJob:   {name: model.ClearJoinTeamMutexCronJob, run: clearJoinTeamMutexTask},
	}
)

func exec(name string, task func() model.RetVal) func() {
	return func() {
		start := time.Now()
		ret := task()
		now := time.Now()
		if updateRet := db.InitCronJobRepo(db.DB).UpdateByName(name, db.UpdateCronJobOptions{Last: &now}); !updateRet.OK {
			log.Logger.Warningf("Failed to update cron last runtime %s: %s", name, updateRet.Msg)
		}
		duration := time.Since(start).Seconds()
		prometheus.RecordCronJob(name, duration, ret.OK)
		if !ret.OK {
			log.Logger.Warningf("%s failed: %s, processing time: %s", name, ret.Msg, time.Duration(duration*float64(time.Second)))
		} else if duration > time.Second.Seconds() {
			log.Logger.Warningf("%s processing time: %s", name, time.Duration(duration*float64(time.Second)))
		} else {
			log.Logger.Debugf("%s processing time: %s", name, time.Duration(duration*float64(time.Second)))
		}
	}
}

func Init() {
	c = cron.New(cron.WithSeconds())
}

func Start() {
	log.Logger.Info("Cron started")
	if ret := reloadAll(); !ret.OK {
		log.Logger.Warningf("Failed to load cron jobs: %s %v", ret.Msg, ret.Attr)
	}
	c.Start()
}

func Stop() {
	if c != nil {
		c.Stop()
	}
}

func ReloadCronJob(name string) model.RetVal {
	def, ok := taskMap[name]
	if !ok {
		return model.RetVal{Msg: "Cron job not registered", Attr: map[string]any{"Name": name}}
	}
	cronJob, ret := db.InitCronJobRepo(db.DB).GetByUniqueKey("name", name)
	if !ret.OK {
		return ret
	}
	return registerCronJob(cronJob, def)
}

func reloadAll() model.RetVal {
	cronJobs, _, ret := db.InitCronJobRepo(db.DB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	for _, cronJob := range cronJobs {
		def, ok := taskMap[cronJob.Name]
		if !ok {
			log.Logger.Warningf("Skip unknown cron job: %s", cronJob.Name)
			continue
		}
		if ret = registerCronJob(cronJob, def); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}

func registerCronJob(cronJob model.CronJob, def taskDefinition) model.RetVal {
	if value, ok := taskEntries.Load(cronJob.Name); ok {
		c.Remove(value.(cron.EntryID))
		taskEntries.Delete(cronJob.Name)
	}
	spec := "@every " + cronJob.Schedule.String()
	entryID, err := c.AddFunc(spec, exec(def.name, def.run))
	if err != nil {
		return model.RetVal{Msg: "Invalid cron schedule", Attr: map[string]any{"Name": cronJob.Name, "Schedule": spec, "Error": err.Error()}}
	}
	taskEntries.Store(cronJob.Name, entryID)
	log.Logger.Infof("Cron job loaded: %s (%s)", cronJob.Name, spec)
	return model.SuccessRetVal()
}

func mustSchedule(name, schedule string, task func() model.RetVal) {
	entryID, err := c.AddFunc(schedule, exec(name, task))
	if err != nil {
		log.Logger.Fatalf("Failed to schedule cron job %s: %s", name, err)
	}
	taskEntries.Store(name, entryID)
}

func RegisterCronJobForTest(name, schedule string, task func() model.RetVal) error {
	if c == nil {
		return fmt.Errorf("cron is not initialized")
	}
	if value, ok := taskEntries.Load(name); ok {
		c.Remove(value.(cron.EntryID))
		taskEntries.Delete(name)
	}
	entryID, err := c.AddFunc(schedule, exec(name, task))
	if err != nil {
		return err
	}
	taskEntries.Store(name, entryID)
	return nil
}
