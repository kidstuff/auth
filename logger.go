package auth

type Logger interface {
	// Debugf formats its arguments according to the format, analogous to fmt.Printf,
	// and records the text as a log message at Debug level.
	Debugf(format string, args ...interface{})
	// Infof is like Debugf, but at Info level.
	Infof(format string, args ...interface{})
	// Warningf is like Debugf, but at Warning level.
	Warningf(format string, args ...interface{})
	// Errorf is like Debugf, but at Error level.
	Errorf(format string, args ...interface{})
	// Criticalf is like Debugf, but at Critical level.
	Criticalf(format string, args ...interface{})
}
