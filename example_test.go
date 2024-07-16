package logger

import (
	"time"
)

func ExampleNewLogger() {
	f, _ := CreateDatedLogFile("./", "kayos_logger_example")
	log := NewLogger(f)
	log.C().SetPrefix(f.Name())
	log.Z().Info().Msg("hello")
	log.C().Info("world")
	time.Sleep(50 * time.Millisecond)
}
