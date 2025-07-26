package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/keshu12345/overlap-avalara/constants"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewLogger)


type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type logrusLogger struct {
	entry *logrus.Entry
}

func NewLogger() Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.SetReportCaller(true)

	// Create log directory with date folder and PID file
	pid := os.Getpid()
	date := time.Now().Format("2006-01-02")
	logDir := filepath.Join("logs", date)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create log dir: %v", err))
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", constants.OverlapFile.String(), pid))
	log.Out = &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // enable gzip compression
	}

	return &logrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// Implement Logger interface
func (l *logrusLogger) Info(args ...interface{})                 { l.entry.Info(args...) }
func (l *logrusLogger) Infof(format string, args ...interface{}) { l.entry.Infof(format, args...) }

func (l *logrusLogger) Error(args ...interface{})                 { l.entry.Error(args...) }
func (l *logrusLogger) Errorf(format string, args ...interface{}) { l.entry.Errorf(format, args...) }

func (l *logrusLogger) Warn(args ...interface{})                 { l.entry.Warn(args...) }
func (l *logrusLogger) Warnf(format string, args ...interface{}) { l.entry.Warnf(format, args...) }

func (l *logrusLogger) Debug(args ...interface{})                 { l.entry.Debug(args...) }
func (l *logrusLogger) Debugf(format string, args ...interface{}) { l.entry.Debugf(format, args...) }
