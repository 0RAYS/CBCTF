package resp

import (
	"encoding/json"
	"strings"
	"testing"

	"CBCTF/internal/model"
)

func TestGeneratorResponseIncludesRuntimeImageWithoutWorkerCredentials(t *testing.T) {
	generator := model.Generator{Image: "registry.example/generator:v2", WorkerToken: "worker-private-token", Status: model.RunningGeneratorStatus}
	response := GetGeneratorResp(generator)
	if response["image"] != generator.Image {
		t.Fatal("runtime image omitted")
	}
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), generator.WorkerToken) {
		t.Fatal("worker credentials exposed")
	}
}
