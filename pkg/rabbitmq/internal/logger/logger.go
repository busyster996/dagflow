package logger

// Logger is describes a logging structure. It can be set using
// WithPublisherOptionsLogger() or WithConsumerOptionsLogger().
type Logger interface {
	Errorf(format string, args ...any)
	Warnf(format string, args ...any)
	Infof(format string, args ...any)
	Debugf(format string, args ...any)
}
