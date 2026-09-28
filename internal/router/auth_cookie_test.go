package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
)

func TestAuthCookieOriginPolicy(t *testing.T) {
	i18n.Init()
	previous := config.Env
	t.Cleanup(func() { config.Env = previous })
	for _, test := range []struct {
		name, host, origin string
		secure             bool
		sameSite           http.SameSite
	}{
		{"chart HTTP", "http://cbctf.local", "http://cbctf.local", false, http.SameSiteLaxMode},
		{"same-origin HTTPS", "https://ctf.example.com", "https://ctf.example.com", true, http.SameSiteLaxMode},
		{"cross-origin HTTPS", "https://api.example.com", "https://ctf.example.com", true, http.SameSiteNoneMode},
		{"no Origin", "http://cbctf.local", "", false, http.SameSiteLaxMode},
	} {
		t.Run(test.name, func(t *testing.T) {
			config.Env = &config.Config{Host: test.host}
			config.Env.Gin.Origins = []string{test.origin}
			for _, logout := range []bool{false, true} {
				recorder := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(recorder)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/login", nil)
				ctx.Request.Header.Set("Origin", test.origin)
				if logout {
					Logout(ctx)
				} else {
					setAuthCookie(ctx, "token")
				}
				cookies := recorder.Result().Cookies()
				if len(cookies) != 1 || cookies[0].Secure != test.secure || cookies[0].SameSite != test.sameSite {
					t.Fatalf("logout=%v: unexpected cookies: %+v", logout, cookies)
				}
			}
		})
	}
}
