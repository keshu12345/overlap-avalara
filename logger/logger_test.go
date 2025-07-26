package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOverlapFile string

func (m mockOverlapFile) String() string {
	return string(m)
}

func createTestLogger() (*logrusLogger, *bytes.Buffer) {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)

	buffer := &bytes.Buffer{}
	log.Out = buffer

	return &logrusLogger{
		entry: logrus.NewEntry(log),
	}, buffer
}

func TestNewLogger(t *testing.T) {
	defer func() {
		os.RemoveAll("logs")
	}()

	logger := NewLogger()
	assert.NotNil(t, logger)

	_, ok := logger.(Logger)
	assert.True(t, ok, "NewLogger should return a Logger interface")

	pid := os.Getpid()
	date := time.Now().Format("2006-01-02")
	expectedLogDir := filepath.Join("logs", date)

	_, err := os.Stat(expectedLogDir)
	assert.NoError(t, err, "Log directory should be created")
	expectedLogFile := filepath.Join(expectedLogDir, fmt.Sprintf("%s.%d.log", "overlap", pid))
	assert.NotEmpty(t, expectedLogFile)
}

func TestLoggerInterface(t *testing.T) {
	logger := NewLogger()
	assert.Implements(t, (*Logger)(nil), logger)
}

func TestLogrusLogger_Info(t *testing.T) {
	logger, buffer := createTestLogger()

	testMessage := "This is an info message"
	logger.Info(testMessage)

	output := buffer.String()
	assert.Contains(t, output, testMessage)
	assert.Contains(t, output, "level=info")
}

func TestLogrusLogger_Infof(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.Infof("This is an info message with value: %d", 42)

	output := buffer.String()
	assert.Contains(t, output, "This is an info message with value: 42")
	assert.Contains(t, output, "level=info")
}

func TestLogrusLogger_Error(t *testing.T) {
	logger, buffer := createTestLogger()

	testMessage := "This is an error message"
	logger.Error(testMessage)

	output := buffer.String()
	assert.Contains(t, output, testMessage)
	assert.Contains(t, output, "level=error")
}

func TestLogrusLogger_Errorf(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.Errorf("This is an error with code: %d", 500)

	output := buffer.String()
	assert.Contains(t, output, "This is an error with code: 500")
	assert.Contains(t, output, "level=error")
}

func TestLogrusLogger_Warn(t *testing.T) {
	logger, buffer := createTestLogger()

	testMessage := "This is a warning message"
	logger.Warn(testMessage)

	output := buffer.String()
	assert.Contains(t, output, testMessage)
	assert.Contains(t, output, "level=warning")
}

func TestLogrusLogger_Warnf(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.Warnf("Warning: %s is deprecated", "function_name")

	output := buffer.String()
	assert.Contains(t, output, "Warning: function_name is deprecated")
	assert.Contains(t, output, "level=warning")
}

func TestLogrusLogger_Debug(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.entry.Logger.SetLevel(logrus.DebugLevel)

	testMessage := "This is a debug message"
	logger.Debug(testMessage)

	output := buffer.String()
	assert.Contains(t, output, testMessage)
	assert.Contains(t, output, "level=debug")
}

func TestLogrusLogger_Debugf(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.entry.Logger.SetLevel(logrus.DebugLevel)

	logger.Debugf("Debug value: %v", map[string]int{"count": 10})

	output := buffer.String()
	assert.Contains(t, output, "Debug value: map[count:10]")
	assert.Contains(t, output, "level=debug")
}

func TestLoggerWithMultipleArguments(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.Info("Multiple", "arguments", "test", 123)

	output := buffer.String()
	assert.Contains(t, output, "Multiple")
	assert.Contains(t, output, "arguments")
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "123")
}

func TestLoggerTimestamp(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.Info("Timestamp test")

	output := buffer.String()
	currentYear := time.Now().Format("2006")
	assert.Contains(t, output, currentYear)
}

func TestLogDirectoryCreation(t *testing.T) {
	os.RemoveAll("logs")
	defer os.RemoveAll("logs")

	_, err := os.Stat("logs")
	require.True(t, os.IsNotExist(err))

	NewLogger()

	date := time.Now().Format("2006-01-02")
	expectedDir := filepath.Join("logs", date)
	info, err := os.Stat(expectedDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestLogDirectoryCreationFailure(t *testing.T) {
	t.Skip("Testing directory creation failure requires mocking os.MkdirAll")
}

func TestLoggerInterfaceCompleteness(t *testing.T) {
	logger := NewLogger()

	assert.NotPanics(t, func() {
		logger.Info("test")
		logger.Infof("test %s", "formatted")
		logger.Error("test")
		logger.Errorf("test %s", "formatted")
		logger.Warn("test")
		logger.Warnf("test %s", "formatted")
		logger.Debug("test")
		logger.Debugf("test %s", "formatted")
	})
}

func TestLoggerConcurrency(t *testing.T) {
	logger, buffer := createTestLogger()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Infof("Concurrent log %d", id)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	output := buffer.String()

	for i := 0; i < 10; i++ {
		assert.Contains(t, output, fmt.Sprintf("Concurrent log %d", i))
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	logger, _ := createTestLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark test message")
	}
}

func BenchmarkLoggerInfof(b *testing.B) {
	logger, _ := createTestLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("Benchmark test message %d", i)
	}
}

func TestLogLevels(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.entry.Logger.SetLevel(logrus.InfoLevel)

	logger.Debug("This debug should not appear")
	logger.Info("This info should appear")
	logger.Warn("This warning should appear")
	logger.Error("This error should appear")

	output := buffer.String()

	assert.NotContains(t, output, "This debug should not appear")
	assert.Contains(t, output, "This info should appear")
	assert.Contains(t, output, "This warning should appear")
	assert.Contains(t, output, "This error should appear")
}

func TestLogFormatting(t *testing.T) {
	logger, buffer := createTestLogger()

	logger.Info("Test message")

	output := buffer.String()
	assert.Regexp(t, `time="[^"]*"`, output)
	assert.Contains(t, output, `level=info`)
	assert.Contains(t, output, `msg="Test message"`)
	assert.Contains(t, output, `func=`)
	assert.Contains(t, output, `file=`)
}
