package router

import (
	"CBCTF/internal/model"
	"testing"
)

func TestGeneratorStatusRequiresSamePermissionAsLogs(t *testing.T) {
	for _, base := range []string{"GET /admin/generators/:generatorID", "GET /admin/contests/:contestID/generators/:generatorID"} {
		logs, ok := model.RoutePermissions[base+"/logs"]
		status, exists := model.RoutePermissions[base+"/status"]
		if !ok || !exists || logs != status || status == "" {
			t.Fatalf("missing/mismatched permission for %s", base)
		}
	}
}
