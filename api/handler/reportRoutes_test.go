package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"main/pkg/helpers"
	"main/pkg/middleware"
	"main/pkg/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockReportService struct {
	Repository helpers.MockReportRepository
}

func (s *MockReportService) GetAllReports() ([]models.Report, error) {
	return []models.Report{
		{Topic: "test", Author: "Test"},
	}, nil
}

func (s *MockReportService) GetSingleReport(id string) (models.Report, error) {
	objID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
	return models.Report{
		ID:     objID,
		Topic:  "test",
		Author: "Test",
		UserID: "123456789012345678901234",
	}, nil
}

func (s *MockReportService) CreateReport(report models.Report) (string, error) {
	return "", nil
}

func (s *MockReportService) UpdateReport(newReport models.Report) error {
	return nil
}

func (s *MockReportService) DeleteReport(id string) (int64, error) {
	return int64(0), nil
}

func (s *MockReportService) UpdateReportReferences(userID primitive.ObjectID, reportID primitive.ObjectID) error {
	return nil
}

func TestGetAllReports(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	r.GET("/api/reports", reportHandler.Get)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/reports", nil)
	r.ServeHTTP(w, req)

	var reports []models.Report
	json.Unmarshal(w.Body.Bytes(), &reports)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, reports)
}

func TestGetSingleReport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	r.GET("/api/reports/123", reportHandler.GetByID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/reports/123", nil)
	r.ServeHTTP(w, req)

	objID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
	shouldBe := models.Report{
		ID: objID, Topic: "test", Author: "Test", UserID: "123456789012345678901234",
	}
	var reports models.Report
	json.Unmarshal(w.Body.Bytes(), &reports)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, shouldBe, reports)
}

func TestPostReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)

	r := gin.Default()
	r.POST("/api/reports", reportHandler.Post)
	test := models.Report{
		Topic:       "test",
		Author:      "test",
		Description: "test",
		UserID:      "111122223333444455556666",
	}
	reportJson, _ := json.Marshal(test)
	req, _ := http.NewRequest(http.MethodPost, "/api/reports", strings.NewReader(string(reportJson)))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPostReport_ShouldErrorIfNoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)

	r := gin.Default()
	r.POST("/api/reports", reportHandler.Post)
	test := models.Report{
		Topic:       "test",
		Author:      "test",
		Description: "test",
	}
	reportJson, _ := json.Marshal(test)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPost, "/api/reports", strings.NewReader(string(reportJson)))

	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestUpdateReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	w := httptest.NewRecorder()
	content := models.Report{
		Topic:       "",
		Description: "",
		Author:      "",
		UserID:      "123",
	}
	reportJson, _ := json.Marshal(content)

	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "123")
	c.Set("isAdmin", true)
	c.Request = &http.Request{
		Method: http.MethodPut,
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(bytes.NewBuffer(reportJson)),
	}
	reportHandler.Update(c)
	assert.Equal(t, 200, w.Code)
}

func TestUpdateReport_ShouldFailWhenUserIdNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	w := httptest.NewRecorder()
	content := models.Report{
		Topic:       "",
		Description: "",
		Author:      "",
		UserID:      "123",
	}
	reportJson, _ := json.Marshal(content)

	c, _ := gin.CreateTestContext(w)
	c.Set("isAdmin", true)
	c.Request = &http.Request{
		Method: http.MethodPut,
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(bytes.NewBuffer(reportJson)),
	}
	reportHandler.Update(c)
	assert.Equal(t, 500, w.Code)
}

func TestUpdateReport_ShouldFailWhenIsAdminNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	w := httptest.NewRecorder()
	content := models.Report{
		Topic:       "",
		Description: "",
		Author:      "",
		UserID:      "123",
	}
	reportJson, _ := json.Marshal(content)

	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "123")

	c.Request = &http.Request{
		Method: http.MethodPut,
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(bytes.NewBuffer(reportJson)),
	}
	reportHandler.Update(c)
	assert.Equal(t, 500, w.Code)
}

func TestDeleteReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "123")
	c.Set("isAdmin", true)
	c.Request = &http.Request{
		Method: http.MethodPut,
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	c.Params = gin.Params{
		{Key: "id", Value: "12345"},
	}

	reportHandler.Delete(c)
	assert.Equal(t, 200, w.Code)
}

func TestDeleteReport_ShouldFailWhenUserIdNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	w := httptest.NewRecorder()
	content := models.Report{
		Topic:       "",
		Description: "",
		Author:      "",
		UserID:      "123",
	}
	reportJson, _ := json.Marshal(content)

	c, _ := gin.CreateTestContext(w)
	c.Set("isAdmin", true)
	c.Request = &http.Request{
		Method: http.MethodPut,
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(bytes.NewBuffer(reportJson)),
	}
	reportHandler.Update(c)
	assert.Equal(t, 500, w.Code)
}

func TestDeleteReport_ShouldFailWhenIsAdminNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := helpers.InitMockReportRepository()
	s := &MockReportService{Repository: repo}
	l := middleware.NewSyslogger(false)
	reportHandler := NewReportRouter(s, l)
	w := httptest.NewRecorder()
	content := models.Report{
		Topic:       "",
		Description: "",
		Author:      "",
		UserID:      "123",
	}
	reportJson, _ := json.Marshal(content)

	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "123")

	c.Request = &http.Request{
		Method: http.MethodPut,
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(bytes.NewBuffer(reportJson)),
	}
	reportHandler.Update(c)
	assert.Equal(t, 500, w.Code)
}

func TestConvertStringToObjectID(t *testing.T) {
	shouldBe, _ := primitive.ObjectIDFromHex("111122223333444455556666")

	s, err := convertStringToPrimitiveID("111122223333444455556666")
	if err != nil {
		t.Errorf("Failed on: %v", err)
	}

	assert.Equal(t, shouldBe, s)
}
