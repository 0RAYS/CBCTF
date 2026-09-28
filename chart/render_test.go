package chart_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

func render(t *testing.T, overrides ...string) []map[string]any {
	t.Helper()
	helm := os.Getenv("HELM_BIN")
	if helm == "" {
		var err error
		helm, err = exec.LookPath("helm")
		if err != nil {
			t.Skip("helm is required; set HELM_BIN or add helm to PATH")
		}
	}
	args := []string{"template", "ctf", ".", "--namespace", "competition"}
	for _, override := range overrides {
		flag := "--set"
		if strings.HasSuffix(override, "=[]") {
			flag = "--set-json"
		}
		args = append(args, flag, override)
	}
	output, err := exec.Command(helm, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("helm template: %v\n%s", err, output)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(output))
	var documents []map[string]any
	for {
		var document map[string]any
		if err := decoder.Decode(&document); err == io.EOF {
			return documents
		} else if err != nil {
			t.Fatal(err)
		}
		if document != nil {
			documents = append(documents, document)
		}
	}
}

func appConfig(t *testing.T, documents []map[string]any) map[string]any {
	t.Helper()
	for _, document := range documents {
		if document["kind"] != "ConfigMap" {
			continue
		}
		if raw, ok := document["data"].(map[string]any)["config.yaml"].(string); ok {
			var config map[string]any
			if err := yaml.Unmarshal([]byte(raw), &config); err != nil {
				t.Fatal(err)
			}
			return config
		}
	}
	t.Fatal("missing application configuration")
	return nil
}

func TestExternalHostsAreNotNamespaceQualified(t *testing.T) {
	for _, host := range []string{"database.example.com", "192.0.2.10", "database.other.svc.cluster.local"} {
		t.Run(host, func(t *testing.T) {
			documents := render(t, "postgres.enabled=false", "redis.enabled=false", "postgres.externalHost="+host, "redis.externalHost="+host)
			config := appConfig(t, documents)
			postgres := config["gorm"].(map[string]any)["postgres"].(map[string]any)
			redis := config["redis"].(map[string]any)
			if postgres["host"] != host || redis["host"] != host {
				t.Fatalf("external host changed: postgres=%v redis=%v", postgres["host"], redis["host"])
			}
			for _, document := range documents {
				if document["kind"] == "StatefulSet" {
					t.Fatal("external mode rendered a bundled database")
				}
			}
		})
	}
}

func TestRuntimeAndApplicationShareExistingClaim(t *testing.T) {
	documents := render(t, "persistence.existingClaim=shared-data")
	config := appConfig(t, documents)
	if config["k8s"].(map[string]any)["shared_volume_claim"] != "shared-data" {
		t.Fatal("runtime PVC does not match existingClaim")
	}
	for _, document := range documents {
		if document["kind"] == "PersistentVolumeClaim" {
			t.Fatal("existingClaim must not create another PVC")
		}
		if document["kind"] != "Deployment" {
			continue
		}
		spec := document["spec"].(map[string]any)["template"].(map[string]any)["spec"].(map[string]any)
		for _, raw := range spec["volumes"].([]any) {
			volume := raw.(map[string]any)
			if volume["name"] == "data" && volume["persistentVolumeClaim"].(map[string]any)["claimName"] != "shared-data" {
				t.Fatal("application PVC does not match existingClaim")
			}
		}
	}
}

func TestEmptyPolicyListsRemainLists(t *testing.T) {
	config := appConfig(t, render(t, "cbctf.gin.origins=[]", "cbctf.gin.proxies=[]"))
	gin := config["gin"].(map[string]any)
	for _, key := range []string{"origins", "proxies"} {
		list, ok := gin[key].([]any)
		if !ok || len(list) != 0 {
			t.Fatalf("%s must be an explicit empty list, got %#v", key, gin[key])
		}
	}
}

func TestDatabaseEphemeralVolumesAndPullSecrets(t *testing.T) {
	documents := render(t, "postgres.persistence.enabled=false", "redis.persistence.enabled=false", "imagePullSecrets[0].name=registry-key")
	count := 0
	for _, document := range documents {
		if document["kind"] != "StatefulSet" {
			continue
		}
		count++
		spec := document["spec"].(map[string]any)
		if _, ok := spec["volumeClaimTemplates"]; ok {
			t.Fatal("ephemeral database rendered a PVC template")
		}
		pod := spec["template"].(map[string]any)["spec"].(map[string]any)
		secrets := pod["imagePullSecrets"].([]any)
		if secrets[0].(map[string]any)["name"] != "registry-key" {
			t.Fatal("database image pull secret is missing")
		}
		found := false
		for _, raw := range pod["volumes"].([]any) {
			volume := raw.(map[string]any)
			if volume["name"] == "data" {
				_, found = volume["emptyDir"]
			}
		}
		if !found {
			t.Fatal("database data mount has no emptyDir volume")
		}
	}
	if count != 2 {
		t.Fatalf("expected both bundled databases, got %d", count)
	}
}
