package dto

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLeaderboardFormBinding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/x?limit=20&offset=40&show_hidden=true", nil)
	var form LeaderboardForm
	ret := Bind(ctx, &form)
	if !ret.OK {
		t.Fatalf("bind failed: %+v", ret)
	}
	if form.Limit != 20 || form.Offset != 40 || !form.ShowHidden {
		t.Fatalf("unexpected form: %+v", form)
	}
}

func TestLeaderboardFormDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/x", nil)
	var form LeaderboardForm
	ret := Bind(ctx, &form)
	if !ret.OK {
		t.Fatalf("bind failed: %+v", ret)
	}
	if form.Limit != 10 || form.Offset != 0 || form.ShowHidden {
		t.Fatalf("unexpected form: %+v", form)
	}
}
