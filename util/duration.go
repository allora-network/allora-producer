package util

import (
	"time"

	"github.com/rs/zerolog/log"
)

func LogExecutionTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Str("name", name).Dur("elapsed", elapsed).Msg("execution time")
}
