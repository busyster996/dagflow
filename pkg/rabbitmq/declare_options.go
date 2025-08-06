package rabbitmq

// DeclareOptions are used to describe a declare configuration.
// Logger is a custom logging interface.
type DeclareOptions struct {
	Logger      Logger
	ConfirmMode bool
}

// getDefaultChannelOptions describes the options that will be used when a value isn't provided
func getDefaultDeclareOptions() DeclareOptions {
	return DeclareOptions{
		Logger:      stdDebugLogger{},
		ConfirmMode: false,
	}
}

// WithDeclareOptionsLogging sets logging to true on the channel options
// and sets the
func WithDeclareOptionsLogging(options *DeclareOptions) {
	options.Logger = &stdDebugLogger{}
}

// WithDeclareOptionsLogger sets logging to a custom interface.
// Use WithDeclareOptionsLogging to just log to stdout.
func WithDeclareOptionsLogger(log Logger) func(options *DeclareOptions) {
	return func(options *DeclareOptions) {
		options.Logger = log
	}
}

// WithDeclareOptionsConfirm enables confirm mode on the connection
// this is required if channel confirmations should be used
func WithDeclareOptionsConfirm(options *DeclareOptions) {
	options.ConfirmMode = true
}
