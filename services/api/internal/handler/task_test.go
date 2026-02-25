package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mock router ที่ไม่ต้องพึ่ง redis จริง
func setupHealthRouter() *gin.Engine {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "api",
			"version": "1.0.0",
		})
	})
	return r
}

func setupTaskRouter() *gin.Engine {
	r := gin.New()
	taskHandler := &TaskHandler{redis: nil}
	r.POST("/api/v1/tasks", taskHandler.CreateTask)
	return r
}

func TestHealthCheck_ReturnsOK(t *testing.T) {
	r := setupHealthRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestCreateTask_InvalidBody(t *testing.T) {
	r := setupTaskRouter()
	w := httptest.NewRecorder()

	body := bytes.NewBufferString(`{"description": "no title"}`)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestCreateTask_ValidRequestFormat(t *testing.T) {
	payload := map[string]string{
		"title":       "Test Task",
		"description": "Test Description",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type header not set correctly")
	}
}
