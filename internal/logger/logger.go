package logger

import (
	"os"
	"time"
)

type Logger interface {
	Log(message LogMessage)
}

type LogMessage struct {
	Time    time.Time
	Level   string
	Sender  string
	Message string
}

type loggerImpl struct {
	File *os.File
}

func NewLogger() (Logger, error) {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &loggerImpl{
		File: file,
	}, nil
}

func (l *loggerImpl) Log(message LogMessage) {
	logEntry := "[" + message.Level + "] "
	if message.Sender != "" {
		logEntry += "[" + message.Sender + "] "
	}
	logEntry += message.Message + "\n"

	l.File.WriteString(logEntry)
}
