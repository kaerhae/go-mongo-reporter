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
func InitLogger(useSyslog bool, logType syslog.Priority) *log.Logger {
	if useSyslog {
		syslogger, err := syslog.NewLogger(logType, 0)
		if err != nil {
			log.Fatalf("Error while initializing logger: %v", err)
		}
		syslogger.SetFlags(log.Ldate | log.Ltime)
		mw := io.MultiWriter(syslogger.Writer(), os.Stdout)
		syslogger.SetOutput(mw)

		return syslogger
	}
	return log.Default()
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

func NewSyslogger(useSyslog bool) Logger {
	return &logger{
		InfoLogger:    InitLogger(useSyslog, syslog.LOG_INFO),
		WarningLogger: InitLogger(useSyslog, syslog.LOG_WARNING),
		ErrorLogger:   InitLogger(useSyslog, syslog.LOG_ERR),
	}
}
