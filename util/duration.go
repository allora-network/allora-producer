package util

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogExecutionTime(start time.Time, name string, fields interface{}, logger *zerolog.Logger) {
	elapsed := time.Since(start)
	if logger == nil {
		log.Info().Str("name", name).Fields(fields).Dur("elapsed (ms)", elapsed).Msg("execution time")
	} else {
		logger.Info().Str("name", name).Fields(fields).Dur("elapsed (ms)", elapsed).Msg("execution time")
	}
}
