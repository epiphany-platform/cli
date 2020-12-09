package cmd

func debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}
