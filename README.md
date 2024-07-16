# logger

[![Go Reference](https://pkg.go.dev/badge/git.tcp.direct/kayos/logger.svg)](https://pkg.go.dev/git.tcp.direct/kayos/logger)

Convenience wrapper for [zerolog](https://github.com/rs/zerolog) + [zwrap](https://github.com/yunginnanet/zwrap) along with log file helper utilities.

## Basic Usage

```golang
package logger

import (
	"git.tcp.direct/kayos/logger"
)

func ExampleNewLogger() {
	// "kayos_logger_example-1721120070.log" will be created in the current directory
	f, _ := CreateDatedLogFile("./", "kayos_logger_example")

	log := NewLogger(f)         // omit 'f' to only log to console
	log.C().SetPrefix(f.Name()) // will apply to both [Z] and [C] calls
	log.Z().Info().Msg("hello") // zerolog syntax
	log.C().Info("world")       // stdlib "log" syntax
}
```
