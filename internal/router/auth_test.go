package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func InitAuthTest() *gin.Engine {
	i18n.Init()
	config.Env = viper.New()
	config.Env.Set("gorm.file", ":memory:")
	config.Env.Set("gorm.log.level", "silent")
	config.Env.Set("log.level", "info")
	config.Env.Set("log.file", false)
	config.Env.Set("gin.mode", "release")
	log.Init()
	db.Init()
	var ctx context.Context
	db.CreateUser(ctx, "user1", "password", "user1@0rays.club")
	db.CreateAdmin(ctx, "admin1", "password", "admin1@0rays.club")
	return Init()
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := InitAuthTest()

	w := httptest.NewRecorder()
	body := `{"name":"user2","password":"password","email":"user2@0rays.club"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	body = `{"name":"user1","password":"password","email":"user1@0rays.club"}`
	req, _ = http.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d but got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := InitAuthTest()

	w := httptest.NewRecorder()
	body := `{"name":"user1","password":"password"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	body = `{"name":"user1@0rays.club","password":"password"}`
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	body = `{"name":"user1","password":"error_pwd"}`
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d but got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAdminLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := InitAuthTest()

	w := httptest.NewRecorder()
	body := `{"name":"admin1","password":"password"}`
	req, _ := http.NewRequest("POST", "/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	body = `{"name":"admin1@0rays.club","password":"password"}`
	req, _ = http.NewRequest("POST", "/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	body = `{"name":"admin1","password":"error_pwd"}`
	req, _ = http.NewRequest("POST", "/admin/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d but got %d", http.StatusUnauthorized, w.Code)
	}
}
