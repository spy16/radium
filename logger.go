package radium

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger implementation should provide logging
// functionality to the radium instance. Log levels
// should be managed externally.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// defaultLogger implements Logger using log package
type defaultLogger struct {
	logger *logrus.Logger
}

func (dl defaultLogger) Debugf(format string, args ...interface{}) {
	dl.logger.WithFields(logrus.Fields(getSourceInfoFields())).Debugf(format, args...)
}

func (dl defaultLogger) Infof(format string, args ...interface{}) {
	dl.logger.WithFields(logrus.Fields(getSourceInfoFields())).Infof(format, args...)
}

func (dl defaultLogger) Warnf(format string, args ...interface{}) {
	dl.logger.WithFields(logrus.Fields(getSourceInfoFields())).Warnf(format, args...)
}

func (dl defaultLogger) Errorf(format string, args ...interface{}) {
	dl.logger.WithFields(logrus.Fields(getSourceInfoFields())).Errorf(format, args...)
}

func (dl defaultLogger) Fatalf(format string, args ...interface{}) {
	dl.logger.WithFields(logrus.Fields(getSourceInfoFields())).Fatalf(format, args...)
}

func getSourceInfoFields() map[string]interface{} {
	file, line := getFileInfo(3)
	m := map[string]interface{}{
		"file": file,
		"line": line,
	}
	return m
}

func getFileInfo(subtractStackLevels int) (string, int) {
	_, file, line, _ := runtime.Caller(subtractStackLevels)
	return chopPath(file), line
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i != -1 {
		return original[i+1:]
	}
	return original
}
