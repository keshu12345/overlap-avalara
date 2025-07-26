package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keshu12345/overlap-avalara/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

type MockOverlapService struct {
	mock.Mock
}

func (m *MockOverlapService) Check(r1, r2 data.DateRange) bool {
	args := m.Called(r1, r2)
	return args.Bool(0)
}

func setupTestRouter() (*gin.Engine, *MockOverlapService, *MockLogger) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockOverlapService{}
	mockLogger := &MockLogger{}

	RegisterEndpoint(router, mockService, mockLogger)

	return router, mockService, mockLogger
}

func createDateRange(start, end string) data.DateRange {
	startTime, _ := time.Parse(time.RFC3339, start)
	endTime, _ := time.Parse(time.RFC3339, end)
	return data.DateRange{
		Start: startTime,
		End:   endTime,
	}
}

func TestCheckOverlap_ValidOverlappingRanges(t *testing.T) {
	router, mockService, mockLogger := setupTestRouter()

	request := data.OverlapRequest{
		Range1: createDateRange("2025-07-01T10:00:00Z", "2025-07-01T12:00:00Z"),
		Range2: createDateRange("2025-07-01T11:00:00Z", "2025-07-01T13:00:00Z"),
	}

	mockService.On("Check", request.Range1, request.Range2).Return(true)
	mockLogger.On("Infof", "isOverlap the time range %v", []interface{}{true}).Return()

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/overlap-check", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Equal(t, true, response["data"])

	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCheckOverlap_ValidNonOverlappingRanges(t *testing.T) {
	router, mockService, mockLogger := setupTestRouter()

	// Create test request with non-overlapping ranges
	request := data.OverlapRequest{
		Range1: createDateRange("2025-07-01T10:00:00Z", "2025-07-01T11:00:00Z"),
		Range2: createDateRange("2025-07-01T12:00:00Z", "2025-07-01T13:00:00Z"),
	}

	// Set up mock expectations
	mockService.On("Check", request.Range1, request.Range2).Return(false)
	mockLogger.On("Infof", "isOverlap the time range %v", []interface{}{false}).Return()

	// Create HTTP request
	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/overlap-check", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Equal(t, false, response["data"])

	// Verify mocks were called
	mockService.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCheckOverlap_EdgeCaseDateRanges(t *testing.T) {
	testCases := []struct {
		name           string
		range1Start    string
		range1End      string
		range2Start    string
		range2End      string
		expectedResult bool
		description    string
	}{
		{
			name:           "Identical Ranges",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "Identical time ranges should overlap",
		},
		{
			name:           "Adjacent Ranges (No Overlap)",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: false,
			description:    "Adjacent ranges with same boundary should not overlap",
		},
		{
			name:           "One Range Inside Another",
			range1Start:    "2025-07-01T09:00:00Z",
			range1End:      "2025-07-01T15:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "Range inside another should overlap",
		},
		{
			name:           "One Minute Overlap",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:01:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "One minute overlap should be detected",
		},
		{
			name:           "Same Start Different End",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "Same start time should overlap",
		},
		{
			name:           "Same End Different Start",
			range1Start:    "2025-07-01T09:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T11:00:00Z",
			expectedResult: true,
			description:    "Same end time should overlap",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, mockService, mockLogger := setupTestRouter()

			// Create test request
			request := data.OverlapRequest{
				Range1: createDateRange(tc.range1Start, tc.range1End),
				Range2: createDateRange(tc.range2Start, tc.range2End),
			}

			// Set up mock expectations
			mockService.On("Check", request.Range1, request.Range2).Return(tc.expectedResult)
			mockLogger.On("Infof", "isOverlap the time range %v", []interface{}{tc.expectedResult}).Return()

			// Create HTTP request
			requestBody, _ := json.Marshal(request)
			req, _ := http.NewRequest("POST", "/api/v1/overlap-check", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, http.StatusOK, w.Code, tc.description)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Contains(t, response, "data")
			assert.Equal(t, tc.expectedResult, response["data"], tc.description)

			// Verify mocks were called
			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestCheckOverlap_DifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(fmt.Sprintf("Method_%s", method), func(t *testing.T) {
			router, mockService, mockLogger := setupTestRouter()

			// Create request with different HTTP method
			req, _ := http.NewRequest(method, "/api/v1/overlap-check", nil)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 405 Method Not Allowed or 404 Not Found
			assert.True(t, w.Code == http.StatusMethodNotAllowed || w.Code == http.StatusNotFound)

			// Service and logger should not be called
			mockService.AssertNotCalled(t, "Check")
			mockLogger.AssertNotCalled(t, "Infof")
			mockLogger.AssertNotCalled(t, "Errorf")
		})
	}
}

func TestRegisterEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockOverlapService{}
	mockLogger := &MockLogger{}

	// Test that registration doesn't panic
	require.NotPanics(t, func() {
		RegisterEndpoint(router, mockService, mockLogger)
	})

	// Test that the route exists
	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Path == "/api/v1/overlap-check" && route.Method == "POST" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected route /api/v1/overlap-check POST should be registered")
}

func TestCheckOverlap_ConcurrentRequests(t *testing.T) {
	router, mockService, mockLogger := setupTestRouter()

	// Set up mock expectations for multiple calls
	mockService.On("Check", mock.AnythingOfType("data.DateRange"), mock.AnythingOfType("data.DateRange")).Return(true)
	mockLogger.On("Infof", "isOverlap the time range %v", []interface{}{true}).Return()

	// Create test request
	request := data.OverlapRequest{
		Range1: createDateRange("2025-07-01T10:00:00Z", "2025-07-01T12:00:00Z"),
		Range2: createDateRange("2025-07-01T11:00:00Z", "2025-07-01T13:00:00Z"),
	}
	requestBody, _ := json.Marshal(request)

	// Make concurrent requests
	const numRequests = 10
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req, _ := http.NewRequest("POST", "/api/v1/overlap-check", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			results <- w.Code
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		statusCode := <-results
		assert.Equal(t, http.StatusOK, statusCode)
	}
}

func TestCheckOverlap_ValidRequests(t *testing.T) {
	testCases := []struct {
		name           string
		range1Start    string
		range1End      string
		range2Start    string
		range2End      string
		expectedResult bool
		description    string
	}{
		{
			name:           "Overlapping Ranges",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T13:00:00Z",
			expectedResult: true,
			description:    "Standard overlapping ranges should return true",
		},
		{
			name:           "Non-overlapping Ranges",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T12:00:00Z",
			range2End:      "2025-07-01T13:00:00Z",
			expectedResult: false,
			description:    "Non-overlapping ranges should return false",
		},
		{
			name:           "Adjacent Ranges",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: false,
			description:    "Adjacent ranges should return false",
		},
		{
			name:           "Identical Ranges",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "Identical ranges should return true",
		},
		{
			name:           "One Range Contains Another",
			range1Start:    "2025-07-01T09:00:00Z",
			range1End:      "2025-07-01T15:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "When one range contains another, should return true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()
			mockService := &MockOverlapService{}
			mockLogger := &MockLogger{}

			// Register endpoints
			RegisterEndpoint(router, mockService, mockLogger)

			// Create test request
			request := data.OverlapRequest{
				Range1: createDateRange(tc.range1Start, tc.range1End),
				Range2: createDateRange(tc.range2Start, tc.range2End),
			}

			mockService.On("Check", request.Range1, request.Range2).Return(tc.expectedResult)

			mockLogger.On("Infof", "isOverlap the time range %v", mock.Anything).Return()

			requestBody, _ := json.Marshal(request)
			req, _ := http.NewRequest("POST", "/api/v1/overlap-check", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.description)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			if data, exists := response["data"]; exists {
				assert.Equal(t, tc.expectedResult, data, tc.description)
			}

			mockService.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestCheckOverlap_InvalidJSON(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
		description string
	}{
		{
			name:        "Malformed JSON",
			requestBody: `{"range1": {"start": "2025-07-01T10:00:00Z", "end": "2025-07-01T12:00:00Z"`,
			description: "Malformed JSON should return BadRequest",
		},
		{
			name:        "Invalid Date Format",
			requestBody: `{"range1": {"start": "invalid-date", "end": "2025-07-01T12:00:00Z"}, "range2": {"start": "2025-07-01T11:00:00Z", "end": "2025-07-01T13:00:00Z"}}`,
			description: "Invalid date format should return BadRequest",
		},
		{
			name:        "Missing Required Fields",
			requestBody: `{"range1": {"start": "2025-07-01T10:00:00Z"}}`, 
			description: "Missing required fields should return BadRequest",
		},
		{
			name:        "Empty JSON Object",
			requestBody: `{}`,
			description: "Empty JSON object should return BadRequest",
		},
		{
			name:        "Null Values",
			requestBody: `{"range1": null, "range2": null}`,
			description: "Null values should return BadRequest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()
			mockService := &MockOverlapService{}
			mockLogger := &MockLogger{}

		
			RegisterEndpoint(router, mockService, mockLogger)

			mockLogger.On("Errorf", "Unable to bind with json body :%v", mock.Anything).Return()
			req, _ := http.NewRequest("POST", "/api/v1/overlap-check", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, tc.description)

			mockLogger.AssertExpectations(t)

			mockService.AssertNotCalled(t, "Check")
		})
	}
}



func TestCheckOverlap_MissingFields(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
	}{
		{
			name:        "Missing Range1",
			requestBody: `{"range2": {"start": "2025-07-01T11:00:00Z", "end": "2025-07-01T13:00:00Z"}}`,
		},
		{
			name:        "Missing Range2",
			requestBody: `{"range1": {"start": "2025-07-01T10:00:00Z", "end": "2025-07-01T12:00:00Z"}}`,
		},
		{
			name:        "Missing Start in Range1",
			requestBody: `{"range1": {"end": "2025-07-01T12:00:00Z"}, "range2": {"start": "2025-07-01T11:00:00Z", "end": "2025-07-01T13:00:00Z"}}`,
		},
		{
			name:        "Missing End in Range2",
			requestBody: `{"range1": {"start": "2025-07-01T10:00:00Z", "end": "2025-07-01T12:00:00Z"}, "range2": {"start": "2025-07-01T11:00:00Z"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()
			mockService := &MockOverlapService{}
			mockLogger := &MockLogger{}

			RegisterEndpoint(router, mockService, mockLogger)

			mockLogger.On("Errorf", "Unable to bind with json body :%v", mock.Anything).Return()

			req, _ := http.NewRequest("POST", "/api/v1/overlap-check", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			mockLogger.AssertExpectations(t)
			mockService.AssertNotCalled(t, "Check")
		})
	}
}