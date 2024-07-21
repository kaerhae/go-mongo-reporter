package middleware

import (
	"io"
	"log"
	"log/syslog"
	"os"
)

type Logger interface {
	LogInfo(message string)
	LogWarning(message string)
	LogError(message string)
}

type logger struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

// LogError implements Logger.
func InitLogger(logType syslog.Priority) *log.Logger {
	syslogger, err := syslog.NewLogger(logType, 0)
	if err != nil {
		log.Fatalf("Error while initializing logger: %v", err)
	}
	syslogger.SetFlags(log.Ldate | log.Ltime)
	mw := io.MultiWriter(syslogger.Writer(), os.Stdout)
	syslogger.SetOutput(mw)

	return syslogger
}

// LogInfo implements Logger.
func (l *logger) LogInfo(message string) {
	l.InfoLogger.Printf("[INFO]: %s", message)
}

// LogWarning implements Logger.
func (l *logger) LogWarning(message string) {
	l.InfoLogger.Printf("[WARNING]: %s", message)
}

// LogError implements Logger.
func (l *logger) LogError(message string) {
	l.ErrorLogger.Printf("[ERROR]: %s", message)
}

func NewSyslogger() Logger {
	return &logger{
		InfoLogger:    InitLogger(syslog.LOG_INFO),
		WarningLogger: InitLogger(syslog.LOG_WARNING),
		ErrorLogger:   InitLogger(syslog.LOG_ERR),
	}
}
