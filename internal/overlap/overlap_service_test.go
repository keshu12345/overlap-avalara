package overlap

import (
	"testing"
	"time"

	"github.com/keshu12345/overlap-avalara/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func mustParseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}
	return t
}


func createDateRange(start, end string) data.DateRange {
	return data.DateRange{
		Start: mustParseTime(start),
		End:   mustParseTime(end),
	}
}

func TestOverlapService_Check(t *testing.T) {
	testCases := []struct {
		name           string
		range1Start    string
		range1End      string
		range2Start    string
		range2End      string
		expectedResult bool
		description    string
	}{
		// Overlapping cases
		{
			name:           "Standard Overlap - Partial",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T13:00:00Z",
			expectedResult: true,
			description:    "Ranges with partial overlap should return true",
		},
		{
			name:           "Range1 Contains Range2",
			range1Start:    "2025-07-01T09:00:00Z",
			range1End:      "2025-07-01T15:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: true,
			description:    "When range1 completely contains range2, should return true",
		},
		{
			name:           "Range2 Contains Range1",
			range1Start:    "2025-07-01T11:00:00Z",
			range1End:      "2025-07-01T13:00:00Z",
			range2Start:    "2025-07-01T09:00:00Z",
			range2End:      "2025-07-01T15:00:00Z",
			expectedResult: true,
			description:    "When range2 completely contains range1, should return true",
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
			name:           "Overlap at Start",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T08:00:00Z",
			range2End:      "2025-07-01T11:00:00Z",
			expectedResult: true,
			description:    "Ranges overlapping at the start should return true",
		},
		{
			name:           "Overlap at End",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T14:00:00Z",
			expectedResult: true,
			description:    "Ranges overlapping at the end should return true",
		},
		{
			name:           "Minimal Overlap - One Second",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T12:00:00Z",
			range2Start:    "2025-07-01T11:59:59Z",
			range2End:      "2025-07-01T13:00:00Z",
			expectedResult: true,
			description:    "Ranges with minimal overlap should return true",
		},

		// Non-overlapping cases
		{
			name:           "Non-overlapping - Gap Between",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T12:00:00Z",
			range2End:      "2025-07-01T13:00:00Z",
			expectedResult: false,
			description:    "Non-overlapping ranges with gap should return false",
		},
		{
			name:           "Adjacent Ranges - Touching Endpoints",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: false,
			description:    "Adjacent ranges (touching endpoints) should return false",
		},
		{
			name:           "Range1 Before Range2",
			range1Start:    "2025-07-01T08:00:00Z",
			range1End:      "2025-07-01T09:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: false,
			description:    "When range1 is completely before range2, should return false",
		},
		{
			name:           "Range2 Before Range1",
			range1Start:    "2025-07-01T12:00:00Z",
			range1End:      "2025-07-01T14:00:00Z",
			range2Start:    "2025-07-01T08:00:00Z",
			range2End:      "2025-07-01T10:00:00Z",
			expectedResult: false,
			description:    "When range2 is completely before range1, should return false",
		},
		{
			name:           "Adjacent with One Second Gap",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T11:00:01Z",
			range2End:      "2025-07-01T12:00:00Z",
			expectedResult: false,
			description:    "Ranges with one second gap should return false",
		},

		// Edge cases
		{
			name:           "Zero Duration Range1",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T10:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T11:00:00Z",
			expectedResult: false,
			description:    "Zero duration range1 at start of range2 should return false",
		},
		{
			name:           "Zero Duration Range2",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T11:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T10:00:00Z",
			expectedResult: false,
			description:    "Zero duration range2 at start of range1 should return false",
		},
		{
			name:           "Both Zero Duration - Same Time",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T10:00:00Z",
			range2Start:    "2025-07-01T10:00:00Z",
			range2End:      "2025-07-01T10:00:00Z",
			expectedResult: false,
			description:    "Both zero duration at same time should return false",
		},
		{
			name:           "Both Zero Duration - Different Times",
			range1Start:    "2025-07-01T10:00:00Z",
			range1End:      "2025-07-01T10:00:00Z",
			range2Start:    "2025-07-01T11:00:00Z",
			range2End:      "2025-07-01T11:00:00Z",
			expectedResult: false,
			description:    "Both zero duration at different times should return false",
		},

		// Cross-day scenarios
		{
			name:           "Cross Day Overlap",
			range1Start:    "2025-07-01T23:00:00Z",
			range1End:      "2025-07-02T02:00:00Z",
			range2Start:    "2025-07-02T01:00:00Z",
			range2End:      "2025-07-02T03:00:00Z",
			expectedResult: true,
			description:    "Ranges spanning across days with overlap should return true",
		},
		{
			name:           "Cross Day Non-overlap",
			range1Start:    "2025-07-01T22:00:00Z",
			range1End:      "2025-07-01T23:00:00Z",
			range2Start:    "2025-07-02T01:00:00Z",
			range2End:      "2025-07-02T03:00:00Z",
			expectedResult: false,
			description:    "Ranges spanning across days without overlap should return false",
		},

		// Different time zones (but using UTC for consistency)
		{
			name:           "Long Duration Ranges",
			range1Start:    "2025-07-01T00:00:00Z",
			range1End:      "2025-07-10T00:00:00Z",
			range2Start:    "2025-07-05T00:00:00Z",
			range2End:      "2025-07-15T00:00:00Z",
			expectedResult: true,
			description:    "Long duration ranges with overlap should return true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			service := New(mockLogger)

			mockLogger.On("Info", mock.Anything).Return()

			range1 := createDateRange(tc.range1Start, tc.range1End)
			range2 := createDateRange(tc.range2Start, tc.range2End)

			result := service.Check(range1, range2)

			assert.Equal(t, tc.expectedResult, result, tc.description)

			mockLogger.AssertExpectations(t)
		})
	}
}

func TestOverlapService_New(t *testing.T) {
	t.Run("Constructor Creates Service Correctly", func(t *testing.T) {
		// Setup
		mockLogger := &MockLogger{}

		// Execute
		service := New(mockLogger)

		// Assertions
		assert.NotNil(t, service)

		// Verify it implements the interface
		var _ OverlapService = service
	})
}

func BenchmarkOverlapService_Check_Overlapping(b *testing.B) {
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything).Return()
	service := New(mockLogger)

	range1 := createDateRange("2025-07-01T10:00:00Z", "2025-07-01T12:00:00Z")
	range2 := createDateRange("2025-07-01T11:00:00Z", "2025-07-01T13:00:00Z")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Check(range1, range2)
	}
}

func BenchmarkOverlapService_Check_NonOverlapping(b *testing.B) {
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything).Return()
	service := New(mockLogger)

	range1 := createDateRange("2025-07-01T10:00:00Z", "2025-07-01T11:00:00Z")
	range2 := createDateRange("2025-07-01T12:00:00Z", "2025-07-01T13:00:00Z")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Check(range1, range2)
	}
}
