package log

import (
	"io"

	"github.com/rs/zerolog"
)

func New(level zerolog.Level) *zerolog.Logger {
	var output io.Writer

	logger := zerolog.New(output).With().Timestamp().Logger().Level(level)
	return &logger
}
