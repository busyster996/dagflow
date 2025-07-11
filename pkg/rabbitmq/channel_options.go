package rabbitmq

// ChannelOptions are used to describe a channel's configuration.
// Logger is a custom logging interface.
type ChannelOptions struct {
	Logger      Logger
	QOSPrefetch int
	QOSGlobal   bool
	ConfirmMode bool
}

// getDefaultChannelOptions describes the options that will be used when a value isn't provided
func getDefaultChannelOptions() ChannelOptions {
	return ChannelOptions{
		Logger:      stdDebugLogger{},
		QOSPrefetch: 10,
		QOSGlobal:   false,
	}
}

// WithChannelOptionsQOSPrefetch returns a function that sets the prefetch count, which means that
// many messages will be fetched from the server in advance to help with throughput.
// This doesn't affect the handler, messages are still processed one at a time.
func WithChannelOptionsQOSPrefetch(prefetchCount int) func(*ChannelOptions) {
	return func(options *ChannelOptions) {
		options.QOSPrefetch = prefetchCount
	}
}

// WithChannelOptionsQOSGlobal sets the qos on the channel to global, which means
// these QOS settings apply to ALL existing and future
// channels on the same connection
func WithChannelOptionsQOSGlobal(options *ChannelOptions) {
	options.QOSGlobal = true
}

// WithChannelOptionsLogging sets logging to true on the channel options
// and sets the
func WithChannelOptionsLogging(options *ChannelOptions) {
	options.Logger = &stdDebugLogger{}
}

// WithChannelOptionsLogger sets logging to a custom interface.
// Use WithChannelOptionsLogging to just log to stdout.
func WithChannelOptionsLogger(log Logger) func(options *ChannelOptions) {
	return func(options *ChannelOptions) {
		options.Logger = log
	}
}

// WithChannelOptionsConfirm enables confirm mode on the connection
// this is required if channel confirmations should be used
func WithChannelOptionsConfirm(options *ChannelOptions) {
	options.ConfirmMode = true
}
