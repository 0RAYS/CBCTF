package router

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"CBCTF/internal/config"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
)

func TestSystemConfigExposesEditableSettings(t *testing.T) {
	previous, bundle := config.Env, i18n.Bundle
	t.Cleanup(func() { config.Env, i18n.Bundle = previous, bundle })
	config.Env = &config.Config{Path: t.TempDir()}
	config.Env.K8S.CaptureEnabled = false
	config.Env.K8S.PriorityClassName = ""
	config.Env.K8S.WorkerImage = "registry.example/worker:v3"
	config.Env.K8S.GeneratorPoolSize = 0
	i18n.Init()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/admin/system/config", nil)
	SystemConfig(ctx)
	var response struct {
		Code int                        `json:"code"`
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Code != 200 {
		t.Fatal(recorder.Body.String())
	}
	formType := reflect.TypeOf(dto.UpdateSettingForm{})
	for i := 0; i < formType.NumField(); i++ {
		key := formType.Field(i).Tag.Get("json")
		if key == "gin_jwt_secret" {
			continue
		}
		if _, exists := response.Data[key]; !exists {
			t.Errorf("editable setting %s is missing from GET", key)
		}
	}
	for key, want := range map[string]string{
		"k8s_capture_enabled":     "false",
		"k8s_priority_class_name": `""`,
		"k8s_worker_image":        `"registry.example/worker:v3"`,
		"k8s_generator_pool_size": "0",
	} {
		if string(response.Data[key]) != want {
			t.Errorf("%s = %s, want %s", key, response.Data[key], want)
		}
	}
}

func TestWorkloadSettingsPreserveExplicitZeroValues(t *testing.T) {
	var form dto.UpdateSettingForm
	if err := json.Unmarshal([]byte(`{"k8s_capture_enabled":false,"k8s_priority_class_name":"","k8s_worker_image":"worker:v1","k8s_generator_pool_size":0}`), &form); err != nil {
		t.Fatal(err)
	}
	if err := binding.Validator.ValidateStruct(form); err != nil {
		t.Fatal(err)
	}
	if form.K8SCaptureEnabled == nil || *form.K8SCaptureEnabled || form.K8SGeneratorPoolSize == nil || *form.K8SGeneratorPoolSize != 0 || form.K8SPriorityClassName == nil || *form.K8SPriorityClassName != "" {
		t.Fatal("explicit clear/disable was treated as missing")
	}
	form.K8SGeneratorPoolSize = new(-1)
	if err := binding.Validator.ValidateStruct(form); err == nil {
		t.Fatal("negative pool size accepted")
	}
	form.K8SGeneratorPoolSize = new(0)
	form.K8SWorkerImage = new("")
	if err := binding.Validator.ValidateStruct(form); err == nil {
		t.Fatal("empty worker image accepted")
	}
}
