package logger

func ExampleNewLogger() {
	// "kayos_logger_example-1721120070.log" will be created in the current directory
	f, _ := CreateDatedLogFile("./", "kayos_logger_example")

	log := NewLogger(f)         // omit 'f' to only log to console
	log.C().SetPrefix(f.Name()) // will apply to both [Z] and [C] calls
	log.Z().Info().Msg("hello") // zerolog syntax
	log.C().Info("world")       // stdlib "log" syntax
}
