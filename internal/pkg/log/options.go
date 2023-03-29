package log

import "go.uber.org/zap/zapcore"

type Options struct {
	// DisableCaller disables automatic caller annotation.
	DisableCaller bool
	// DisableStacktrace disables automatic stacktrace annotation.
	DisableStacktrace bool
	// Level is the minimum enabled logging level.
	Level string
	// Format is the output format of the logger.
	Format string
	// OutputPaths is a list of URLs or file paths to write logging output to.
	OutputPaths []string
}

func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             zapcore.InfoLevel.String(),
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}
