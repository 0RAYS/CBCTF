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
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid := binding.Validator.ValidateStruct(tc.form) == nil && validateK8sSettings(tc.form).OK
			if valid != tc.valid {
				t.Fatalf("valid=%t, want=%t", valid, tc.valid)
			}
		})
	}
}
