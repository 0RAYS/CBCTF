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

const releaseNamespace = "competition"

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
	args := []string{"template", "ctf", ".", "--namespace", releaseNamespace}
	for _, override := range overrides {
		flag := "--set"
		if strings.HasSuffix(override, "=[]") {
			flag = "--set-json"
		}
		args = append(args, flag, override)
	}
	output, err := exec.Command(helm, args...).CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "missing in charts/ directory") {
			t.Skip("chart dependencies are not vendored; run `helm dependency build chart/`")
		}
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

// valkeyDocuments returns everything the bundled upstream valkey chart rendered.
func valkeyDocuments(documents []map[string]any) []map[string]any {
	var found []map[string]any
	for _, document := range documents {
		if ownedByValkey(document) {
			found = append(found, document)
		}
	}
	return found
}

// ownedByValkey reports whether a document came from the upstream valkey chart
// rather than from this chart's own templates.
func ownedByValkey(document map[string]any) bool {
	metadata, ok := document["metadata"].(map[string]any)
	if !ok {
		return false
	}
	labels, _ := metadata["labels"].(map[string]any)
	return labels["app.kubernetes.io/name"] == "valkey"
}

// valkeyService returns the Service the bundled Valkey is reached through.
func valkeyService(t *testing.T, documents []map[string]any) (string, int) {
	t.Helper()
	for _, document := range valkeyDocuments(documents) {
		if document["kind"] != "Service" {
			continue
		}
		metadata := document["metadata"].(map[string]any)
		if metadata["labels"].(map[string]any)["app.kubernetes.io/component"] != "primary" {
			continue
		}
		ports := document["spec"].(map[string]any)["ports"].([]any)
		return metadata["name"].(string), ports[0].(map[string]any)["port"].(int)
	}
	t.Fatal("bundled Valkey rendered no primary Service")
	return "", 0
}

func TestExternalHostsAreNotNamespaceQualified(t *testing.T) {
	for _, host := range []string{"database.example.com", "192.0.2.10", "database.other.svc.cluster.local"} {
		t.Run(host, func(t *testing.T) {
			documents := render(t, "postgres.enabled=false", "valkey.enabled=false", "postgres.externalHost="+host, "valkey.externalHost="+host, "valkey.password=external-secret")
			config := appConfig(t, documents)
			postgres := config["gorm"].(map[string]any)["postgres"].(map[string]any)
			valkey := config["redis"].(map[string]any)
			if postgres["host"] != host || valkey["host"] != host {
				t.Fatalf("external host changed: postgres=%v valkey=%v", postgres["host"], valkey["host"])
			}
			if bundled := valkeyDocuments(documents); len(bundled) != 0 {
				t.Fatalf("external mode rendered the bundled Valkey: %d resources", len(bundled))
			}
			for _, document := range documents {
				if document["kind"] == "StatefulSet" {
					t.Fatal("external mode rendered a bundled database")
				}
			}
		})
	}
}

// TestBundledValkeyAddressMatchesRenderedService pins cbctf.valkey.fullname, the
// only place this chart reimplements the upstream valkey chart's naming rule.
// The application configuration has to reach the Service the subchart renders,
// and nothing else in this chart could notice the two drifting apart.
func TestBundledValkeyAddressMatchesRenderedService(t *testing.T) {
	documents := render(t, "valkey.auth.aclUsers.default.password=bundled-secret")
	config := appConfig(t, documents)
	connection := config["redis"].(map[string]any)
	if connection["pwd"] != "bundled-secret" {
		t.Fatalf("platform password does not follow the bundled Valkey ACL entry: %v", connection["pwd"])
	}
	name, port := valkeyService(t, documents)
	if connection["host"] != name+"."+releaseNamespace {
		t.Fatalf("platform addresses %v but the bundled Valkey Service is %s", connection["host"], name)
	}
	if connection["port"] != port {
		t.Fatalf("platform dials port %v but the bundled Valkey Service exposes %d", connection["port"], port)
	}
}

func TestValkeyEphemeralVolumeAndPullSecrets(t *testing.T) {
	documents := render(t, "valkey.dataStorage.enabled=false", "valkey.imagePullSecrets[0]=registry-key")
	for _, document := range documents {
		if document["kind"] == "PersistentVolumeClaim" && ownedByValkey(document) {
			t.Fatal("ephemeral Valkey rendered a PVC")
		}
	}
	var pod map[string]any
	for _, document := range valkeyDocuments(documents) {
		if document["kind"] == "Deployment" {
			pod = document["spec"].(map[string]any)["template"].(map[string]any)["spec"].(map[string]any)
		}
	}
	if pod == nil {
		t.Fatal("bundled Valkey rendered no Deployment")
	}
	if secret := pod["imagePullSecrets"].([]any)[0].(map[string]any)["name"]; secret != "registry-key" {
		t.Fatal("Valkey image pull secret is missing")
	}
	found := false
	for _, raw := range pod["volumes"].([]any) {
		volume := raw.(map[string]any)
		if volume["name"] == "valkey-data" {
			_, found = volume["emptyDir"]
		}
	}
	if !found {
		t.Fatal("Valkey data mount has no emptyDir volume")
	}
}

func TestValkeyPersistenceIsRetainedOnUninstall(t *testing.T) {
	documents := render(t)
	var volume bool
	for _, document := range documents {
		metadata := document["metadata"].(map[string]any)
		labels, _ := metadata["labels"].(map[string]any)
		if document["kind"] != "PersistentVolumeClaim" || labels["app.kubernetes.io/name"] != "valkey" {
			continue
		}
		annotations, _ := metadata["annotations"].(map[string]any)
		if annotations["helm.sh/resource-policy"] != "keep" {
			t.Fatal("Valkey volume is deleted on helm uninstall")
		}
		volume = true
	}
	if !volume {
		t.Fatal("bundled Valkey rendered no volume")
	}
}

func TestRuntimeAndApplicationShareExistingClaim(t *testing.T) {
	documents := render(t, "persistence.existingClaim=shared-data")
	config := appConfig(t, documents)
	if config["k8s"].(map[string]any)["shared_volume_claim"] != "shared-data" {
		t.Fatal("runtime PVC does not match existingClaim")
	}
	// The only PVCs rendered as standalone objects are the shared data volume
	// and the one the valkey subchart asks for; PostgreSQL uses a claim template.
	for _, document := range documents {
		if document["kind"] == "PersistentVolumeClaim" && !ownedByValkey(document) {
			t.Fatal("existingClaim must not create another PVC")
		}
		if document["kind"] != "Deployment" || ownedByValkey(document) {
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

func TestPostgresEphemeralVolumeAndPullSecrets(t *testing.T) {
	documents := render(t, "postgres.persistence.enabled=false", "imagePullSecrets[0].name=registry-key")
	for _, document := range documents {
		if document["kind"] != "StatefulSet" {
			continue
		}
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
}
