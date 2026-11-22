package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/config"
	p "CBCTF/internal/prometheus"
)

func TestMetricsEndpointExposesRegisteredMetrics(t *testing.T) {
	previous := config.Env
	t.Cleanup(func() { config.Env = previous })
	config.Env = &config.Config{}
	config.Env.Gin.Metrics.Whitelist = []string{"127.0.0.1"}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterMetricsRouter(router)

	// Counter vectors only appear in the exposition once they hold a child.
	p.HttpRequestsTotal.WithLabelValues("GET", "/login", "200", "2xx").Inc()
	p.TaskQueueDepth.WithLabelValues("tasks:victim:start", "pending").Set(3)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	request.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", recorder.Code, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, name := range []string{
		"cbctf_http_requests_total",
		"cbctf_http_in_flight_requests",
		"cbctf_task_queue_depth",
		"cbctf_build_info",
		"cbctf_uptime_seconds",
		"cbctf_dependency_ready",
		"cbctf_ctf_collection_success",
		"cbctf_redis_collection_success",
		"postgres_collection_success",
	} {
		if !strings.Contains(body, name) {
			t.Errorf("metric %s missing from the exposition", name)
		}
	}
}
