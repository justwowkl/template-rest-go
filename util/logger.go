package util

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger asdf
var Logger zerolog.Logger

func loggerInit() {
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
