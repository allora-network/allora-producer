package util

import (
	"time"

	"github.com/rs/zerolog/log"
)

func LogExecutionTime(start time.Time, name string, fields interface{}) {
	elapsed := time.Since(start)
	log.Info().Str("name", name).Fields(fields).Dur("elapsed (ms)", elapsed).Msg("execution time")
}
