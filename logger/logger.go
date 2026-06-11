package logger

import (
	"os"
	"time"
)

type Logger struct {
	logBuffer string
	logFile   *os.File

	Listener func(string)
}

var logger *Logger

func (l *Logger) Log(s string) {
	t := time.Now()
	newLine := t.Format("2006-01-02 15:04:05") + " > " + s + "\n"
	l.logBuffer = l.logBuffer + newLine
	l.logFile.WriteString(newLine)
	l.logFile.Sync()
	l.Listener(l.logBuffer)
}

func GetLogger() *Logger {
	if logger == nil {
		logFile, _ := os.OpenFile("C:\\superbot\\log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		logger = &Logger{logBuffer: "", logFile: logFile}
	}
	return logger
}
