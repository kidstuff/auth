// +build !windows

package auth

import (
	"fmt"
	"log/syslog"
)

type SysLogger struct {
	*syslog.Writer
}

func NewSysLogger(tag string) (*SysLogger, error) {
	w, err := syslog.New(syslog.LOG_INFO, tag)
	if err != nil {
		return nil, err
	}

	return &SysLogger{w}, nil
}

// Debugf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as a log message at Debug level.
func (l *SysLogger) Debugf(format string, args ...interface{}) {
	l.Writer.Debug(fmt.Sprintf(format, args...))
}

// Infof is like Debugf, but at Info level.
func (l *SysLogger) Infof(format string, args ...interface{}) {
	l.Writer.Info(fmt.Sprintf(format, args...))
}

// Warningf is like Debugf, but at Warning level.
func (l *SysLogger) Warningf(format string, args ...interface{}) {
	l.Writer.Warning(fmt.Sprintf(format, args...))
}

// Errorf is like Debugf, but at Error level.
func (l *SysLogger) Errorf(format string, args ...interface{}) {
	l.Writer.Err(fmt.Sprintf(format, args...))
}

// Criticalf is like Debugf, but at Critical level.
func (l *SysLogger) Criticalf(format string, args ...interface{}) {
	l.Writer.Crit(fmt.Sprintf(format, args...))
}
