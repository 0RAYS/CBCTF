package cron

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

func TestVictimExpiryRequiresRunning(t *testing.T) {
	now := time.Now()
	for _, status := range []string{model.WaitingVictimStatus, model.PendingVictimStatus, model.TerminatingVictimStatus, model.RunningVictimStatus} {
		victim := model.Victim{Status: status, Start: now.Add(-2 * time.Hour), Duration: time.Hour}
		if victimExpired(victim, now) != (status == model.RunningVictimStatus) {
			t.Fatal(status)
		}
	}
	if victimExpired(model.Victim{Status: model.RunningVictimStatus}, now) {
		t.Fatal("zero start must not expire")
	}
}

func TestOrphanScanFailsClosedAndDeduplicates(t *testing.T) {
	var pods []corev1.Pod
	for _, id := range []string{"", "invalid", "-1", "0", "1", "1", "2"} {
		pods = append(pods, corev1.Pod{Labels: map[string]string{"victim_id": id}})
	}
	calls := 0
	ids, ret := orphanVictimIDs(pods, func(id uint) model.RetVal {
		calls++
		if id == 1 {
			return model.RetVal{Msg: i18n.Model.NotFound}
		}
		return model.RetVal{Msg: i18n.DB.Unavailable}
	})
	if ret.OK || ids != nil || calls != 2 {
		t.Fatalf("unsafe scan: %v %+v calls=%d", ids, ret, calls)
	}
	ids, ret = orphanVictimIDs(pods, func(id uint) model.RetVal {
		if id == 1 {
			return model.RetVal{Msg: i18n.Model.NotFound}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK || len(ids) != 1 || ids[0] != "1" {
		t.Fatalf("wrong candidates: %v %+v", ids, ret)
	}
}
