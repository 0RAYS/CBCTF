package service

import (
	"testing"

	"github.com/gin-gonic/gin/binding"

	"CBCTF/internal/dto"
)

func TestWorkloadSettingNamesFollowKubernetesRules(t *testing.T) {
	for _, tc := range []struct {
		name  string
		form  dto.UpdateSettingForm
		valid bool
	}{
		{name: "hyphenated namespace", form: dto.UpdateSettingForm{K8SNamespace: new("ctf-platform")}, valid: true},
		{name: "uppercase namespace", form: dto.UpdateSettingForm{K8SNamespace: new("CTF")}},
		{name: "empty namespace", form: dto.UpdateSettingForm{K8SNamespace: new("")}},
		{name: "clear priority", form: dto.UpdateSettingForm{K8SPriorityClassName: new("")}, valid: true},
		{name: "priority subdomain", form: dto.UpdateSettingForm{K8SPriorityClassName: new("ctf.high-priority")}, valid: true},
		{name: "invalid priority", form: dto.UpdateSettingForm{K8SPriorityClassName: new("high priority")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid := binding.Validator.ValidateStruct(tc.form) == nil && validateK8sSettings(tc.form).OK
			if valid != tc.valid {
				t.Fatalf("valid=%t, want=%t", valid, tc.valid)
			}
		})
	}
}
