package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"slices"
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
	return enqueueTask(loadTrafficTaskType, task, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute))
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

// LoadTraffic enrich pcap、打包归档，并提取流量涉及的所有 IP 写入 traffics 表。
func LoadTraffic(ctx context.Context, root *gorm.DB, victim model.Victim) model.RetVal {
	trafficRepo := db.InitTrafficRepo(root)

	count, _ := trafficRepo.Count(db.CountOptions{Conditions: map[string]any{"victim_id": victim.ID}})
	if count > 0 {
		log.Logger.Debugf("Traffic already loaded: victim_id=%d", victim.ID)
		return model.SuccessRetVal()
	}

	if err := ctx.Err(); err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	archive := func(victim model.Victim) (model.File, bool) {
		start := time.Now()
		log.Logger.Debugf("Enrich pcap with process info: victim_id=%d path=%s", victim.ID, victim.TrafficBasePath())
		if errs := utils.EnrichPcapDirWithContext(ctx, victim.TrafficBasePath()); len(errs) > 0 {
			for _, err := range errs {
				log.Logger.Warningf("Enrich pcap error: %v", err)
			}
		}
		if err := ctx.Err(); err != nil {
			log.Logger.Warningf("Traffic archive cancelled after enrichment: victim_id=%d error=%s", victim.ID, err)
			return model.File{}, false
		}
		log.Logger.Debugf("Archiving victim traffic pcaps: victim_id=%d path=%s", victim.ID, victim.TrafficBasePath())
		if err := utils.ZipWithContext(ctx, victim.TrafficBasePath(), victim.TrafficZipPath()); err != nil {
			log.Logger.Warningf("Failed to archive victim traffic pcaps: victim_id=%d error=%s", victim.ID, err)
			return model.File{}, false
		}
		if err := ctx.Err(); err != nil {
			log.Logger.Warningf("Traffic archive cancelled after zip: victim_id=%d error=%s", victim.ID, err)
			return model.File{}, false
		}
		size, hash, err := utils.GetFileInfoByPath(ctx, victim.TrafficZipPath())
		if err != nil {
			log.Logger.Warningf("Failed to get traffic archive info: victim_id=%d path=%s error=%s", victim.ID, victim.TrafficZipPath(), err)
			return model.File{}, false
		}
		log.Logger.Debugf("Archived victim traffic pcaps: victim_id=%d size=%d duration=%s", victim.ID, size, time.Since(start))
		return model.File{
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
	result, err := utils.ReadPcapDirWithContext(ctx, victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read victim pcaps: victim_id=%d path=%s error=%s", victim.ID, victim.TrafficBasePath(), err)
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	if err = ctx.Err(); err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}

	// 收集所有涉及的 IP：普通 pod src/dst + frpc proxy protocol 真实 IP
	ipSet := make(map[string]struct{})
	for _, conn := range result.Connections {
		if conn.SrcIP != "" {
			ipSet[conn.SrcIP] = struct{}{}
		}
		if conn.DstIP != "" {
			ipSet[conn.DstIP] = struct{}{}
		}
	}
	for _, ip := range result.FrpcIPs {
		ipSet[ip] = struct{}{}
	}

	ips := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		ips = append(ips, ip)
	}
	slices.Sort(ips)

	log.Logger.Debugf("Collected IPs from pcaps: victim_id=%d connections=%d frpc_ips=%d unique_ips=%d duration=%s",
		victim.ID, len(result.Connections), len(result.FrpcIPs), len(ips), time.Since(start))

	ret := db.WithTransactionDB(root, func(tx *db.Tx) model.RetVal {
		if hasArchive {
			if _, ret := db.InitFileRepo(tx).Create(fileOptions); !ret.OK {
				return model.RetVal{Msg: fmt.Sprintf("create traffic archive record failed: %s", ret.Msg)}
			}
		}
		return db.InitTrafficRepo(tx).UpsertIPs(victim.ID, ips)
	})
	if !ret.OK {
		return ret
	}
	log.Logger.Infof("Traffic loaded: victim_id=%d unique_ips=%d", victim.ID, len(ips))
	return model.SuccessRetVal()
}
