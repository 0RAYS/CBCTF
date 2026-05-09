package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
)

const loadTrafficTaskType = "tasks:traffic:load"

type LoadTrafficPayload struct {
	Victim model.Victim
}

func EnqueueLoadTrafficTask(victim model.Victim) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(LoadTrafficPayload{Victim: victim})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(loadTrafficTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(loadTrafficTaskType), asynq.MaxRetry(3), asynq.Timeout(5*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(loadTrafficTaskType)
	}
	return info, err
}

func HandleLoadTrafficTask(ctx context.Context, t *asynq.Task) error {
	var payload LoadTrafficPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	victim := payload.Victim
	log.Logger.Infof("Loading victim traffic: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
	ret := LoadTraffic(ctx, db.DB, victim)
	if !ret.OK {
		return fmt.Errorf("load traffic failed: %s", ret.Msg)
	}
	return nil
}

// LoadTraffic 简单记录涉及到的 IP 地址
func LoadTraffic(ctx context.Context, root *gorm.DB, victim model.Victim) model.RetVal {
	trafficRepo := db.InitTrafficRepo(root)
	optionsL := make(map[string]db.CreateTrafficOptions)
	count, _ := trafficRepo.Count(db.CountOptions{Conditions: map[string]any{"victim_id": victim.ID}})
	if count > 0 {
		log.Logger.Debugf("Traffic already loaded: victim_id=%d records=%d", victim.ID, count)
		return model.SuccessRetVal()
	}
	if err := ctx.Err(); err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	archive := func(victim model.Victim) (db.CreateFileOptions, bool) {
		start := time.Now()
		log.Logger.Debugf("Enrich pcap with process info: victim_id=%d path=%s", victim.ID, victim.TrafficBasePath())
		if errs := utils.EnrichPcapDirWithContext(ctx, victim.TrafficBasePath()); len(errs) > 0 {
			for _, err := range errs {
				log.Logger.Warningf("Enrich pcap error: %v", err)
			}
		}
		if err := ctx.Err(); err != nil {
			log.Logger.Warningf("Traffic archive cancelled after enrichment: victim_id=%d error=%s", victim.ID, err)
			return db.CreateFileOptions{}, false
		}
		log.Logger.Debugf("Archiving victim traffic pcaps: victim_id=%d path=%s", victim.ID, victim.TrafficBasePath())
		if err := utils.ZipWithContext(ctx, victim.TrafficBasePath(), victim.TrafficZipPath()); err != nil {
			log.Logger.Warningf("Failed to archive victim traffic pcaps: victim_id=%d error=%s", victim.ID, err)
			return db.CreateFileOptions{}, false
		}
		if err := ctx.Err(); err != nil {
			log.Logger.Warningf("Traffic archive cancelled after zip: victim_id=%d error=%s", victim.ID, err)
			return db.CreateFileOptions{}, false
		}
		size, hash, err := utils.GetFileInfoByPathWithContext(ctx, victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to get traffic archive info: victim_id=%d path=%s error=%s", victim.ID, victim.TrafficZipPath(), err)
			return db.CreateFileOptions{}, false
		}
		log.Logger.Debugf("Archived victim traffic pcaps: victim_id=%d size=%d duration=%s", victim.ID, size, time.Since(start))
		return db.CreateFileOptions{
			RandID:   utils.UUID(),
			Filename: "traffics.zip",
			Size:     size,
			Path:     model.FilePath(victim.TrafficZipPath()),
			Model:    model.Name(victim),
			ModelID:  victim.ID,
			Suffix:   ".zip",
			Hash:     hash,
			Type:     model.TrafficFileType,
		}, true
	}
	fileOptions, hasArchive := archive(victim)
	if err := ctx.Err(); err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	start := time.Now()
	connections, err := utils.ReadPcapDirWithContext(ctx, victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read victim pcaps: victim_id=%d path=%s error=%s", victim.ID, victim.TrafficBasePath(), err)
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	if err := ctx.Err(); err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	log.Logger.Debugf("Read victim pcaps: victim_id=%d packets=%d duration=%s", victim.ID, len(connections), time.Since(start))
	for _, conn := range connections {
		connID := fmt.Sprintf("%s-%s-%s-%s", conn.SrcIP, conn.DstIP, conn.Type, conn.Subtype)
		if options, ok := optionsL[connID]; ok {
			options.Count += 1
			options.Size += conn.Size
			optionsL[connID] = options
		} else {
			optionsL[connID] = db.CreateTrafficOptions{
				VictimID: victim.ID,
				SrcIP:    conn.SrcIP,
				DstIP:    conn.DstIP,
				Type:     conn.Type,
				Subtype:  conn.Subtype,
				Size:     conn.Size,
				Count:    1,
			}
		}
	}
	ret := db.WithTransactionDB(root, func(tx *db.Tx) model.RetVal {
		if hasArchive {
			if _, ret := db.InitFileRepo(tx).Create(fileOptions); !ret.OK {
				return model.RetVal{Msg: fmt.Sprintf("create traffic archive record failed: %s", ret.Msg)}
			}
		}
		trafficRepo := db.InitTrafficRepo(tx)
		for _, options := range optionsL {
			if err := ctx.Err(); err != nil {
				return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
			}
			_, ret := trafficRepo.Create(options)
			if !ret.OK {
				return ret
			}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}
	log.Logger.Infof("Traffic loaded: victim_id=%d packets=%d records=%d", victim.ID, len(connections), len(optionsL))
	return model.SuccessRetVal()
}
